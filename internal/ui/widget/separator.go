package widget

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var Separator = co.DefineType(&separatorPresenter{})

type separatorPresenter struct {
	LayoutData mat.LayoutData `co:"layout"`
}

var _ ui.ElementRenderHandler = (*separatorPresenter)(nil)

func (p *separatorPresenter) Render() co.Instance {
	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: p,
		})
		co.WithLayoutData(p.LayoutData)
	})
}

func (p *separatorPresenter) OnRender(element *ui.Element, canvas *ui.Canvas) {
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
