package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Separator = co.Define(&separatorComponent{})

type separatorComponent struct {
	co.BaseComponent
}

var _ ui.ElementRenderHandler = (*separatorComponent)(nil)

func (c *separatorComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence: c,
		})
	})
}

func (c *separatorComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	bounds := element.Bounds() // take padding into consideration
	area := sprec.Vec2{
		X: float32(bounds.Width),
		Y: float32(bounds.Height),
	}

	canvas.Reset()
	canvas.Rectangle(sprec.ZeroVec2(), area)
	canvas.Fill(ui.Fill{
		Rule:  ui.FillRuleSimple,
		Color: ui.Black(),
	})
}
