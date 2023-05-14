package widget

import (
	"math"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

var Loading = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	essence := co.UseState(func() *loadingEssence {
		return &loadingEssence{
			lastTick: time.Now(),
		}
	}).Get()

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   essence,
			IdealSize: opt.V(ui.NewSize(300, 300)),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
})

var _ ui.ElementRenderHandler = (*loadingEssence)(nil)

type loadingEssence struct {
	lastTick   time.Time
	greenAngle sprec.Angle
	redAngle   sprec.Angle
}

func (e *loadingEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	const (
		radius          = float32(140.0)
		redAngleSpeed   = float32(210.0)
		greenAngleSpeed = float32(90.0)
		anglePrecision  = 2
	)

	currentTime := time.Now()
	elapsedSeconds := float32(currentTime.Sub(e.lastTick).Seconds())
	e.redAngle -= sprec.Degrees(elapsedSeconds * redAngleSpeed)
	e.greenAngle += sprec.Degrees(elapsedSeconds * greenAngleSpeed)
	e.lastTick = currentTime

	canvas.Push()
	bounds := element.ContentBounds()
	canvas.Translate(sprec.Vec2{
		X: float32(bounds.Width) / 2.0,
		Y: float32(bounds.Height) / 2.0,
	})

	canvas.Reset()
	for angle := sprec.Degrees(0); angle < sprec.Degrees(360); angle += sprec.Degrees(anglePrecision) {
		distanceToRed := e.angleDistance(e.redAngle, angle).Degrees()
		distanceToGreen := e.angleDistance(e.greenAngle, angle).Degrees()
		canvas.SetStrokeSize(e.sizeFromDistances(distanceToRed, distanceToGreen))
		canvas.SetStrokeColor(e.colorFromDistances(distanceToRed, distanceToGreen))
		position := sprec.NewVec2(
			sprec.Cos(angle)*radius,
			-sprec.Sin(angle)*radius,
		)
		if angle == 0 {
			canvas.MoveTo(position)
		} else {
			canvas.LineTo(position)
		}
	}
	canvas.CloseLoop()
	canvas.Stroke()

	canvas.Pop()

	element.Invalidate() // force redraw
}

func (e *loadingEssence) angleMod360(angle sprec.Angle) sprec.Angle {
	degrees := float64(angle.Degrees())
	degrees = math.Mod(degrees, 360.0)
	if degrees < 0.0 {
		degrees += 360.0
	}
	return sprec.Degrees(float32(degrees))
}

func (e *loadingEssence) angleDistance(a, b sprec.Angle) sprec.Angle {
	modA := e.angleMod360(a)
	modB := e.angleMod360(b)
	rotation360 := sprec.Degrees(360.0)
	if modA > modB {
		forwardDelta := sprec.Abs(modA - modB)
		reverseDelta := sprec.Abs(modA - modB - rotation360)
		return sprec.Min(forwardDelta, reverseDelta)
	} else {
		forwardDelta := sprec.Abs(modB - modA)
		reverseDelta := sprec.Abs(modB - modA - rotation360)
		return sprec.Min(forwardDelta, reverseDelta)
	}
}

func (*loadingEssence) colorFromDistances(redDistance, greenDistance float32) ui.Color {
	return ui.RGB(
		uint8(2550.0/(redDistance+10.0)),
		uint8(2550.0/(greenDistance+10.0)),
		0xFF,
	)
}

func (*loadingEssence) sizeFromDistances(redDistance, greenDistance float32) float32 {
	redSize := 100.0 / (redDistance/4.0 + 5.0)
	greenSize := 100.0 / (greenDistance/4.0 + 5.0)
	return 3.0 + redSize + greenSize
}
