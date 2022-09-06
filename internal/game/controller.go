package game

import (
	"time"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/util/metrics"
	"github.com/mokiat/rally-mka/internal/ecssys"
)

func NewController(reg asset.Registry, gfxEngine *graphics.Engine) *Controller {
	ioWorker := async.NewWorker(1)
	go func() {
		ioWorker.ProcessAll()
	}()

	controller := &Controller{
		lastFrameTime: time.Now(),
	}
	controller.engine = game.NewEngine(
		game.WithGFXWorker(game.WorkerFunc(func(fn func() error) game.Operation {
			operation := game.NewOperation()
			controller.Schedule(func() error {
				err := fn()
				operation.Complete(err)
				return err
			})
			return operation
		})),
		game.WithIOWorker(game.WorkerFunc(func(fn func() error) game.Operation {
			operation := game.NewOperation()
			ioWorker.ScheduleFunc(func() error {
				err := fn()
				operation.Complete(err)
				return err
			})
			return operation
		})),
		game.WithRegistry(reg),
		game.WithPhysics(physics.NewEngine()),
		game.WithGraphics(gfxEngine),
		game.WithECS(ecs.NewEngine()),
	)
	return controller
}

type Controller struct {
	app.NopController

	window app.Window
	engine *game.Engine
	scene  *game.Scene

	lastFrameTime time.Time
	freezeFrame   bool

	width  int
	height int

	vehicleSystem     *ecssys.VehicleSystem
	cameraStandSystem *ecssys.CameraStandSystem

	camera *graphics.Camera

	OnUpdate      func()
	OnFrameFinish func()
}

func (c *Controller) Schedule(fn func() error) {
	c.window.Schedule(fn)
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

func (c *Controller) OnMouseEvent(window app.Window, event app.MouseEvent) bool {
	return c.vehicleSystem.OnMouseEvent(event, graphics.NewViewport(0, 0, c.width, c.height), c.camera, c.scene.Graphics())
}

func (c *Controller) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	if event.Code == app.KeyCodeEscape {
		window.Close()
		return true
	}
	if event.Code == app.KeyCodeC && event.Type == app.KeyboardEventTypeKeyDown {
		graphics.ShowLightView = !graphics.ShowLightView
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

func (c *Controller) ResetFrameTime() {
	c.lastFrameTime = time.Now()
}

func (c *Controller) OnRender(window app.Window) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(c.lastFrameTime).Seconds()
	c.lastFrameTime = currentTime

	frameSpan := metrics.BeginSpan("frame")
	defer frameSpan.End()

	if !c.freezeFrame {
		var gamepad *app.GamepadState
		if state, ok := window.GamepadState(0); ok {
			gamepad = &state
		}

		vehSpan := metrics.BeginSpan("vehicle system")
		c.vehicleSystem.Update(elapsedSeconds, gamepad)
		vehSpan.End()
		updateSpan := metrics.BeginSpan("scene update")
		c.scene.Update(elapsedSeconds)
		updateSpan.End()
		cameraSpan := metrics.BeginSpan("camera system")
		c.cameraStandSystem.Update(elapsedSeconds, gamepad)
		cameraSpan.End()
		callbackSpan := metrics.BeginSpan("update callback")
		if c.OnUpdate != nil {
			c.OnUpdate()
		}
		callbackSpan.End()

		renderSpan := metrics.BeginSpan("render")
		c.scene.Graphics().Render(graphics.NewViewport(0, 0, c.width, c.height), c.camera)
		renderSpan.End()
	}

	if c.OnFrameFinish != nil {
		c.OnFrameFinish()
	}

	window.Invalidate() // force redraw
}
