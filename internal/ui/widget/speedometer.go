package widget

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type SpeedometerSource interface {
	Velocity() float64
}

type SpeedometerData struct {
	Source SpeedometerSource
}

var Speedometer = co.Define(&speedometerComponent{})

type speedometerComponent struct {
	co.BaseComponent

	panelImage      *ui.Image
	blankDigitImage *ui.Image
	digitImages     [10]*ui.Image

	source SpeedometerSource

	speed       int
	updateAfter time.Duration
}

func (c *speedometerComponent) OnUpsert() {
	c.panelImage = co.OpenImage(c.Scope(), "ui/images/speed-panel.png")
	c.blankDigitImage = co.OpenImage(c.Scope(), "ui/images/digit-blank.png")
	c.digitImages = [10]*ui.Image{
		co.OpenImage(c.Scope(), "ui/images/digit-0.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-1.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-2.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-3.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-4.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-5.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-6.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-7.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-8.png"),
		co.OpenImage(c.Scope(), "ui/images/digit-9.png"),
	}

	data := co.GetData[SpeedometerData](c.Properties())
	c.source = data.Source
}

func (c *speedometerComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(320, 128)),
		})
		co.WithChildren(c.Properties().Children())
	})
}

func (c *speedometerComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	c.updateVelocity(canvas.ElapsedTime())

	drawBounds := canvas.DrawBounds(element, false)

	canvas.Push()
	canvas.Translate(drawBounds.Position)

	canvas.Reset()
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

	canvas.Push()
	canvas.Translate(sprec.NewVec2(40.0, 70.0))
	c.drawNumber(canvas, c.speed, 3)
	canvas.Pop()

	canvas.Pop()

	element.Invalidate() // force redraw
}

func (c *speedometerComponent) updateVelocity(deltaTime time.Duration) {
	c.updateAfter -= deltaTime
	if c.updateAfter > 0 {
		return
	}
	c.updateAfter = 500 * time.Millisecond

	c.speed = int(c.source.Velocity() * 3.6) // from m/s to km/h
}

func (c *speedometerComponent) drawNumber(canvas *ui.Canvas, number int, digits int) {
	const digitOffset = 32.0 + 8.0
	canvas.Push()
	canvas.Translate(sprec.NewVec2(float32(digits-1)*digitOffset, 0.0))
	denominator := 1
	for i := range digits {
		normalized := number / denominator
		if (normalized > 0) || (i == 0) {
			c.drawDigit(canvas, normalized%10)
		} else {
			c.drawDigit(canvas, -1)
		}
		canvas.Translate(sprec.NewVec2(-digitOffset, 0.0))
		denominator *= 10
	}
	canvas.Pop()
}

func (c *speedometerComponent) drawDigit(canvas *ui.Canvas, digit int) {
	var image *ui.Image
	switch {
	case digit >= 0 && digit <= 9:
		image = c.digitImages[digit]
	default:
		image = c.blankDigitImage
	}

	size := sprec.Vec2{
		X: 32.0,
		Y: 64.0,
	}
	offset := sprec.Vec2{
		X: -size.X / 2.0,
		Y: -size.Y / 2.0,
	}

	canvas.Reset()
	canvas.Rectangle(offset, size)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       image,
		ImageOffset: offset,
		ImageSize:   size,
	})
}
