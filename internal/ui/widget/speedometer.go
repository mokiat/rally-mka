package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type SpeedometerSource interface {
	Velocity() float64
}

type SpeedometerData struct {
	Source SpeedometerSource
}

var Speedometer = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data = co.GetData[SpeedometerData](props)
	)

	essence := co.UseState(func() *speedometerEssence {
		return &speedometerEssence{
			speedometerImage: co.OpenImage(scope, "ui/images/speedometer.png"),
			needleImage:      co.OpenImage(scope, "ui/images/needle.png"),
			source:           data.Source,
		}
	}).Get()

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   essence,
			IdealSize: opt.V(ui.NewSize(300, 150)),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
})

var _ ui.ElementRenderHandler = (*speedometerEssence)(nil)

type speedometerEssence struct {
	speedometerImage *ui.Image
	needleImage      *ui.Image
	source           SpeedometerSource
}

func (e *speedometerEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	const (
		maxVelocity = 200.0 // TODO: Configurable
	)

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
		Image:       e.speedometerImage,
		ImageOffset: sprec.ZeroVec2(),
		ImageSize:   area,
	})

	needleSize := sprec.NewVec2(34.0, 150.0)
	canvas.Translate(sprec.NewVec2(
		area.X/2.0,
		area.Y-20,
	))
	velocity := e.source.Velocity() * 3.6 // from m/s to km/h

	canvas.Rotate(sprec.Degrees(-90 + 180.0*(float32(velocity/maxVelocity))))
	canvas.Reset()
	canvas.Rectangle(sprec.NewVec2(-needleSize.X/2.0, 20-needleSize.Y), needleSize)
	canvas.Fill(ui.Fill{
		Rule:        ui.FillRuleSimple,
		Color:       ui.White(),
		Image:       e.needleImage,
		ImageOffset: sprec.NewVec2(-needleSize.X/2.0, 20-needleSize.Y),
		ImageSize:   needleSize,
	})

	canvas.Pop()

	element.Invalidate() // force redraw
}
