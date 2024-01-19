package view

import (
	"fmt"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/theme"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var sections = func() []creditsSection {
	return []creditsSection{
		createSection("ART & PROGRAMMING",
			"Momchil Atanasov",
		),
		createSection("NOTABLE TOOLING",
			"Visual Studio Code",
			"Blender",
			"Affinity Designer",
			"Procreate",
			"GIMP",
		),
		createSection("SPECIAL THANKS",
			"Go Developers for the brilliant programming language",
			"\"GameDev БГ\" Discord server for provided support",
			"Open-source developers for used libraries and tools",
			"Grant Abbitt for video tutorials",
			"Erin Catto for articles and videos",
		),
	}
}()

var CreditsScreen = co.Define(&creditsScreenComponent{})

type CreditsScreenData struct {
	AppModel *model.Application
}

type creditsScreenComponent struct {
	co.BaseComponent

	appModel *model.Application

	fadeInVisible  bool
	fadeOutVisible bool
}

func (c *creditsScreenComponent) OnCreate() {
	data := co.GetData[CreditsScreenData](c.Properties())
	c.appModel = data.AppModel

	c.fadeInVisible = true
	c.fadeOutVisible = false
}

func (c *creditsScreenComponent) Render() co.Instance {
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
				Layout:          layout.Fill(),
			})

			co.WithChild("scroll-pane", co.New(widget.AutoScroll, func() {
				co.WithData(widget.AutoScrollData{
					Velocity: 50.0,
				})
				co.WithCallbackData(widget.AutoScrollCallbackData{
					OnFinished: c.onCreditsFinished,
				})

				co.WithChild("credits-list", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Layout: layout.Vertical(layout.VerticalSettings{
							ContentAlignment: layout.HorizontalAlignmentCenter,
							ContentSpacing:   20,
						}),
					})

					co.WithChild("header-spacing", co.New(std.Spacing, func() {
						co.WithData(std.SpacingData{
							Size: ui.NewSize(10, 500),
						})
					}))

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

					co.WithChild("section-spacing", co.New(std.Spacing, func() {
						co.WithData(std.SpacingData{
							Size: ui.NewSize(10, 100),
						})
					}))

					for i, section := range sections {
						co.WithChild(fmt.Sprintf("section-%d-title", i), co.New(std.Label, func() {
							co.WithData(std.LabelData{
								Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
								FontSize:  opt.V(float32(24)),
								FontColor: opt.V(theme.PrimaryColor),
								Text:      section.Title,
							})
						}))
						for j, item := range section.Items {
							co.WithChild(fmt.Sprintf("section-%d-item-%d", i, j), co.New(std.Label, func() {
								co.WithData(std.LabelData{
									Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
									FontSize:  opt.V(float32(32)),
									FontColor: opt.V(theme.PrimaryOverColor),
									Text:      item,
								})
							}))
						}
						co.WithChild(fmt.Sprintf("post-section-%d-spacing", i), co.New(std.Spacing, func() {
							co.WithData(std.SpacingData{
								Size: ui.NewSize(10, 20),
							})
						}))
					}

					co.WithChild("thank-you-spacing", co.New(std.Spacing, func() {
						co.WithData(std.SpacingData{
							Size: ui.NewSize(10, 300),
						})
					}))

					co.WithChild("thank-you", co.New(std.Label, func() {
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(64.0)),
							FontColor: opt.V(ui.White()),
							Text:      "THANK YOU",
						})
					}))

					co.WithChild("footer-spacing", co.New(std.Spacing, func() {
						co.WithData(std.SpacingData{
							Size: ui.NewSize(10, 300),
						})
					}))
				}))
			}))

			if c.fadeInVisible {
				co.WithChild("fade-in", co.New(widget.FadeIn, func() {
					co.WithData(widget.FadeInData{
						Duration: time.Second,
					})
					co.WithCallbackData(widget.FadeInCallbackData{
						OnFinished: c.onFadeInFinished,
					})
				}))
			}

			if c.fadeOutVisible {
				co.WithChild("fade-out", co.New(widget.FadeOut, func() {
					co.WithData(widget.FadeOutData{
						Duration: time.Second,
					})
					co.WithCallbackData(widget.FadeOutCallbackData{
						OnFinished: c.onFadeOutFinished,
					})
				}))
			}
		}))
	})
}

func (c *creditsScreenComponent) onBackClicked() {
	c.appModel.SetActiveView(model.ViewNameHome)
}

func (c *creditsScreenComponent) onCreditsFinished() {
	c.fadeOutVisible = true
	c.Invalidate()
}

func (c *creditsScreenComponent) onFadeInFinished() {
	c.fadeInVisible = false
	c.Invalidate()
}

func (c *creditsScreenComponent) onFadeOutFinished() {
	c.onBackClicked()
	c.Invalidate()
}

func createSection(title string, items ...string) creditsSection {
	return creditsSection{
		Title: title,
		Items: items,
	}
}

type creditsSection struct {
	Title string
	Items []string
}
