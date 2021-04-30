package game

import (
	"time"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/loading"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/simulation"
)

type View interface {
	Load(window app.Window, cb func())
	Unload(window app.Window)

	Open(window app.Window)
	Close(window app.Window)

	OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool

	Update(window app.Window, elapsedTime time.Duration)
	Render(window app.Window, width, height int, pipeline *graphics.Pipeline)
}

func NewController() *Controller {
	gfxWorker := async.NewWorker(1024)
	gfxRenderer := graphics.NewRenderer()
	locator := resource.FileLocator{}
	registry := resource.NewRegistry(locator, gfxWorker)

	return &Controller{
		gfxWorker:   gfxWorker,
		gfxRenderer: gfxRenderer,
		registry:    registry,

		loadingView:    loading.NewView(registry),
		simulationView: simulation.NewView(registry, gfxWorker),

		lastFrameTime: time.Now(),
	}
}

type Controller struct {
	app.NopController

	window      app.Window
	gfxWorker   *async.Worker
	gfxRenderer *graphics.Renderer
	registry    *resource.Registry

	activeView     View
	loadingView    View
	simulationView View

	lastFrameTime time.Time
	width         int
	height        int
}

func (c *Controller) OnCreate(window app.Window) {
	c.window = window
	c.width, c.height = window.Size()

	c.loadingView.Load(window, c.onLoadingAvailable)
	c.simulationView.Load(window, c.onSimulationAvailable)
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
	if c.activeView != nil {
		return c.activeView.OnKeyboardEvent(window, event)
	}
	return false
}

func (c *Controller) OnRender(window app.Window) {
	c.gfxWorker.ProcessTryMultiple(10)

	currentTime := time.Now()
	elapsedTime := currentTime.Sub(c.lastFrameTime)
	c.lastFrameTime = currentTime

	if c.activeView != nil {
		c.activeView.Update(window, elapsedTime)
		pipeline := c.gfxRenderer.BeginPipeline()
		c.activeView.Render(window, c.width, c.height, pipeline)
		c.gfxRenderer.EndPipeline(pipeline)
		c.gfxRenderer.Render()
	}

	window.Invalidate() // force redraw
}

func (c *Controller) OnDestroy(window app.Window) {
	c.changeView(nil)

	c.loadingView.Unload(window)
	c.simulationView.Unload(window)
}

func (c *Controller) onLoadingAvailable() {
	c.changeView(c.loadingView)
}

func (c *Controller) onSimulationAvailable() {
	c.changeView(c.simulationView)
}

func (c *Controller) changeView(view View) {
	if c.activeView != nil {
		c.activeView.Close(c.window)
	}
	c.activeView = view
	if c.activeView != nil {
		c.activeView.Open(c.window)
	}
}
