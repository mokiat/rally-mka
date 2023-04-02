package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

const (
	toggleFontSize = float32(24.0)
)

type ToggleData struct {
	Text     string
	Selected bool
}

type ToggleCallbackData struct {
	ClickListener mat.ClickListener
}

var Toggle = co.DefineType(&togglePresenter{})

var _ ui.ElementRenderHandler = (*togglePresenter)(nil)

type togglePresenter struct {
	*mat.ButtonBaseEssence

	Scope        co.Scope           `co:"scope"`
	Data         ToggleData         `co:"data"`
	CallbackData ToggleCallbackData `co:"callback"`
	LayoutData   mat.LayoutData     `co:"layout"`

	font *ui.Font
}

func (p *togglePresenter) OnCreate() {
	if p.CallbackData.ClickListener == nil {
		p.CallbackData.ClickListener = func() {}
	}
	p.ButtonBaseEssence = mat.NewButtonBaseEssence(p.CallbackData.ClickListener)

	p.font = co.OpenFont(p.Scope, "mat:///roboto-bold.ttf")
}

func (p *togglePresenter) Render() co.Instance {
	padding := ui.Spacing{
		Left:   5,
		Right:  5,
		Top:    2,
		Bottom: 2,
	}

	txtSize := p.font.TextSize(p.Data.Text, toggleFontSize)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: p,
			Padding: padding,
			IdealSize: opt.V(
				ui.NewSize(int(txtSize.X)+int(txtSize.Y)+5, int(txtSize.Y)).Grow(padding.Size()),
			),
		})
		co.WithLayoutData(p.LayoutData)
	})
}

func (p *togglePresenter) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var fontColor ui.Color
	switch p.State() {
	case mat.ButtonStateOver:
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	case mat.ButtonStateDown:
		fontColor = ui.RGB(0x00, 0x33, 0x00)
	default:
		fontColor = ui.White()
	}
	if p.Data.Selected {
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
	canvas.FillText(p.Data.Text, sprec.NewVec2(
		float32(area.Y+10),
		float32(0.0),
	), ui.Typography{
		Font:  p.font,
		Size:  toggleFontSize,
		Color: fontColor,
	})
}
