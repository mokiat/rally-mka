package view

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

type IntroScreenData struct {
	Home         *model.Home
	LoadingModel *model.Loading
}

var IntroScreen = co.Define(&introScreenComponent{})

type introScreenComponent struct {
	co.BaseComponent
}

func (c *introScreenComponent) OnCreate() {
	co.Window(c.Scope()).SetCursorVisible(false)

	globalContext := co.TypedValue[global.Context](c.Scope())
	engine := globalContext.Engine
	resourceSet := globalContext.ResourceSet

	screenData := co.GetData[IntroScreenData](c.Properties())
	homeModel := screenData.Home
	loadingModel := screenData.LoadingModel

	homeModel.SetData(data.LoadHomeData(engine, resourceSet))

	co.After(c.Scope(), time.Second, func() {
		promise := homeModel.Data()
		if promise.Ready() {
			// TODO: Handle errors!!!
			mvc.Dispatch(c.Scope(), action.ChangeView{
				ViewName: model.ViewNameHome,
			})
		} else {
			loadingModel.SetPromise(promise)
			loadingModel.SetNextViewName(model.ViewNameHome)
			mvc.Dispatch(c.Scope(), action.ChangeView{
				ViewName: model.ViewNameLoading,
			})
		}
	})
}

func (c *introScreenComponent) OnDelete() {
	co.Window(c.Scope()).SetCursorVisible(true)
}

func (c *introScreenComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          layout.Anchor(),
		})

		co.WithChild("logo-picture", co.New(std.Picture, func() {
			co.WithLayoutData(layout.Data{
				Width:            opt.V(512),
				Height:           opt.V(128),
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
			co.WithData(std.PictureData{
				BackgroundColor: opt.V(ui.Transparent()),
				Image:           co.OpenImage(c.Scope(), "ui/images/logo.png"),
				Mode:            std.ImageModeFit,
			})
		}))
	})
}
