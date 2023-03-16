package play

import (
	"fmt"
	"time"

	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/preset"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/metrics"
	"github.com/mokiat/lacking/util/shape"
	"github.com/mokiat/rally-mka/internal/global"
	"github.com/mokiat/rally-mka/internal/scene"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

const (
	anchorDistance = 6.0
	cameraDistance = 10.0
)

type ViewData struct {
	GameData *scene.Data
}

var View = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		context = co.GetContext[global.Context]()
		data    = co.GetData[ViewData](props)
	)

	regionsVisible := co.UseState(func() bool {
		return false
	})

	regionsState := co.UseState(func() []metrics.RegionStat {
		return nil
	})

	frames := co.UseState(func() *int {
		return gog.PtrOf(0)
	}).Get()

	lifecycle := co.UseState(func() *playLifecycle {
		return &playLifecycle{
			window:       co.Window(scope),
			engine:       context.Engine,
			gameData:     data.GameData,
			spansVisible: regionsVisible,
			frames:       frames,
		}
	}).Get()

	speedState := co.UseState(func() float64 {
		return float64(0.0)
	})

	co.Once(func() {
		co.Window(scope).SetCursorVisible(false)
	})

	co.Defer(func() {
		co.Window(scope).SetCursorVisible(true)
	})

	co.Once(func() {
		var refreshSpeed func()
		refreshSpeed = func() {
			var velocity dprec.Vec3
			if car := lifecycle.car; car != nil {
				velocity = car.Body().Velocity()
			}
			speedState.Set(velocity.Length() * 3.6)
			co.After(100*time.Millisecond, refreshSpeed)
		}
		refreshSpeed()
	})

	co.Once(func() {
		lastTime := time.Now()
		var refreshFPS func()
		refreshFPS = func() {
			currentTime := time.Now()
			deltaTime := currentTime.Sub(lastTime)
			if deltaTime > 2*time.Second {
				fps := float64(*frames) / deltaTime.Seconds()
				*frames = 0

				co.Window(scope).SetTitle(fmt.Sprintf("FPS: %.1f", fps))
				lastTime = currentTime
			}
			co.After(time.Second, refreshFPS)
		}
		refreshFPS()
	})

	co.Once(func() {
		var refreshSpans func()
		refreshSpans = func() {
			regionsState.Set(metrics.RegionStats())
			co.After(time.Second, refreshSpans)
		}
		refreshSpans()
	})

	co.Once(func() {
		lifecycle.onCreate()
	})

	co.Defer(func() {
		lifecycle.onDestroy()
	})

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   lifecycle,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Padding: ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    5,
				Bottom: 5,
			},
			Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
				ContentAlignment: mat.AlignmentLeft,
				ContentSpacing:   10,
			}),
		})

		if regionsVisible.Get() {
			co.WithChild("regions", co.New(widget.RegionBlock, func() {
				co.WithData(widget.RegionBlockData{
					Regions: regionsState.Get(),
				})
				co.WithLayoutData(mat.LayoutData{
					GrowHorizontally: true,
				})
			}))
		}

		co.WithChild("speed-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(24.0)),
				FontColor: opt.V(ui.White()),
				Text:      fmt.Sprintf("Speed: %.4f", speedState.Get()),
			})
		}))
	})
})

var _ ui.ElementKeyboardHandler = (*playLifecycle)(nil)
var _ ui.ElementMouseHandler = (*playLifecycle)(nil)

type playLifecycle struct {
	window       *ui.Window
	engine       *game.Engine
	gameData     *scene.Data
	spansVisible *co.State[bool]
	frames       *int

	scene        *game.Scene
	gfxScene     *graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	carSystem         *preset.CarSystem
	cameraStandSystem *preset.FollowCameraSystem

	orbitCamera  *graphics.Camera
	bonnetCamera *graphics.Camera
	car          *game.Node
}

func (h *playLifecycle) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if event.Code == app.KeyCodeEscape {
		element.Window().Close()
		return true
	}
	if event.Code == ui.KeyCodeTab {
		if event.Type == ui.KeyboardEventTypeKeyUp {
			h.spansVisible.Set(!h.spansVisible.Get())
		}
		return true
	}
	if event.Code == app.KeyCodeX {
		if event.Type == app.KeyboardEventTypeKeyDown {
			if h.gfxScene.ActiveCamera() == h.orbitCamera {
				h.gfxScene.SetActiveCamera(h.bonnetCamera)
			} else {
				h.gfxScene.SetActiveCamera(h.orbitCamera)
			}
		}
		return true
	}
	return h.carSystem.OnKeyboardEvent(event) || h.cameraStandSystem.OnKeyboardEvent(event)
}

func (h *playLifecycle) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	// bounds := element.Bounds()
	// viewport := graphics.NewViewport(bounds.X, bounds.Y, bounds.Width, bounds.Height)
	// camera := h.gfxScene.ActiveCamera()
	// return h.vehicleSystem.OnMouseEvent(event, viewport, camera, h.gfxScene)
	return false
}

