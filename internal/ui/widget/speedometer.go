package widget

import (
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
	MaxVelocity float64
	Source      SpeedometerSource
}

var Speedometer = co.Define(&speedometerComponent{})

type speedometerComponent struct {
	co.BaseComponent

	speedometerImage *ui.Image
	needleImage      *ui.Image

	maxVelocity float64
	source      SpeedometerSource
}

func (c *speedometerComponent) OnUpsert() {
	c.speedometerImage = co.OpenImage(c.Scope(), "ui/images/speedometer.png")
	c.needleImage = co.OpenImage(c.Scope(), "ui/images/needle.png")

	data := co.GetData[SpeedometerData](c.Properties())
	c.maxVelocity = data.MaxVelocity
	c.source = data.Source
}

func (c *speedometerComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(300, 150)),
		})
		co.WithChildren(c.Properties().Children())
	})
}

func (c *speedometerComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	bounds := element.Bounds()
	area := sprec.Vec2{
		X: float32(bounds.Width),
		Y: float32(bounds.Height),
	}

	canvas.Push()
	canvas.Reset()
	canvas.Rectangle(sprec.ZeroVec2(), area)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       c.speedometerImage,
		ImageOffset: sprec.ZeroVec2(),
		ImageSize:   area,
	})

	needleSize := sprec.NewVec2(34.0, 150.0)
	canvas.Translate(sprec.NewVec2(
		area.X/2.0,
		area.Y-20,
	))
	velocity := c.source.Velocity() * 3.6 // from m/s to km/h

	canvas.Rotate(sprec.Degrees(-90 + 180.0*(float32(velocity/c.maxVelocity))))
	canvas.Reset()
	canvas.Rectangle(sprec.NewVec2(-needleSize.X/2.0, 20-needleSize.Y), needleSize)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       c.needleImage,
		ImageOffset: sprec.NewVec2(-needleSize.X/2.0, 20-needleSize.Y),
		ImageSize:   needleSize,
	})

	canvas.Pop()

	element.Invalidate() // force redraw
}
