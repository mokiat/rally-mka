package internal

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/view"
)

func BootstrapApplication(window *ui.Window, gameController *game.Controller) {
	engine := gameController.Engine()
	eventBus := mvc.NewEventBus()

	scope := co.RootScope(window)
	scope = co.TypedValueScope(scope, eventBus)
	scope = co.TypedValueScope(scope, global.Context{
		Engine:      engine,
		ResourceSet: engine.CreateResourceSet(),
	})
	co.Initialize(scope, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(&bootstrapComponent{})

type bootstrapComponent struct {
	co.BaseComponent
}

func (c *bootstrapComponent) Render() co.Instance {
	return co.New(view.Application, nil)
}
