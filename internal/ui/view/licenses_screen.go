package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
	"github.com/mokiat/rally-mka/resources"
)

var LicensesScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	onBackClicked := func() {
		mvc.Dispatch(scope, action.ChangeView{
			ViewName: model.ViewNameHome,
		})
	}

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("menu-pane", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: opt.V(ui.Black()),
				Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(0),
				Width:  opt.V(200),
			})

			co.WithChild("button", co.New(widget.HomeButton, func() {
				co.WithData(widget.HomeButtonData{ // FIXME: This is not HomeButton
					Text: "Back",
				})
				co.WithLayoutData(mat.LayoutData{
					HorizontalCenter: opt.V(0),
					Bottom:           opt.V(100),
				})
				co.WithCallbackData(widget.HomeButtonCallbackData{
					ClickListener: onBackClicked,
				})
			}))
		}))

		co.WithChild("content-pane", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: opt.V(ui.RGB(0x11, 0x11, 0x11)),
				Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(200),
				Right:  opt.V(0),
			})

			co.WithChild("title", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(32)),
					FontColor: opt.V(ui.White()),
					Text:      "Open-Source Licenses",
				})

				co.WithLayoutData(mat.LayoutData{
					Top:              opt.V(15),
					Height:           opt.V(32),
					HorizontalCenter: opt.V(0),
				})
			}))

			co.WithChild("sub-title", co.New(mat.Label, func() {
				co.WithData(mat.LabelData{
					Font:      co.OpenFont(scope, "mat:///roboto-italic.ttf"),
					FontSize:  opt.V(float32(20)),
					FontColor: opt.V(ui.White()),
					Text:      "- scroll to view all -",
				})

				co.WithLayoutData(mat.LayoutData{
					Top:              opt.V(50),
					Height:           opt.V(20),
					HorizontalCenter: opt.V(0),
				})
			}))

			co.WithChild("license-scroll-pane", co.New(mat.ScrollPane, func() {
				co.WithData(mat.ScrollPaneData{
					DisableHorizontal: true,
					DisableVertical:   false,
				})

				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(80),
					Bottom: opt.V(0),
					Left:   opt.V(0),
					Right:  opt.V(0),
				})

				co.WithChild("license-holder", co.New(mat.Element, func() {
					co.WithData(mat.ElementData{
						Padding: ui.Spacing{
							Top:    100,
							Bottom: 100,
						},
						Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
					})
					co.WithLayoutData(mat.LayoutData{
						GrowHorizontally: true,
					})
					co.WithChild("license-text", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(16)),
							FontColor: opt.V(ui.White()),
							Text:      resources.Licenses,
						})
						co.WithLayoutData(mat.LayoutData{
							HorizontalCenter: opt.V(0),
							VerticalCenter:   opt.V(0),
						})
					}))
				}))
			}))
		}))
	})
})
