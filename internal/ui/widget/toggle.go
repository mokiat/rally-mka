package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

const (
	toggleFontSize = float32(24.0)
)

type ToggleData struct {
	Text     string
	Selected bool
}

var toggleDefaultData = ToggleData{}

type ToggleCallbackData struct {
	OnToggle std.OnActionFunc
}

var toggleDefaultCallbackData = ToggleCallbackData{
	OnToggle: func() {},
}

var Toggle = co.Define(&toggleComponent{})

var _ ui.ElementRenderHandler = (*toggleComponent)(nil)

type toggleComponent struct {
	std.BaseButtonComponent

	Scope      co.Scope      `co:"scope"`
	Properties co.Properties `co:"properties"`

	font       *ui.Font
	text       string
	isSelected bool
}

func (c *toggleComponent) OnUpsert() {
	c.font = co.OpenFont(c.Scope, "ui:///roboto-bold.ttf")

	data := co.GetOptionalData(c.Properties, toggleDefaultData)
	c.text = data.Text
	c.isSelected = data.Selected

	callbackData := co.GetOptionalCallbackData(c.Properties, toggleDefaultCallbackData)
	c.SetOnClickFunc(callbackData.OnToggle)
}

func (c *toggleComponent) Render() co.Instance {
	padding := ui.Spacing{
		Left:   5,
		Right:  5,
		Top:    2,
		Bottom: 2,
	}
	txtSize := c.font.TextSize(c.text, toggleFontSize)

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties.LayoutData())
		co.WithData(std.ElementData{
			Essence: c,
			Padding: padding,
			IdealSize: opt.V(
				ui.NewSize(int(txtSize.X)+int(txtSize.Y)+5, int(txtSize.Y)).Grow(padding.Size()),
			),
		})
	})
}

func (c *toggleComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var fontColor ui.Color
	switch c.State() {
	case std.ButtonStateOver:
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	case std.ButtonStateDown:
		fontColor = ui.RGB(0x00, 0x33, 0x00)
	default:
		fontColor = ui.White()
	}
	if c.isSelected {
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	}

	bounds := element.ContentBounds() // take padding into consideration
	area := sprec.Vec2{
		X: float32(bounds.Width),
		Y: float32(bounds.Height),
	}

	canvas.Reset()
	canvas.Circle(sprec.NewVec2(area.Y/2.0, area.Y/2.0), area.Y/2.0)
	canvas.Fill(ui.Fill{
		Rule:  ui.FillRuleSimple,
		Color: fontColor,
	})

	canvas.Reset()
	canvas.FillText(c.text, sprec.NewVec2(
		float32(area.Y+10),
		float32(0.0),
	), ui.Typography{
		Font:  c.font,
		Size:  toggleFontSize,
		Color: fontColor,
	})
}
