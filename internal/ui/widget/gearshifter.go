package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type GearShifterSource interface {
	IsDrive() bool
}

type GearShifterData struct {
	Source GearShifterSource
}

var GearShifter = co.Define(&gearShifterComponent{})

type gearShifterComponent struct {
	co.BaseComponent

	driveImage   *ui.Image
	reverseImage *ui.Image
	source       GearShifterSource
}

func (c *gearShifterComponent) OnCreate() {
	data := co.GetData[GearShifterData](c.Properties())
	c.source = data.Source
	c.driveImage = co.OpenImage(c.Scope(), "ui/images/drive.png")
	c.reverseImage = co.OpenImage(c.Scope(), "ui/images/reverse.png")
}

func (c *gearShifterComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(200, 150)),
		})
		co.WithLayoutData(c.Properties().LayoutData())
	})
}

func (c *gearShifterComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	drawBounds := canvas.DrawBounds(element, false)

	image := c.reverseImage
	if c.source.IsDrive() {
		image = c.driveImage
	}

	canvas.Reset()
	canvas.Rectangle(
		drawBounds.Position,
		drawBounds.Size,
	)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       image,
		ImageOffset: drawBounds.Position,
		ImageSize:   drawBounds.Size,
	})
	element.Invalidate() // force redraw
}
