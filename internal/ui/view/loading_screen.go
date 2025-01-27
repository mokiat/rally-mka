package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/theme"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var LoadingScreen = co.Define(&loadingScreenComponent{})

type LoadingScreenData struct {
	AppModel     *model.ApplicationModel
	LoadingModel *model.LoadingModel
}

type loadingScreenComponent struct {
	co.BaseComponent
}

func (c *loadingScreenComponent) OnCreate() {
	screenData := co.GetData[LoadingScreenData](c.Properties())
	appModel := screenData.AppModel
	loadingModel := screenData.LoadingModel

	state := loadingModel.State()
	state.Promise.OnSuccess(func() {
		appModel.SetActiveView(state.SuccessViewName)
	})
	state.Promise.OnError(func() {
		appModel.SetActiveView(state.ErrorViewName)
	})
}

func (c *loadingScreenComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          layout.Anchor(),
		})

		co.WithChild("loading", co.New(widget.Loading, func() {
			co.WithLayoutData(layout.Data{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
		}))

		co.WithChild("info-label", co.New(std.Label, func() {
			co.WithLayoutData(layout.Data{
				Right:  opt.V(40),
				Bottom: opt.V(40),
			})
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-italic.ttf"),
				FontSize:  opt.V(float32(32)),
				FontColor: opt.V(theme.PrimaryColor),
				Text:      "Loading...",
			})
		}))
	})
}
