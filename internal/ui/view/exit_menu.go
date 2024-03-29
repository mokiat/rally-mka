package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var ExitMenu = co.Define(&exitMenuComponent{})

type ExitMenuCallback struct {
	OnContinue std.OnActionFunc
	OnHome     std.OnActionFunc
	OnExit     std.OnActionFunc
}

type exitMenuComponent struct {
	co.BaseComponent

	onContinue std.OnActionFunc
	onHome     std.OnActionFunc
	onExit     std.OnActionFunc
}

func (c *exitMenuComponent) OnUpsert() {
	callbackData := co.GetCallbackData[ExitMenuCallback](c.Properties())
	c.onContinue = callbackData.OnContinue
	c.onHome = callbackData.OnHome
	c.onExit = callbackData.OnExit
}

func (c *exitMenuComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {
	case ui.KeyCodeEscape:
		if event.Action == ui.KeyboardActionUp {
			c.onContinue()
		}
		return true
	default:
		return false
	}
}

func (c *exitMenuComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    layout.Fill(),
		})

		co.WithChild("background", co.New(std.Container, func() {
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(ui.RGBA(0x00, 0x00, 0x00, 0xAA)),
				Layout:          layout.Anchor(),
			})

			co.WithChild("pane", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					Top:    opt.V(0),
					Bottom: opt.V(0),
					Left:   opt.V(0),
					Width:  opt.V(320),
				})
				co.WithData(std.ContainerData{
					BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 192)),
					Layout:          layout.Anchor(),
				})

				co.WithChild("holder", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						Left:           opt.V(75),
						VerticalCenter: opt.V(0),
					})
					co.WithData(std.ElementData{
						Layout: layout.Vertical(layout.VerticalSettings{
							ContentAlignment: layout.HorizontalAlignmentLeft,
							ContentSpacing:   15,
						}),
					})

					co.WithChild("continue-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Continue",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onContinue,
						})
					}))

					co.WithChild("home-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Main Menu",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onHome,
						})
					}))

					co.WithChild("exit-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Exit",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onExit,
						})
					}))
				}))
			}))
		}))
	})
}