func (h *playLifecycle) onCreate() {
	h.scene = h.engine.CreateScene()

	h.physicsScene = h.scene.Physics()
	h.gfxScene = h.scene.Graphics()
	h.ecsScene = h.scene.ECS()

	h.carSystem = preset.NewCarSystem(h.ecsScene, h.window)
	h.carSystem.UseDefaults()
	h.cameraStandSystem = preset.NewFollowCameraSystem(h.ecsScene, h.window)
	h.cameraStandSystem.UseDefaults()

	h.orbitCamera = h.gfxScene.CreateCamera()
	h.orbitCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	h.orbitCamera.SetFoV(sprec.Degrees(66))
	h.orbitCamera.SetAutoExposure(true)
	h.orbitCamera.SetExposure(1.0)
	h.orbitCamera.SetAutoFocus(false)
	h.gfxScene.SetActiveCamera(h.orbitCamera)

	h.bonnetCamera = h.gfxScene.CreateCamera()
	h.bonnetCamera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	h.bonnetCamera.SetFoV(sprec.Degrees(80))
	h.bonnetCamera.SetAutoExposure(true)
	h.bonnetCamera.SetExposure(1.0)
	h.bonnetCamera.SetAutoFocus(false)

	h.engine.SubscribePreUpdate(h.onPreUpdate)
	h.engine.SubscribePostUpdate(h.onPostUpdate)
	h.setupLevel(h.gameData)
}

func (h *playLifecycle) onDestroy() {
	h.gameData.Dismiss()
}

func (h *playLifecycle) onPreUpdate(engine *game.Engine, scene *game.Scene, elapsedSeconds float64) {
	h.carSystem.Update(elapsedSeconds)
}

func (h *playLifecycle) onPostUpdate(engine *game.Engine, scene *game.Scene, elapsedSeconds float64) {
	h.cameraStandSystem.Update(elapsedSeconds)
	*h.frames += 1
}

func (h *playLifecycle) setupLevel(gameData *scene.Data) {
	h.scene.Initialize(gameData.Level)

	sunLight := h.gfxScene.CreateDirectionalLight(graphics.DirectionalLightInfo{
		EmitColor: dprec.NewVec3(0.5, 0.5, 0.3),
		EmitRange: 16000, // FIXME
	})

	// FIXME: Should not be needed
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
	lightNode.SetDirectionalLight(sunLight)

	carInstance := h.scene.CreateModel(game.ModelInfo{
		Name:       "SUV",
		Definition: gameData.CarModel,
		Position:   dprec.NewVec3(0.0, 3.0, 0.0),
		Rotation:   dprec.IdentityQuat(),
		Scale:      dprec.NewVec3(1.0, 1.0, 1.0),
		IsDynamic:  true,
	})
	carChassisNode := h.setupCarDemo(
		carInstance,
		dprec.NewVec3(0.0, 2.0, 0.0),
		dprec.IdentityQuat(),
		true,
	)

	h.car = carChassisNode
	h.car.AppendChild(lightNode)

	cameraDeltaNode := game.NewNode()
	cameraDeltaNode.SetCamera(h.bonnetCamera)
	cameraDeltaNode.SetRotation(dprec.RotationQuat(dprec.Degrees(180), dprec.BasisYVec3()))
	cameraDeltaNode.SetPosition(dprec.NewVec3(0.0, 1.0, 0.0))
	carChassisNode.AppendChild(cameraDeltaNode)

	orbitCameraNode := game.NewNode()
	orbitCameraNode.SetAttachable(h.orbitCamera)
	h.scene.Root().AppendChild(orbitCameraNode)

	camTarget := h.car
	cameraEntity := h.ecsScene.CreateEntity()
	ecs.AttachComponent(cameraEntity, &preset.NodeComponent{
		Node: orbitCameraNode,
	})
	ecs.AttachComponent(cameraEntity, &preset.FollowCameraComponent{
		Target:         camTarget,
		AnchorPosition: dprec.Vec3Sum(camTarget.Position(), dprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
		PitchAngle:     dprec.Degrees(-20),
		YawAngle:       dprec.Degrees(0),
		Zoom:           1.0,
	})
	ecs.AttachComponent(cameraEntity, &preset.ControlledComponent{
		Inputs: preset.ControlInputGamepad0 | preset.ControlInputKeyboard,
	})

	h.engine.ResetDeltaTime()
}

func (h *playLifecycle) setupCarDemo(model *game.Model, position dprec.Vec3, rotation dprec.Quat, controlled bool) *game.Node {
	// TODO: The definition can be instantiated just once and the
	// ApplyToModel can be used multiple times on new instances.

	collisionGroup := physics.NewCollisionGroup()

	chassisBodyDef := h.physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   260,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(208),
		DragFactor:             0.0,
		AngularDragFactor:      0.0,
		RestitutionCoefficient: 0.0,
		CollisionGroup:         collisionGroup,
		CollisionShapes: []physics.CollisionShape{
			shape.NewPlacement[shape.Shape](
				shape.NewTransform(
					dprec.NewVec3(0.0, 0.3, -0.4),
					dprec.IdentityQuat(),
				),
				shape.NewStaticBox(1.6, 1.4, 4.0),
			),
		},
	})

	wheelBodyDef := h.physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   20,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(0.9),
		DragFactor:             0.0,
		AngularDragFactor:      0.0,
		RestitutionCoefficient: 0.0,
		CollisionGroup:         collisionGroup,
		CollisionShapes: []physics.CollisionShape{
			shape.NewPlacement[shape.Shape](
				shape.IdentityTransform(),
				shape.NewStaticSphere(0.25),
			),
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
		WithMaxSteeringAngle(dprec.Degrees(30)).
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

	info := preset.CarApplyInfo{
		Model:    model,
		Position: position,
		Rotation: rotation,
	}
	if controlled {
		info.Inputs = preset.ControlInputKeyboard | preset.ControlInputMouse | preset.ControlInputGamepad0
	}
	car := carDef.ApplyToModel(h.scene, info)
	return car.Chassis().Node()
}
