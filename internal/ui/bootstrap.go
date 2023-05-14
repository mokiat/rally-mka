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
	resourceSet := engine.CreateResourceSet()
	co.RegisterContext(global.Context{
		Engine:      engine,
		ResourceSet: resourceSet,
	})
	co.Initialize(window, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	appModel := co.UseState(func() *model.Application {
		return model.NewApplication()
	})
	appController := co.UseState(func() *controller.Application {
		return controller.NewApplication(appModel.Get())
	})
	return co.New(view.Application, func() {
		co.WithData(appModel.Get())
		co.WithScope(mvc.UseReducer(scope, appController.Get()))
	})
})
