package controller

import (
	"math"
	"runtime"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/hierarchy"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/acceleration"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/shape3d"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/game/level"
	"github.com/mokiat/rally-mka/internal/game/preset"
)

const (
	anchorDistance = 6.0
	cameraDistance = 15.0
	pitchAngle     = 35.0
)

func NewPlayController(window app.Window, engine *game.Engine, playData *data.PlayData) *PlayController {
	return &PlayController{
		window:   window,
		engine:   engine,
		playData: playData,
	}
}

type PlayController struct {
	window   app.Window
	engine   *game.Engine
	playData *data.PlayData

	scene        *preset.Scene
	hierScene    *hierarchy.Scene
	gfxScene     *graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *preset.ECSScene

	followCameraSystem *preset.FollowCameraSystem
	followCamera       *graphics.Camera
	bonnetCamera       *graphics.Camera

	carSystem         *preset.CarSystem
	vehicleDefinition *preset.CarDefinition
	vehicle           *preset.Car
}

func tilePosition(coord level.Coord) dprec.Vec3 {
	const tileSize = 80.0
	x, y := coord.X, coord.Y
	xShift := tileSize * float64(math.Sqrt(3)) / 2.0
	yShift := tileSize * 3.0 / 4.0
	xOffset := float64(0.0)
	if max(y, -y)%2 == 1 {
		xOffset = xShift / 2.0
	}
	return dprec.Vec3{
		X: float64(x)*xShift + xOffset,
		Y: 0.0,
		Z: float64(y) * yShift,
	}
}

