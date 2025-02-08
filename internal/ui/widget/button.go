package widget

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type ButtonData struct {
	Text        string
	Selected    bool
	AppearAfter time.Duration
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

	font        *ui.Font
	fontSize    float32
	text        []rune
	selected    bool
	appearAfter time.Duration
}

func (c *buttonComponent) OnCreate() {
	data := co.GetOptionalData(c.Properties(), defaultButtonData)
	c.appearAfter = data.AppearAfter
}

func (c *buttonComponent) OnUpsert() {
	c.font = co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf")
	c.fontSize = 26.0

	data := co.GetOptionalData(c.Properties(), defaultButtonData)
	c.text = []rune(data.Text)
	c.selected = data.Selected

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
	txtSize := sprec.Vec2{
		X: c.font.LineWidth(c.text, c.fontSize),
		Y: c.font.LineHeight(c.fontSize),
	}
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
	c.appearAfter -= canvas.ElapsedTime()

	var fontColor ui.Color

	switch c.State() {
	case std.ButtonStateOver:
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	case std.ButtonStateDown:
		fontColor = ui.RGB(0x00, 0x33, 0x00)
	default:
		fontColor = ui.White()
	}

	const appearDuration = time.Second
	if c.appearAfter > 0 {
		appearProgress := max(0.0, 1.0-c.appearAfter.Seconds()/appearDuration.Seconds())
		fontColor.A = byte(255 * appearProgress)
	}

	drawBounds := canvas.DrawBounds(element, true)

	canvas.Push()
	canvas.Translate(drawBounds.Position)

	canvas.Reset()
	canvas.FillTextLine(c.text, sprec.ZeroVec2(), ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: fontColor,
	})
	if c.selected {
		canvas.Reset()
		canvas.SetStrokeColor(fontColor)
		canvas.SetStrokeSize(2.0)
		canvas.MoveTo(sprec.Vec2{
			X: 0,
			Y: drawBounds.Height() + 1,
		})
		canvas.LineTo(sprec.Vec2{
			X: drawBounds.Width(),
			Y: drawBounds.Height() + 1,
		})
		canvas.Stroke()
	}

	canvas.Pop()

	if c.appearAfter > 0 {
		element.Invalidate()
	}
}
