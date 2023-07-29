package widget

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type FadeInData struct {
	Duration time.Duration
}

var defaultFadeInData = FadeInData{
	Duration: time.Second,
}

type FadeInCallbackData struct {
	OnFinished func()
}

var defaultFadeInCallbackData = FadeInCallbackData{
	OnFinished: func() {},
}

var FadeIn = co.Define(&fadeInComponent{})

type fadeInComponent struct {
	co.BaseComponent

	opacity  float64
	duration float64
	lastTick time.Time

	onFinished std.OnActionFunc
}

func (c *fadeInComponent) OnCreate() {
	c.lastTick = time.Now()
	c.opacity = 0.0

	data := co.GetOptionalData(c.Properties(), defaultFadeInData)
	c.duration = data.Duration.Seconds()

	callbackData := co.GetOptionalCallbackData(c.Properties(), defaultFadeInCallbackData)
	c.onFinished = callbackData.OnFinished
}

func (c *fadeInComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(false),
		})
		co.WithChildren(c.Properties().Children())
	})
}

func (c *fadeInComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(c.lastTick).Seconds()
	c.lastTick = currentTime

	wasRunning := c.opacity < 1.0
	c.opacity += elapsedSeconds / c.duration
	isRunning := c.opacity < 1.0
	c.opacity = dprec.Clamp(c.opacity, 0.0, 1.0)

	if wasRunning && !isRunning {
		c.onFinished()
	}

	drawBounds := canvas.DrawBounds(element, false)

	canvas.Reset()
	canvas.Rectangle(
		drawBounds.Position,
		drawBounds.Size,
	)
	canvas.Fill(ui.Fill{
		Color: ui.RGBA(0, 0, 0, 255-uint8(c.opacity*255)),
	})

	// Force redraw.
	if c.opacity < 1.0 {
		element.Invalidate()
	}
}
