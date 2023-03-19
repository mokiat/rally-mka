package view

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

type IntroScreenData struct {
	Home         *model.Home
	LoadingModel *model.Loading
}

var IntroScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		globalContext = co.GetContext[global.Context]()
		screenData    = co.GetData[IntroScreenData](props)

		engine       = globalContext.Engine
		homeModel    = screenData.Home
		loadingModel = screenData.LoadingModel
	)

	co.Once(func() {
		co.Window(scope).SetCursorVisible(false)
	})
	co.Defer(func() {
		co.Window(scope).SetCursorVisible(true)
	})

	co.Once(func() {
		resourceSet := engine.CreateResourceSet()
		homeModel.SetData(data.LoadHomeData(engine, resourceSet))
	})

	co.After(time.Second, func() {
		promise := homeModel.Data()
		if promise.Ready() {
			// TODO: Handle errors!!!
			mvc.Dispatch(scope, action.ChangeView{
				ViewName: model.ViewNameHome,
			})
		} else {
			loadingModel.SetPromise(promise)
			loadingModel.SetNextViewName(model.ViewNameHome)
			mvc.Dispatch(scope, action.ChangeView{
				ViewName: model.ViewNameLoading,
			})
		}
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("logo-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				BackgroundColor: opt.V(ui.Transparent()),
				Image:           co.OpenImage(scope, "ui/images/logo.png"),
				Mode:            mat.ImageModeFit,
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            opt.V(512),
				Height:           opt.V(128),
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
		}))
	})
})
