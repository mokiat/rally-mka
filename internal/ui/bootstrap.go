package internal

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game"
	"github.com/mokiat/rally-mka/internal/global"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/view"
)

func BootstrapApplication(window *ui.Window, gameController *game.Controller) {
	co.RegisterContext(global.Context{
		GameController: gameController,
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
