package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/theme"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

type LoadingScreenData struct {
	Model *model.Loading
}

var LoadingScreen = co.Define(&loadingScreenComponent{})

type loadingScreenComponent struct {
	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`
}

func (c *loadingScreenComponent) OnCreate() {
	screenData := co.GetData[LoadingScreenData](c.Properties)
	loadingModel := screenData.Model
	loadingModel.Promise().OnReady(func() {
		// TODO: Handle errors!
		mvc.Dispatch(c.Scope, action.ChangeView{
			ViewName: loadingModel.NextViewName(),
		})
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
				Font:      co.OpenFont(c.Scope, "ui:///roboto-italic.ttf"),
				FontSize:  opt.V(float32(32)),
				FontColor: opt.V(theme.PrimaryColor),
				Text:      "Loading...",
			})
		}))
	})
}
