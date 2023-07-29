package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type ButtonData struct {
	Text string
}

var defaultButtonData = ButtonData{
	Text: "",
}

type ButtonCallbackData struct {
	OnClick std.OnActionFunc
}

var defaultButtonCallbackData = ButtonCallbackData{
	OnClick: func() {},
}

var Button = co.Define(&buttonComponent{})

type buttonComponent struct {
	co.BaseComponent
	std.BaseButtonComponent

	font     *ui.Font
	fontSize float32
	text     string
}

func (c *buttonComponent) OnUpsert() {
	c.font = co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf")
	c.fontSize = 26.0

	data := co.GetOptionalData(c.Properties(), defaultButtonData)
	c.text = data.Text

	callbackData := co.GetOptionalCallbackData(c.Properties(), defaultButtonCallbackData)
	c.SetOnClickFunc(callbackData.OnClick)
}

func (c *buttonComponent) Render() co.Instance {
	padding := ui.Spacing{
		Left:   5,
		Right:  5,
		Top:    2,
		Bottom: 2,
	}
	txtSize := c.font.TextSize(c.text, c.fontSize)

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence: c,
			Padding: padding,
			IdealSize: opt.V(
				ui.NewSize(int(txtSize.X), int(txtSize.Y)).Grow(padding.Size()),
			),
		})
		co.WithChildren(c.Properties().Children())
	})
}

func (c *buttonComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var fontColor ui.Color

	switch c.State() {
	case std.ButtonStateOver:
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	case std.ButtonStateDown:
		fontColor = ui.RGB(0x00, 0x33, 0x00)
	default:
		fontColor = ui.White()
	}

	drawBounds := canvas.DrawBounds(element, true)
	textPosition := drawBounds.Position
	canvas.Reset()
	canvas.FillText(c.text, sprec.NewVec2(
		float32(textPosition.X),
		float32(textPosition.Y),
	), ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: fontColor,
	})
}
