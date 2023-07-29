package widget

import (
	"math"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Loading = co.Define(&loadingComponent{})

type loadingComponent struct {
	co.BaseComponent

	lastTick   time.Time
	greenAngle sprec.Angle
	redAngle   sprec.Angle
}

func (c *loadingComponent) OnCreate() {
	c.lastTick = time.Now()
}

func (c *loadingComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(300, 300)),
		})
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithChildren(c.Properties().Children())
	})
}

func (c *loadingComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	const (
		radius          = float32(140.0)
		redAngleSpeed   = float32(210.0)
		greenAngleSpeed = float32(90.0)
		anglePrecision  = 2
	)

	currentTime := time.Now()
	elapsedSeconds := float32(currentTime.Sub(c.lastTick).Seconds())
	c.redAngle -= sprec.Degrees(elapsedSeconds * redAngleSpeed)
	c.greenAngle += sprec.Degrees(elapsedSeconds * greenAngleSpeed)
	c.lastTick = currentTime

	canvas.Push()
	drawBounds := canvas.DrawBounds(element, false)
	canvas.Translate(sprec.Vec2Quot(drawBounds.Size, 2.0))

	canvas.Reset()
	for angle := sprec.Degrees(0); angle < sprec.Degrees(360); angle += sprec.Degrees(anglePrecision) {
		distanceToRed := c.angleDistance(c.redAngle, angle).Degrees()
		distanceToGreen := c.angleDistance(c.greenAngle, angle).Degrees()
		canvas.SetStrokeSize(c.sizeFromDistances(distanceToRed, distanceToGreen))
		canvas.SetStrokeColor(c.colorFromDistances(distanceToRed, distanceToGreen))
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

func (c *loadingComponent) angleMod360(angle sprec.Angle) sprec.Angle {
	degrees := float64(angle.Degrees())
	degrees = math.Mod(degrees, 360.0)
	if degrees < 0.0 {
		degrees += 360.0
	}
	return sprec.Degrees(float32(degrees))
}

func (c *loadingComponent) angleDistance(a, b sprec.Angle) sprec.Angle {
	modA := c.angleMod360(a)
	modB := c.angleMod360(b)
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

func (*loadingComponent) colorFromDistances(redDistance, greenDistance float32) ui.Color {
	return ui.RGB(
		uint8(2550.0/(redDistance+10.0)),
		uint8(2550.0/(greenDistance+10.0)),
		0xFF,
	)
}

func (*loadingComponent) sizeFromDistances(redDistance, greenDistance float32) float32 {
	redSize := 100.0 / (redDistance/4.0 + 5.0)
	greenSize := 100.0 / (greenDistance/4.0 + 5.0)
	return 3.0 + redSize + greenSize
}