func (c *PlayController) Start(environment data.Lighting, controller data.Input, board *level.Board) {
	physics.ImpulseDriftAdjustmentRatio = 0.06 // FIXME: Use default once multi-point collisions are fixed

	c.scene = preset.NewScene(c.engine.CreateScene(game.SceneInfo{}))

	c.scene.InstantiateModel(game.ModelInfo{
		Name:      opt.V("Background"),
		Template:  c.playData.Background,
		IsDynamic: false,
	})

	centerPosition := tilePosition(board.Center())
	for y := range board.Size() {
		for x := range board.Size() {
			tileCoord := level.C(x, y)
			tile := board.Tile(tileCoord)
			nodeName := tile.NodeName()
			if nodeName == "" {
				continue
			}
			tilePosition := tilePosition(tileCoord)
			position := dprec.Vec3Diff(tilePosition, centerPosition)
			c.scene.InstantiateModel(game.ModelInfo{
				SubTreeNode: opt.V(nodeName),
				Position:    opt.V(position),
				Rotation:    opt.V(tile.RotationQuat()),
				Template:    c.playData.Scene,
				IsDynamic:   false,
			})
		}
	}

	c.scene.SubscribeFixedUpdate(c.onFixedUpdate)

	c.hierScene = c.scene.Hierarchy()
	c.gfxScene = c.scene.Graphics()
	c.physicsScene = c.scene.Physics()
	c.ecsScene = c.scene.ECS()

	c.physicsScene.CreateGlobalAccelerator(acceleration.NewGravityDirection())

	c.vehicleDefinition = c.createVehicleDefinition()

	c.followCameraSystem = preset.NewFollowCameraSystem(c.ecsScene, c.window)
	c.followCameraSystem.UseDefaults()

	c.carSystem = preset.NewCarSystem(c.scene)

	carModel := c.scene.InstantiateModel(game.ModelInfo{
		Name:      opt.V("Vehicle"),
		Template:  c.playData.Vehicle,
		IsDynamic: true,
	})
	c.vehicle = c.vehicleDefinition.ApplyToModel(c.scene, preset.CarApplyInfo{
		Model:    carModel,
		Position: dprec.NewVec3(0.0, 0.5, 0.0),
		Rotation: dprec.RotationQuat(dprec.Degrees(90), dprec.BasisYVec3()),
	})

	vehicleNodeComponent := c.vehicle.Entity().RefNodeComponent()
	vehicleNode := vehicleNodeComponent.Node

	vehicleCarComponent := c.vehicle.Entity().RefCarComponent()
	vehicleCarComponent.LightsOn = (environment == data.LightingNight)

	switch controller {
	case data.InputKeyboard:
		c.vehicle.Entity().SetCarKeyboardComponent(preset.CarKeyboardComponent{
			AccelerateKey: ui.KeyCodeArrowUp,
			DecelerateKey: ui.KeyCodeArrowDown,
			TurnLeftKey:   ui.KeyCodeArrowLeft,
			TurnRightKey:  ui.KeyCodeArrowRight,
			ShiftUpKey:    ui.KeyCodeD,
			ShiftDownKey:  ui.KeyCodeR,
			RecoverKey:    ui.KeyCodeLeftShift,

			AccelerationChangeSpeed: 2.0,
			DecelerationChangeSpeed: 4.0,
			SteeringChangeSpeed:     3.0,
			SteeringRestoreSpeed:    6.0,
		})
	case data.InputMouse:
		c.vehicle.Entity().SetCarMouseComponent(preset.CarMouseComponent{
			AccelerationChangeSpeed: 2.0,
			DecelerationChangeSpeed: 4.0,
			Destination:             dprec.ZeroVec3(),
		})
	case data.InputGamepad:
		c.vehicle.Entity().SetCarGamepadComponent(preset.CarGamepadComponent{
			Gamepad: c.window.Gamepads()[0],
		})
	}

	c.followCamera = c.gfxScene.CreateCamera()
	c.followCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	c.followCamera.SetFoV(sprec.Degrees(66))
	c.followCamera.SetAutoExposure(false)
	c.followCamera.SetExposure(1.5)
	c.followCamera.SetAutoFocus(false)
	c.followCamera.SetCascadeDistances([]float32{16.0, 64.0, 256.0})
	c.gfxScene.SetActiveCamera(c.followCamera)

	followCameraNode := c.hierScene.Wrap(c.hierScene.CreateNode())
	followCameraNode.SetPosition(dprec.NewVec3(0.0, 20.0, 10.0))
	c.scene.CameraBindingSet().Bind(followCameraNode.ID(), c.followCamera)

	c.bonnetCamera = c.gfxScene.CreateCamera()
	c.bonnetCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	c.bonnetCamera.SetFoV(sprec.Degrees(80))
	c.bonnetCamera.SetAutoExposure(false)
	c.bonnetCamera.SetExposure(1.5)
	c.bonnetCamera.SetAutoFocus(false)
	c.bonnetCamera.SetCascadeDistances([]float32{16.0, 64.0, 256.0})

	bonnetCameraNode := c.hierScene.Wrap(c.hierScene.CreateNode())
	c.scene.CameraBindingSet().Bind(bonnetCameraNode.ID(), c.bonnetCamera)
	bonnetCameraNode.SetRotation(dprec.RotationQuat(dprec.Degrees(180), dprec.BasisYVec3()))
	bonnetCameraNode.SetPosition(dprec.NewVec3(0.0, 0.75, 0.35))
	vehicleNode.AppendChild(bonnetCameraNode, false)

	followCameraEntity := c.ecsScene.CreateEntity()
	followCameraEntity.SetNodeComponent(preset.NodeComponent{
		Node: followCameraNode,
	})
	followCameraEntity.SetControlledComponent(preset.ControlledComponent{
		Inputs: preset.ControlInputKeyboard | preset.ControlInputMouse | preset.ControlInputGamepad0,
	})
	followCameraEntity.SetFollowCameraComponent(preset.FollowCameraComponent{
		Target:         vehicleNode,
		AnchorPosition: dprec.Vec3Sum(vehicleNode.Position(), dprec.NewVec3(0.0, 0.0, -anchorDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
		PitchAngle:     dprec.Degrees(-pitchAngle),
		YawAngle:       dprec.Degrees(0),
		Zoom:           1.0,
	})

	runtime.GC()
	c.engine.ResetDeltaTime()
}

func (c *PlayController) Stop() {
	c.scene.Delete()
}

func (c *PlayController) Pause() {
	c.scene.Freeze()
}

func (c *PlayController) Resume() {
	c.scene.Unfreeze()
}

func (c *PlayController) Velocity() float64 {
	if c.vehicle == nil {
		return 0.0
	}
	return c.vehicle.Velocity()
}

func (c *PlayController) ToggleCamera() {
	if c.scene.Graphics().ActiveCamera() == c.followCamera {
		c.scene.Graphics().SetActiveCamera(c.bonnetCamera)
	} else {
		c.scene.Graphics().SetActiveCamera(c.followCamera)
	}
}

func (c *PlayController) IsDrive() bool {
	if c.vehicle == nil {
		return true
	}
	carComp := c.vehicle.Entity().RefCarComponent()
	return carComp.Gear == preset.CarGearForward
}

func (c *PlayController) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.carSystem.OnMouseEvent(element, event)
}

func (c *PlayController) OnKeyboardEvent(event ui.KeyboardEvent) bool {
	return c.carSystem.OnKeyboardEvent(event)
}

func (c *PlayController) onFixedUpdate(elapsedTime time.Duration) {
	c.carSystem.OnFixedUpdate(elapsedTime.Seconds())
	c.followCameraSystem.OnFixedUpdate(elapsedTime.Seconds())
}

func (c *PlayController) createVehicleDefinition() *preset.CarDefinition {
	collisionGroup := physics.NewCollisionGroup()

	chassisBodyDef := c.physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   260,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(208),
		DragFactor:             0.0,
		AngularDragFactor:      0.0,
		RestitutionCoefficient: 0.0,
		CollisionGroup:         collisionGroup,
		CollisionBoxes: []shape3d.Box{
			shape3d.NewBox(
				dprec.NewVec3(0.0, 0.34, -0.3),
				dprec.IdentityQuat(),
				dprec.NewVec3(1.6, 1.18, 3.64),
			),
		},
	})

	wheelBodyDef := c.physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   20,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(0.9),
		DragFactor:             0.0,
		AngularDragFactor:      0.0,
		RestitutionCoefficient: 0.0,
		CollisionGroup:         collisionGroup,
		CollisionSpheres: []shape3d.Sphere{
			shape3d.NewSphere(dprec.ZeroVec3(), 0.25),
		},
	})

	hubBodyDef := c.physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   1,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(0.01),
		DragFactor:             0.0,
		AngularDragFactor:      0.0,
		RestitutionCoefficient: 0.0,
	})

	chassisDef := preset.NewChassisDefinition().
		WithNodeName("Chassis").
		WithBodyDefinition(chassisBodyDef).
		WithHeadLightNodeNames("FLLight", "FRLight").
		WithTailLightNodeNames("BLLight", "BRLight").
		WithBeamLightNodeNames("FLBeamLight", "FRBeamLight").
		WithStopLightNodeNames("BLStopLight", "BRStopLight")

	frontLeftWheelDef := preset.NewWheelDefinition().
		WithNodeName("FLWheel").
		WithBodyDefinition(wheelBodyDef)

	frontRightWheelDef := preset.NewWheelDefinition().
		WithNodeName("FRWheel").
		WithBodyDefinition(wheelBodyDef)

	rearLeftWheelDef := preset.NewWheelDefinition().
		WithNodeName("BLWheel").
		WithBodyDefinition(wheelBodyDef)

	rearRightWheelDef := preset.NewWheelDefinition().
		WithNodeName("BRWheel").
		WithBodyDefinition(wheelBodyDef)

	frontLeftHubDef := preset.NewHubDefinition().
		WithNodeName("FLHub").
		WithBodyDefinition(hubBodyDef)

	frontRightHubDef := preset.NewHubDefinition().
		WithNodeName("FRHub").
		WithBodyDefinition(hubBodyDef)

	rearLeftHubDef := preset.NewHubDefinition().
		WithNodeName("BLHub").
		WithBodyDefinition(hubBodyDef)

	rearRightHubDef := preset.NewHubDefinition().
		WithNodeName("BRHub").
		WithBodyDefinition(hubBodyDef)

	frontAxisDef := preset.NewAxisDefinition().
		WithPosition(dprec.NewVec3(0.0, -0.18, 0.96)).
		WithWidth(1.7).
		WithSuspensionLength(0.16).
		WithSpringLength(0.25).
		WithSpringFrequency(2.9).
		WithSpringDamping(0.8).
		WithLeftWheelDefinition(frontLeftWheelDef).
		WithRightWheelDefinition(frontRightWheelDef).
		WithLeftHubDefinition(frontLeftHubDef).
		WithRightHubDefinition(frontRightHubDef).
		WithMaxSteeringAngle(dprec.Degrees(45)).
		WithMaxAcceleration(145 * 1.4).
		WithMaxBraking(250 * 1.5).
		WithReverseRatio(0.5)

	rearAxisDef := preset.NewAxisDefinition().
		WithPosition(dprec.NewVec3(0.0, -0.18, -1.37)).
		WithWidth(1.7).
		WithSuspensionLength(0.16).
		WithSpringLength(0.25).
		WithSpringFrequency(2.7).
		WithSpringDamping(0.8).
		WithLeftWheelDefinition(rearLeftWheelDef).
		WithRightWheelDefinition(rearRightWheelDef).
		WithLeftHubDefinition(rearLeftHubDef).
		WithRightHubDefinition(rearRightHubDef).
		WithMaxSteeringAngle(dprec.Degrees(0)).
		WithMaxAcceleration(145 * 1.4).
		WithMaxBraking(180 * 3.0).
		WithReverseRatio(0.5)

	carDef := preset.NewCarDefinition().
		WithChassisDefinition(chassisDef).
		WithAxisDefinition(frontAxisDef).
		WithAxisDefinition(rearAxisDef)

	return carDef
}
