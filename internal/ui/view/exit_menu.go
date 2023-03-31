package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var ExitMenu = co.DefineType(&ExitMenuPresenter{})

type ExitMenuCallback struct {
	OnContinue func()
	OnHome     func()
	OnExit     func()
}

type ExitMenuPresenter struct {
	CallbackData ExitMenuCallback `co:"callback"`
}

func (p *ExitMenuPresenter) Render() co.Instance {
	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   p,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    mat.NewFillLayout(),
		})

		co.WithChild("background", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: opt.V(ui.RGBA(0x00, 0x00, 0x00, 0xAA)),
				Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})

			co.WithChild("pane", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 192)),
					Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(0),
					Bottom: opt.V(0),
					Left:   opt.V(0),
					Width:  opt.V(320),
				})

				co.WithChild("holder", co.New(mat.Element, func() {
					co.WithData(mat.ElementData{
						Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
							ContentAlignment: mat.AlignmentLeft,
							ContentSpacing:   15,
						}),
					})
					co.WithLayoutData(mat.LayoutData{
						Left:           opt.V(75),
						VerticalCenter: opt.V(0),
					})

					co.WithChild("continue-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Continue",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.CallbackData.OnContinue,
						})
					}))

					co.WithChild("home-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Main Menu",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.CallbackData.OnHome,
						})
					}))

					co.WithChild("exit-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Exit",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.CallbackData.OnExit,
						})
					}))
				}))
			}))
		}))
	})
}
