package widget

import (
	"time"

	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mat"
)

type AutoScrollData struct {
	Velocity float64
}

var defaultAutoScrollData = AutoScrollData{
	Velocity: 50.0,
}

type AutoScrollCallbackData struct {
	OnFinished func()
}

var defaultAutoScrollCallbackData = AutoScrollCallbackData{
	OnFinished: func() {},
}

var AutoScroll = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data         = co.GetOptionalData(props, defaultAutoScrollData)
		callbackData = co.GetOptionalCallbackData(props, defaultAutoScrollCallbackData)
	)

	essence := co.UseState(func() *scrollPaneEssence {
		return &scrollPaneEssence{
			lastTick: time.Now(),
		}
	}).Get()
	essence.velocity = data.Velocity
	essence.onFinished = callbackData.OnFinished

	return co.New(mat.Element, func() {
		co.WithData(co.ElementData{
			Essence: essence,
			Layout:  essence,
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
})

var _ ui.Layout = (*scrollPaneEssence)(nil)
var _ ui.ElementRenderHandler = (*scrollPaneEssence)(nil)

type scrollPaneEssence struct {
	offsetY    float64
	maxOffsetY float64
	velocity   float64
	lastTick   time.Time
	onFinished func()
}

func (e *scrollPaneEssence) Apply(element *ui.Element) {
	var maxChildSize ui.Size

	contentBounds := element.ContentBounds()
	for childElement := element.FirstChild(); childElement != nil; childElement = childElement.RightSibling() {
		layoutConfig := layout.ElementData(childElement)

		childSize := childElement.IdealSize()
		if layoutConfig.Width.Specified {
			childSize.Width = layoutConfig.Width.Value
		}
		childSize.Width = maxInt(childSize.Width, contentBounds.Width)
		if layoutConfig.Height.Specified {
			childSize.Height = layoutConfig.Height.Value
		}

		maxChildSize = ui.Size{
			Width:  maxInt(maxChildSize.Width, childSize.Width),
			Height: maxInt(maxChildSize.Height, childSize.Height),
		}

		childElement.SetBounds(ui.Bounds{
			Position: ui.NewPosition(0, -int(e.offsetY)),
			Size:     childSize,
		})
	}

	e.maxOffsetY = float64(maxInt(0, maxChildSize.Height-contentBounds.Height))

	element.SetIdealSize(ui.Size{
		Width:  maxChildSize.Width + element.Padding().Horizontal(),
		Height: maxChildSize.Height + element.Padding().Vertical(),
	})
}

func (e *scrollPaneEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(e.lastTick).Seconds()
	e.lastTick = currentTime

	wasEnded := e.offsetY > e.maxOffsetY
	e.offsetY += e.velocity * elapsedSeconds
	isEnded := e.offsetY > e.maxOffsetY
	if isEnded && !wasEnded {
		e.onFinished()
	}

	// Relayout and redraw
	e.Apply(element)
	element.Invalidate()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
