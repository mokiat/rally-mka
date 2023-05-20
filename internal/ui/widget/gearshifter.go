package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type GearShifterSource interface {
	IsDrive() bool
}

type GearShifterData struct {
	Source GearShifterSource
}

var GearShifter = co.DefineType(&gearShifterPresenter{})

var _ ui.ElementRenderHandler = (*gearShifterPresenter)(nil)

type gearShifterPresenter struct {
	Scope      co.Scope        `co:"scope"`
	Data       GearShifterData `co:"data"`
	LayoutData any             `co:"layout"`

	driveImage   *ui.Image
	reverseImage *ui.Image
	source       GearShifterSource
}

func (p *gearShifterPresenter) OnCreate() {
	p.source = p.Data.Source
	p.driveImage = co.OpenImage(p.Scope, "ui/images/drive.png")
	p.reverseImage = co.OpenImage(p.Scope, "ui/images/reverse.png")
}

func (p *gearShifterPresenter) Render() co.Instance {
	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   p,
			IdealSize: opt.V(ui.NewSize(200, 150)),
		})
		co.WithLayoutData(p.LayoutData)
	})
}

func (p *gearShifterPresenter) OnRender(element *ui.Element, canvas *ui.Canvas) {
	bounds := element.Bounds()
	area := sprec.Vec2{
		X: float32(bounds.Width),
		Y: float32(bounds.Height),
	}

	image := p.reverseImage
	if p.source.IsDrive() {
		image = p.driveImage
	}

	canvas.Reset()
	canvas.Rectangle(sprec.ZeroVec2(), area)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       image,
		ImageOffset: sprec.ZeroVec2(),
		ImageSize:   area,
	})
	element.Invalidate() // force redraw
}
