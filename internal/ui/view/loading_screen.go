package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/theme"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

type LoadingScreenData struct {
	Model *model.Loading
}

var LoadingScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		screenData   = co.GetData[LoadingScreenData](props)
		loadingModel = screenData.Model
	)

	co.Once(func() {
		loadingModel.Promise().OnReady(func() {
			// TODO: Handle errors!

			mvc.Dispatch(scope, action.ChangeView{
				ViewName: loadingModel.NextViewName(),
			})
		})
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("loading", co.New(widget.Loading, func() {
			co.WithLayoutData(mat.LayoutData{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
		}))

		co.WithChild("info-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-italic.ttf"),
				FontSize:  opt.V(float32(32)),
				FontColor: opt.V(theme.PrimaryColor),
				Text:      "Loading...",
			})
			co.WithLayoutData(mat.LayoutData{
				Right:  opt.V(40),
				Bottom: opt.V(40),
			})
		}))
	})
})
