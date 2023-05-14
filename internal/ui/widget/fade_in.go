package widget

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
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

var FadeIn = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data         = co.GetOptionalData(props, defaultFadeInData)
		callbackData = co.GetOptionalCallbackData(props, defaultFadeInCallbackData)
	)

	essence := co.UseState(func() *fadeInEssence {
		return &fadeInEssence{
			lastTick: time.Now(),
		}
	}).Get()
	essence.duration = data.Duration.Seconds()
	essence.onFinished = callbackData.OnFinished

	return co.New(mat.Element, func() {
		co.WithData(co.ElementData{
			Essence:   essence,
			Focusable: opt.V(false),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
})

var _ ui.ElementRenderHandler = (*fadeInEssence)(nil)

type fadeInEssence struct {
	opacity    float64
	duration   float64
	lastTick   time.Time
	onFinished func()
}

func (e *fadeInEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(e.lastTick).Seconds()
	e.lastTick = currentTime

	wasRunning := e.opacity < 1.0
	e.opacity += elapsedSeconds / e.duration
	isRunning := e.opacity < 1.0
	e.opacity = dprec.Clamp(e.opacity, 0.0, 1.0)

	if wasRunning && !isRunning {
		e.onFinished()
	}

	bounds := element.Bounds()

	canvas.Reset()
	canvas.Rectangle(
		sprec.NewVec2(0, 0),
		sprec.NewVec2(float32(bounds.Width), float32(bounds.Height)),
	)
	canvas.Fill(ui.Fill{
		Color: ui.RGBA(0, 0, 0, 255-uint8(e.opacity*255)),
	})

	// Force redraw.
	if e.opacity < 1.0 {
		element.Invalidate()
	}
}
