package game

import (
	"time"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/internal/ecssys"
)

func NewController(reg asset.Registry, gfxEngine *graphics.Engine) *Controller {
	engine := game.NewEngine(
		game.WithPhysics(physics.NewEngine()),
		game.WithGraphics(gfxEngine),
		game.WithECS(ecs.NewEngine()),
	)

	controller := &Controller{
		engine: engine,

		lastFrameTime: time.Now(),
	}
	controller.registry = resource.NewRegistry(reg, gfxEngine, controller)
	return controller
}

type Controller struct {
	app.NopController

	window   app.Window
	engine   *game.Engine
	scene    *game.Scene
	registry *resource.Registry

	lastFrameTime time.Time
	freezeFrame   bool

	width  int
	height int

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

func (c *Controller) Engine() *game.Engine {
	return c.engine
}

func (c *Controller) Scene() *game.Scene {
	return c.scene
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

	c.engine.Graphics().Create()
	c.scene = c.engine.CreateScene()
	c.camera = c.scene.Graphics().CreateCamera()

	c.vehicleSystem = ecssys.NewVehicleSystem(c.scene.ECS())
	c.cameraStandSystem = ecssys.NewCameraStandSystem(c.scene.ECS())
}

func (c *Controller) OnDestroy(window app.Window) {
	c.vehicleSystem = nil
	c.cameraStandSystem = nil

	c.scene.Delete()
	c.engine.Graphics().Destroy()
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
	return c.vehicleSystem.OnKeyboardEvent(event) ||
		c.cameraStandSystem.OnKeyboardEvent(event)
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

		c.vehicleSystem.Update(elapsedSeconds, gamepad)
		c.scene.Update(elapsedSeconds)
		c.cameraStandSystem.Update(elapsedSeconds, gamepad)
		if c.OnUpdate != nil {
			c.OnUpdate()
		}

		c.scene.Graphics().Render(graphics.NewViewport(0, 0, c.width, c.height), c.camera)
	}

	window.Invalidate() // force redraw
}
