package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
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

	panelImage       *ui.Image
	modeDriveImage   *ui.Image
	modeReverseImage *ui.Image

	source GearShifterSource
}

func (c *gearShifterComponent) OnCreate() {
	data := co.GetData[GearShifterData](c.Properties())
	c.source = data.Source

	c.panelImage = co.OpenImage(c.Scope(), "ui/images/gear-panel.png")
	c.modeDriveImage = co.OpenImage(c.Scope(), "ui/images/mode-drive.png")
	c.modeReverseImage = co.OpenImage(c.Scope(), "ui/images/mode-reverse.png")
}

func (c *gearShifterComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(280, 128)),
		})
		co.WithLayoutData(c.Properties().LayoutData())
	})
}

func (c *gearShifterComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	drawBounds := canvas.DrawBounds(element, false)

	canvas.Reset()
	canvas.Translate(drawBounds.Position)
	canvas.Rectangle(
		sprec.ZeroVec2(),
		drawBounds.Size,
	)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       c.panelImage,
		ImageOffset: sprec.ZeroVec2(),
		ImageSize:   drawBounds.Size,
	})

	modeImage := c.modeReverseImage
	if c.source.IsDrive() {
		modeImage = c.modeDriveImage
	}
	imageSize := sprec.NewVec2(130.0, 76.0)
	canvas.Push()
	canvas.Translate(sprec.NewVec2(130.0, 30.0))
	canvas.Reset()
	canvas.Rectangle(sprec.ZeroVec2(), imageSize)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       modeImage,
		ImageOffset: sprec.ZeroVec2(),
		ImageSize:   imageSize,
	})
	canvas.Pop()

	element.Invalidate() // force redraw
}
