package widget

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

type FadeOutData struct {
	Duration time.Duration
}

var defaultFadeOutData = FadeOutData{
	Duration: time.Second,
}

type FadeOutCallbackData struct {
	OnFinished func()
}

var defaultFadeOutCallbackData = FadeOutCallbackData{
	OnFinished: func() {},
}

var FadeOut = co.Define(&fadeOutComponent{})

type fadeOutComponent struct {
	co.BaseComponent

	lastTick time.Time
	opacity  float64
	duration float64

	onFinished func()
}

func (c *fadeOutComponent) OnCreate() {
	c.lastTick = time.Now()
	c.opacity = 0.0

	data := co.GetOptionalData(c.Properties(), defaultFadeOutData)
	c.duration = data.Duration.Seconds()

	callbackData := co.GetOptionalCallbackData(c.Properties(), defaultFadeOutCallbackData)
	c.onFinished = callbackData.OnFinished
}

func (c *fadeOutComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(false),
		})
		co.WithChildren(c.Properties().Children())
	})
}

func (c *fadeOutComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
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

	bounds := element.Bounds()

	canvas.Reset()
	canvas.Rectangle(
		sprec.NewVec2(0, 0),
		sprec.NewVec2(float32(bounds.Width), float32(bounds.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: ui.RGBA(0, 0, 0, uint8(c.opacity*255)),
	})

	// Force redraw.
	if c.opacity < 1.0 {
		element.Invalidate()
	}
}
