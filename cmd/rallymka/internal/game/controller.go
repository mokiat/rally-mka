package game

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/loading"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/simulation"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const (
	maxQueuedResources = 64
	maxResources       = 1024
	maxEvents          = 64
)

type View interface {
	Load()
	IsAvailable() bool
	Unload()

	Open()
	Close()

	Update(ctx game.UpdateContext)
	Render(ctx game.RenderContext)
}

func NewController() *Controller {
	return &Controller{}
}

type Controller struct {
	resWorker   *resource.Worker
	resRegistry *resource.Registry

	windowSize     game.WindowSize
	activeView     View
	loadingView    View
	simulationView View
}

func (c *Controller) Init(ctx game.InitContext) error {
	c.resWorker = resource.NewWorker(maxQueuedResources)
	c.resRegistry = resource.NewRegistry(c.resWorker, maxResources, maxEvents)

	locator := resource.FileLocator{}
	programOperator := stream.NewProgramOperator(locator, ctx.GFXWorker)
	programOperator.Register(c.resRegistry)
	cubeTextureOperator := stream.NewCubeTextureOperator(locator, ctx.GFXWorker)
	cubeTextureOperator.Register(c.resRegistry)
	twodTextureOperator := stream.NewTwoDTextureOperator(locator, ctx.GFXWorker)
	twodTextureOperator.Register(c.resRegistry)
	modelOperator := stream.NewModelOperator(locator, ctx.GFXWorker)
	modelOperator.Register(c.resRegistry)
	meshOperator := stream.NewMeshOperator(locator, ctx.GFXWorker)
	meshOperator.Register(c.resRegistry)
	levelOperator := stream.NewLevelOperator(locator, ctx.GFXWorker)
	levelOperator.Register(c.resRegistry)

	c.loadingView = loading.NewView(c.resRegistry)
	c.simulationView = simulation.NewView(c.resRegistry, ctx.GFXWorker)

	go c.resWorker.Work()
	c.loadingView.Load()
	c.simulationView.Load()
	return nil
}

func (c *Controller) Update(ctx game.UpdateContext) bool {
	c.resRegistry.Update()
	c.pickView()

	if c.activeView != nil {
		c.activeView.Update(ctx)
	}

	return !ctx.Keyboard.IsPressed(input.KeyEscape)
}

func (c *Controller) Render(ctx game.RenderContext) {
	if c.activeView != nil {
		c.activeView.Render(ctx)
	}
}

func (c *Controller) Release(ctx game.ReleaseContext) error {
	return nil
}

func (c *Controller) pickView() {
	switch c.activeView {
	case nil:
		if c.loadingView.IsAvailable() {
			c.changeView(c.loadingView)
		}
	case c.loadingView:
		if c.simulationView.IsAvailable() {
			c.changeView(c.simulationView)
		}
	}
}

func (c *Controller) changeView(view View) {
	if c.activeView != nil {
		c.activeView.Close()
	}
	c.activeView = view
	if c.activeView != nil {
		c.activeView.Open()
	}
}
