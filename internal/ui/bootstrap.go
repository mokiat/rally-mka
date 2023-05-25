package internal

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/view"
)

func BootstrapApplication(window *ui.Window, gameController *game.Controller) {
	engine := gameController.Engine()

	scope := co.RootScope(window)
	scope = co.TypedValueScope(scope, global.Context{
		Engine:      engine,
		ResourceSet: engine.CreateResourceSet(),
	})
	co.Initialize(scope, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(&bootstrapComponent{})

type bootstrapComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	appModel      *model.Application
	appController *controller.Application
	childrenScope co.Scope
}

func (c *bootstrapComponent) OnCreate() {
	c.appModel = model.NewApplication()
	c.appController = controller.NewApplication(c.appModel)
	c.childrenScope = mvc.UseReducer(c.Scope, c.appController)
}

func (c *bootstrapComponent) Render() co.Instance {
	return co.New(view.Application, func() {
		co.WithData(c.appModel)
		co.WithScope(c.childrenScope)
	})
}
