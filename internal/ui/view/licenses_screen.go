package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
	"github.com/mokiat/rally-mka/resources"
)

var LicensesScreen = co.Define(&licensesScreenComponent{})

type LicensesScreenData struct {
	AppModel *model.Application
}

type licensesScreenComponent struct {
	co.BaseComponent

	appModel *model.Application
}

func (c *licensesScreenComponent) OnCreate() {
	data := co.GetData[LicensesScreenData](c.Properties())
	c.appModel = data.AppModel
}

func (c *licensesScreenComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          layout.Anchor(),
		})

		co.WithChild("menu-pane", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(0),
				Width:  opt.V(200),
			})
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(ui.Black()),
				Layout:          layout.Anchor(),
			})

			co.WithChild("button", co.New(widget.Button, func() {
				co.WithLayoutData(layout.Data{
					HorizontalCenter: opt.V(0),
					Bottom:           opt.V(100),
				})
				co.WithData(widget.ButtonData{
					Text: "Back",
				})
				co.WithCallbackData(widget.ButtonCallbackData{
					OnClick: c.onBackClicked,
				})
			}))
		}))

		co.WithChild("content-pane", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(200),
				Right:  opt.V(0),
			})
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(ui.RGB(0x11, 0x11, 0x11)),
				Layout:          layout.Anchor(),
			})

			co.WithChild("title", co.New(std.Label, func() {
				co.WithLayoutData(layout.Data{
					Top:              opt.V(15),
					Height:           opt.V(32),
					HorizontalCenter: opt.V(0),
				})
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(32)),
					FontColor: opt.V(ui.White()),
					Text:      "Open-Source Licenses",
				})
			}))

			co.WithChild("sub-title", co.New(std.Label, func() {
				co.WithLayoutData(layout.Data{
					Top:              opt.V(50),
					Height:           opt.V(20),
					HorizontalCenter: opt.V(0),
				})
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-italic.ttf"),
					FontSize:  opt.V(float32(20)),
					FontColor: opt.V(ui.White()),
					Text:      "- scroll to view all -",
				})
			}))

			co.WithChild("license-scroll-pane", co.New(std.ScrollPane, func() {
				co.WithLayoutData(layout.Data{
					Top:    opt.V(80),
					Bottom: opt.V(0),
					Left:   opt.V(0),
					Right:  opt.V(0),
				})
				co.WithData(std.ScrollPaneData{
					DisableHorizontal: true,
					DisableVertical:   false,
					Focused:           true,
				})

				co.WithChild("license-holder", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Padding: ui.Spacing{
							Top:    100,
							Bottom: 100,
						},
						Layout: layout.Anchor(),
					})

					co.WithChild("license-text", co.New(std.Label, func() {
						co.WithLayoutData(layout.Data{
							HorizontalCenter: opt.V(0),
							VerticalCenter:   opt.V(0),
						})
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(16)),
							FontColor: opt.V(ui.White()),
							Text:      resources.Licenses,
						})
					}))
				}))
			}))
		}))
	})
}

func (c *licensesScreenComponent) onBackClicked() {
	c.appModel.SetActiveView(model.ViewNameHome)
}
