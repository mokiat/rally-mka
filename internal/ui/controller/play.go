package controller

import (
	"runtime"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/collision"
	"github.com/mokiat/lacking/game/preset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/rally-mka/internal/game/data"
)

const (
	anchorDistance = 6.0
	cameraDistance = 10.0
	pitchAngle     = 20.0
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

	preUpdateSubscription  *game.UpdateSubscription
	postUpdateSubscription *game.UpdateSubscription

	scene        *game.Scene
	gfxScene     *graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	followCameraSystem *preset.FollowCameraSystem
	followCamera       *graphics.Camera
	bonnetCamera       *graphics.Camera

	carSystem         *preset.CarSystem
	vehicleDefinition *preset.CarDefinition
	vehicle           *preset.Car
}

func (c *PlayController) Start(environment data.Environment, controller data.Controller) {
	// TODO: These subscriptions should be attached on the scene and/or the physics.
	c.preUpdateSubscription = c.engine.SubscribePreUpdate(c.onPreUpdate)
	c.postUpdateSubscription = c.engine.SubscribePostUpdate(c.onPostUpdate)

	c.scene = c.engine.CreateScene()
	c.scene.Initialize(c.playData.Scene)

	c.gfxScene = c.scene.Graphics()
	c.physicsScene = c.scene.Physics()
	c.ecsScene = c.scene.ECS()

	c.vehicleDefinition = c.createVehicleDefinition()

	c.followCameraSystem = preset.NewFollowCameraSystem(c.ecsScene, c.window)
	c.followCameraSystem.UseDefaults()

	c.carSystem = preset.NewCarSystem(c.ecsScene, c.gfxScene, c.window)

	var sunLight *graphics.DirectionalLight
	switch environment {
	case data.EnvironmentDay:
		sunLight = c.scene.Graphics().CreateDirectionalLight(graphics.DirectionalLightInfo{
			EmitColor: dprec.NewVec3(0.5, 0.5, 0.3),
			EmitRange: 16000, // FIXME
		})
	case data.EnvironmentNight:
		sunLight = c.scene.Graphics().CreateDirectionalLight(graphics.DirectionalLightInfo{
			EmitColor: dprec.NewVec3(0.001, 0.001, 0.001),
			EmitRange: 16000, // FIXME
		})
	}

	lightNode := game.NewNode()
	lightNode.SetPosition(dprec.NewVec3(-100.0, 100.0, 0.0))
	lightNode.SetRotation(dprec.QuatProd(
		dprec.RotationQuat(dprec.Degrees(-90), dprec.BasisYVec3()),
		dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
	))
	lightNode.UseTransformation(func(parent, current dprec.Mat4) dprec.Mat4 {
		// Remove parent's rotation
		parent.M11 = 1.0
		parent.M12 = 0.0
		parent.M13 = 0.0
		parent.M21 = 0.0
		parent.M22 = 1.0
		parent.M23 = 0.0
		parent.M31 = 0.0
		parent.M32 = 0.0
		parent.M33 = 1.0
		return dprec.Mat4Prod(parent, current)
	})
	lightNode.SetAttachable(sunLight)

	carInstance := c.scene.CreateModel(game.ModelInfo{
		Name:       "SUV",
		Definition: c.playData.Vehicle,
		Position:   dprec.ZeroVec3(),
		Rotation:   dprec.IdentityQuat(),
		Scale:      dprec.NewVec3(1.0, 1.0, 1.0),
		IsDynamic:  true,
	})
	c.vehicle = c.vehicleDefinition.ApplyToModel(c.scene, preset.CarApplyInfo{
		Model:    carInstance,
		Position: dprec.NewVec3(0.0, 0.5, 0.0),
		Rotation: dprec.IdentityQuat(),
	})
	var vehicleNodeComponent *preset.NodeComponent
	ecs.FetchComponent(c.vehicle.Entity(), &vehicleNodeComponent)
	vehicleNode := vehicleNodeComponent.Node
	vehicleNode.AppendChild(lightNode) // FIXME

	switch controller {
	case data.ControllerKeyboard:
		ecs.AttachComponent(c.vehicle.Entity(), &preset.CarKeyboardControl{
			AccelerateKey: ui.KeyCodeArrowUp,
			DecelerateKey: ui.KeyCodeArrowDown,
			TurnLeftKey:   ui.KeyCodeArrowLeft,
			TurnRightKey:  ui.KeyCodeArrowRight,
			ShiftUpKey:    ui.KeyCodeA,
			ShiftDownKey:  ui.KeyCodeZ,
			RecoverKey:    ui.KeyCodeLeftShift,

			AccelerationChangeSpeed: 1.0,
			DecelerationChangeSpeed: 2.0,
			SteeringChangeSpeed:     3.0,
			SteeringRestoreSpeed:    6.0,
		})
	case data.ControllerMouse:
		ecs.AttachComponent(c.vehicle.Entity(), &preset.CarMouseControl{
			AccelerationChangeSpeed: 1.0,
			DecelerationChangeSpeed: 2.0,
			Destination:             dprec.ZeroVec3(),
		})
	case data.ControllerGamepad:
		ecs.AttachComponent(c.vehicle.Entity(), &preset.CarGamepadControl{
			Gamepad: c.window.Gamepads()[0],
		})
	}

	c.followCamera = c.gfxScene.CreateCamera()
	c.followCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	c.followCamera.SetFoV(sprec.Degrees(66))
	c.followCamera.SetAutoExposure(false)
	c.followCamera.SetExposure(15.0)
	c.followCamera.SetAutoFocus(false)
	c.gfxScene.SetActiveCamera(c.followCamera)

	followCameraNode := game.NewNode()
	followCameraNode.SetPosition(dprec.NewVec3(0.0, 20.0, 10.0))
	followCameraNode.SetAttachable(c.followCamera)
	c.scene.Root().AppendChild(followCameraNode)

	c.bonnetCamera = c.gfxScene.CreateCamera()
	c.bonnetCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	c.bonnetCamera.SetFoV(sprec.Degrees(80))
	c.bonnetCamera.SetAutoExposure(false)
	c.bonnetCamera.SetExposure(15.0)
	c.bonnetCamera.SetAutoFocus(false)

	bonnetCameraNode := game.NewNode()
	bonnetCameraNode.SetAttachable(c.bonnetCamera)
	bonnetCameraNode.SetRotation(dprec.RotationQuat(dprec.Degrees(180), dprec.BasisYVec3()))
	bonnetCameraNode.SetPosition(dprec.NewVec3(0.0, 1.0, 0.0))
	vehicleNode.AppendChild(bonnetCameraNode)

	followCameraEntity := c.ecsScene.CreateEntity()
	ecs.AttachComponent(followCameraEntity, &preset.NodeComponent{
		Node: followCameraNode,
	})
	ecs.AttachComponent(followCameraEntity, &preset.ControlledComponent{
		Inputs: preset.ControlInputKeyboard | preset.ControlInputMouse | preset.ControlInputGamepad0,
	})
	ecs.AttachComponent(followCameraEntity, &preset.FollowCameraComponent{
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
	c.preUpdateSubscription.Delete()
	c.postUpdateSubscription.Delete()
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
	var carComp *preset.CarComponent
	ecs.FetchComponent(c.vehicle.Entity(), &carComp)
	return carComp.Gear == preset.CarGearForward
}

func (c *PlayController) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.carSystem.OnMouseEvent(element, event)
}

func (c *PlayController) OnKeyboardEvent(event ui.KeyboardEvent) bool {
	return c.carSystem.OnKeyboardEvent(event)
}

func (c *PlayController) onPreUpdate(engine *game.Engine, scene *game.Scene, elapsedSeconds float64) {
	// TODO: This check will not be necessary if subscription is on Scene.
	if !c.scene.IsFrozen() {
		c.carSystem.Update(elapsedSeconds)
	}
}

func (c *PlayController) onPostUpdate(engine *game.Engine, scene *game.Scene, elapsedSeconds float64) {
	// TODO: This check will not be necessary if subscription is on Scene.
	if !c.scene.IsFrozen() {
		c.followCameraSystem.Update(elapsedSeconds)
	}
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
		CollisionBoxes: []collision.Box{
			collision.NewBox(
				dprec.NewVec3(0.0, 0.3, -0.4),
				dprec.IdentityQuat(),
				dprec.NewVec3(1.6, 1.4, 4.0),
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
		CollisionSpheres: []collision.Sphere{
			collision.NewSphere(dprec.ZeroVec3(), 0.25),
		},
	})

	chassisDef := preset.NewChassisDefinition().
		WithNodeName("Chassis").
		WithBodyDefinition(chassisBodyDef)

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

	frontAxisDef := preset.NewAxisDefinition().
		WithPosition(dprec.NewVec3(0.0, -0.22, 0.96)).
		WithWidth(1.8).
		WithSuspensionLength(0.23).
		WithSpringLength(0.25).
		WithSpringFrequency(2.9).
		WithSpringDamping(0.8).
		WithLeftWheelDefinition(frontLeftWheelDef).
		WithRightWheelDefinition(frontRightWheelDef).
		WithMaxSteeringAngle(dprec.Degrees(45)).
		WithMaxAcceleration(145).
		WithMaxBraking(250).
		WithReverseRatio(0.5)

	rearAxisDef := preset.NewAxisDefinition().
		WithPosition(dprec.NewVec3(0.0, -0.22, -1.37)).
		WithWidth(1.8).
		WithSuspensionLength(0.23).
		WithSpringLength(0.25).
		WithSpringFrequency(2.4).
		WithSpringDamping(0.8).
		WithLeftWheelDefinition(rearLeftWheelDef).
		WithRightWheelDefinition(rearRightWheelDef).
		WithMaxSteeringAngle(dprec.Degrees(0)).
		WithMaxAcceleration(160).
		WithMaxBraking(180).
		WithReverseRatio(0.5)

	carDef := preset.NewCarDefinition().
		WithChassisDefinition(chassisDef).
		WithAxisDefinition(frontAxisDef).
		WithAxisDefinition(rearAxisDef)

	return carDef
}
