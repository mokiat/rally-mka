package game

import (
	"time"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/internal/ecssys"
)

func NewController(reg asset.Registry, gfxEngine *graphics.Engine) *Controller {
	controller := &Controller{
		gfxEngine:     gfxEngine,
		physicsEngine: physics.NewEngine(),
		ecsEngine:     ecs.NewEngine(),

		lastFrameTime: time.Now(),
	}
	controller.registry = resource.NewRegistry(reg, gfxEngine, controller)
	return controller
}

type Controller struct {
	app.NopController

	window        app.Window
	gfxEngine     *graphics.Engine
	physicsEngine *physics.Engine
	ecsEngine     *ecs.Engine
	registry      *resource.Registry

	lastFrameTime time.Time
	freezeFrame   bool

	width  int
	height int

	gfxScene     *graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	renderSystem      *ecssys.Renderer
	vehicleSystem     *ecssys.VehicleSystem
	cameraStandSystem *ecssys.CameraStandSystem

	camera *graphics.Camera

	OnUpdate func()
}

func (c *Controller) Schedule(fn func() error) {
	c.window.Schedule(fn)
}

func (c *Controller) Registry() *resource.Registry {
	return c.registry
}

func (c *Controller) GFXEngine() *graphics.Engine {
	return c.gfxEngine
}

func (c *Controller) GFXScene() *graphics.Scene {
	return c.gfxScene
}

func (c *Controller) PhysicsScene() *physics.Scene {
	return c.physicsScene
}

func (c *Controller) ECSScene() *ecs.Scene {
	return c.ecsScene
}

func (c *Controller) RenderSystem() *ecssys.Renderer {
	return c.renderSystem
}

func (c *Controller) VehicleSystem() *ecssys.VehicleSystem {
	return c.vehicleSystem
}

func (c *Controller) CameraStandSystem() *ecssys.CameraStandSystem {
	return c.cameraStandSystem
}

func (c *Controller) Camera() *graphics.Camera {
	return c.camera
}

func (c *Controller) OnCreate(window app.Window) {
	c.window = window
	c.width, c.height = window.Size()

	c.gfxEngine.Create()

	c.gfxScene = c.gfxEngine.CreateScene()
	c.physicsScene = c.physicsEngine.CreateScene(0.015)
	c.ecsScene = c.ecsEngine.CreateScene()

	c.camera = c.gfxScene.CreateCamera()

	c.renderSystem = ecssys.NewRenderer(c.ecsScene)
	c.vehicleSystem = ecssys.NewVehicleSystem(c.ecsScene)
	c.cameraStandSystem = ecssys.NewCameraStandSystem(c.ecsScene)
}

func (c *Controller) OnResize(window app.Window, width, height int) {
	c.width, c.height = width, height
}

func (c *Controller) OnCloseRequested(window app.Window) {
	window.Close()
}

func (c *Controller) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	if event.Code == app.KeyCodeEscape {
		window.Close()
		return true
	}
	if event.Code == app.KeyCodeF {
		switch event.Type {
		case app.KeyboardEventTypeKeyDown:
			c.freezeFrame = true
			return true
		case app.KeyboardEventTypeKeyUp:
			c.freezeFrame = false
			return true
		}
	}
	return c.vehicleSystem.OnKeyboardEvent(event)
}

func (c *Controller) OnRender(window app.Window) {
	currentTime := time.Now()
	elapsedSeconds := float32(currentTime.Sub(c.lastFrameTime).Seconds())
	c.lastFrameTime = currentTime

	if !c.freezeFrame {
		var gamepad *app.GamepadState
		if state, ok := window.GamepadState(0); ok {
			gamepad = &state
		}

		c.physicsScene.Update(elapsedSeconds)
		c.vehicleSystem.Update(elapsedSeconds, gamepad)
		c.renderSystem.Update()
		c.cameraStandSystem.Update(elapsedSeconds, gamepad)

		if c.OnUpdate != nil {
			c.OnUpdate()
		}

		c.gfxScene.Render(graphics.NewViewport(0, 0, c.width, c.height), c.camera)
	}

	window.Invalidate() // force redraw
}

func (c *Controller) OnDestroy(window app.Window) {
	c.renderSystem = nil
	c.vehicleSystem = nil
	c.cameraStandSystem = nil

	c.ecsScene.Delete()
	c.physicsScene.Delete()
	c.gfxScene.Delete()

	c.gfxEngine.Destroy()
}
