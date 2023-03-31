package view

import (
	"fmt"
	"time"

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

var (
	sections []creditsSection
)

func init() {
	sections = append(sections, createSection("ART & PROGRAMMING",
		"Momchil Atanasov",
	))
	sections = append(sections, createSection("NOTABLE TOOLING",
		"Visual Studio Code",
		"Blender",
		"Krita",
		"GIMP",
		"Procreate",
		"Aseprite",
		"Affinity Designer",
	))
	sections = append(sections, createSection("SPECIAL THANKS",
		"Go Developers",
		"\"GameDev БГ\" Discord Server",
		"Open-Source Code Developers",
	))
}

var CreditsScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	fadeInVisible := co.UseState(func() bool {
		return true
	})

	fadeOutVisible := co.UseState(func() bool {
		return false
	})

	onBackClicked := func() {
		mvc.Dispatch(scope, action.ChangeView{
			ViewName: model.ViewNameHome,
		})
	}

	onCreditsFinished := func() {
		fadeOutVisible.Set(true)
	}

	onFadeInFinished := func() {
		fadeInVisible.Set(false)
	}

	onFadeOutFinished := func() {
		onBackClicked()
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

			co.WithChild("button", co.New(widget.Button, func() {
				co.WithData(widget.ButtonData{
					Text: "Back",
				})
				co.WithLayoutData(mat.LayoutData{
					HorizontalCenter: opt.V(0),
					Bottom:           opt.V(100),
				})
				co.WithCallbackData(widget.ButtonCallbackData{
					ClickListener: onBackClicked,
				})
			}))
		}))

		co.WithChild("content-pane", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: opt.V(ui.RGB(0x11, 0x11, 0x11)),
				Layout:          mat.NewFillLayout(),
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(200),
				Right:  opt.V(0),
			})

			co.WithChild("scroll-pane", co.New(widget.AutoScroll, func() {
				co.WithData(widget.AutoScrollData{
					Velocity: 50.0,
				})
				co.WithCallbackData(widget.AutoScrollCallbackData{
					OnFinished: onCreditsFinished,
				})

				co.WithChild("credits-list", co.New(mat.Element, func() {
					co.WithData(mat.ElementData{
						Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
							ContentAlignment: mat.AlignmentCenter,
							ContentSpacing:   20,
						}),
					})
					co.WithLayoutData(mat.LayoutData{
						GrowHorizontally: true,
					})

					co.WithChild("header-spacing", co.New(mat.Spacing, func() {
						co.WithData(mat.SpacingData{
							Width:  10,
							Height: 500,
						})
					}))

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

					co.WithChild("section-spacing", co.New(mat.Spacing, func() {
						co.WithData(mat.SpacingData{
							Width:  10,
							Height: 100,
						})
					}))

					for i, section := range sections {
						co.WithChild(fmt.Sprintf("section-%d-title", i), co.New(mat.Label, func() {
							co.WithData(mat.LabelData{
								Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
								FontSize:  opt.V(float32(24)),
								FontColor: opt.V(theme.PrimaryColor),
								Text:      section.Title,
							})
						}))
						for j, item := range section.Items {
							co.WithChild(fmt.Sprintf("section-%d-item-%d", i, j), co.New(mat.Label, func() {
								co.WithData(mat.LabelData{
									Font:      co.OpenFont(scope, "mat:///roboto-regular.ttf"),
									FontSize:  opt.V(float32(32)),
									FontColor: opt.V(theme.PrimaryOverColor),
									Text:      item,
								})
							}))
						}
						co.WithChild(fmt.Sprintf("post-section-%d-spacing", i), co.New(mat.Spacing, func() {
							co.WithData(mat.SpacingData{
								Width:  10,
								Height: 20,
							})
						}))
					}

					co.WithChild("thank-you-spacing", co.New(mat.Spacing, func() {
						co.WithData(mat.SpacingData{
							Width:  10,
							Height: 300,
						})
					}))

					co.WithChild("thank-you", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(64.0)),
							FontColor: opt.V(ui.White()),
							Text:      "THANK YOU",
						})
					}))

					co.WithChild("footer-spacing", co.New(mat.Spacing, func() {
						co.WithData(mat.SpacingData{
							Width:  10,
							Height: 300,
						})
					}))
				}))
			}))

			if fadeInVisible.Get() {
				co.WithChild("fade-in", co.New(widget.FadeIn, func() {
					co.WithData(widget.FadeInData{
						Duration: time.Second,
					})
					co.WithCallbackData(widget.FadeInCallbackData{
						OnFinished: onFadeInFinished,
					})
				}))
			}

			if fadeOutVisible.Get() {
				co.WithChild("fade-out", co.New(widget.FadeOut, func() {
					co.WithData(widget.FadeOutData{
						Duration: time.Second,
					})
					co.WithCallbackData(widget.FadeOutCallbackData{
						OnFinished: onFadeOutFinished,
					})
				}))
			}
		}))
	})
})

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
