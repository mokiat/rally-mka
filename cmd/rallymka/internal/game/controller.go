package game

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/input"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/loading"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/simulation"
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
	activeView     View
	loadingView    View
	simulationView View
}

func (c *Controller) Init(ctx game.InitContext) error {
	locator := resource.FileLocator{}
	registry := resource.NewRegistry(locator, ctx.GFXWorker)

	c.loadingView = loading.NewView(registry)
	c.simulationView = simulation.NewView(registry, ctx.GFXWorker)

	c.loadingView.Load()
	c.simulationView.Load()
	return nil
}

func (c *Controller) Update(ctx game.UpdateContext) bool {
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
	if c.activeView != nil {
		c.activeView.Close()
		c.activeView.Unload()
	}
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
