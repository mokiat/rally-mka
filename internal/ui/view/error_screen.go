package view

import (
	"fmt"
	"iter"
	"strings"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"

	"github.com/mokiat/rally-mka/internal/ui/model"
)

type ErrorScreenData struct {
	ErrorModel *model.ErrorModel
}

var ErrorScreen = co.Define(&errorScreenComponent{})

var _ ui.ElementKeyboardHandler = (*errorScreenComponent)(nil)

type errorScreenComponent struct {
	co.BaseComponent

	titleFont     *ui.Font
	titleFontSize float32

	messageFont     *ui.Font
	messageFontSize float32

	message string
}

func (c *errorScreenComponent) OnCreate() {
	screenData := co.GetData[ErrorScreenData](c.Properties())
	errorModel := screenData.ErrorModel
	c.message = c.formatError(errorModel.Error())

	c.titleFont = co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf")
	c.titleFontSize = float32(48.0)

	c.messageFont = co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf")
	c.messageFontSize = float32(24.0)
}

func (c *errorScreenComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          layout.Anchor(),
		})

		co.WithChild("handler", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				Left:   opt.V(0),
				Right:  opt.V(0),
				Top:    opt.V(0),
				Bottom: opt.V(0),
			})
			co.WithData(std.ElementData{
				Essence:   c,
				Enabled:   opt.V(true),
				Focusable: opt.V(true),
				Focused:   opt.V(true),
			})
		}))

		co.WithChild("title", co.New(std.Label, func() {
			co.WithLayoutData(layout.Data{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(-150),
			})
			co.WithData(std.LabelData{
				Text:      "ERROR",
				Font:      c.titleFont,
				FontSize:  opt.V(c.titleFontSize),
				FontColor: opt.V(ui.White()),
			})
		}))

		co.WithChild("info", co.New(std.Label, func() {
			co.WithLayoutData(layout.Data{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
			co.WithData(std.LabelData{
				Text:      c.message,
				Font:      c.messageFont,
				FontSize:  opt.V(c.messageFontSize),
				FontColor: opt.V(ui.White()),
			})
		}))
	})
}

func (c *errorScreenComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if event.Action == ui.KeyboardActionUp && event.Code == ui.KeyCodeEscape {
		co.Window(c.Scope()).Close()
	}
	return true
}

func (c *errorScreenComponent) formatError(err error) string {
	wordWrap := func(text string, maxLineLength int) iter.Seq[string] {
		return func(yield func(string) bool) {
			runes := []rune(text)
			for len(runes) > maxLineLength {
				if !yield(string(runes[:maxLineLength])) {
					return
				}
				runes = runes[maxLineLength:]
			}
			if !yield(string(runes)) {
				return
			}
		}
	}

	var builder strings.Builder
	fmt.Fprintln(&builder, "The game has encountered an error. Press ESCAPE to exit.")
	fmt.Fprintln(&builder)
	fmt.Fprint(&builder, "Error: ")
	for line := range wordWrap(err.Error(), 80) {
		fmt.Fprintln(&builder, line)
	}
	return builder.String()
}
