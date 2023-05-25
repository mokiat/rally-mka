package widget

import (
	"time"

	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
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

var AutoScroll = co.Define(&autoScrollComponent{})

type autoScrollComponent struct {
	Properties co.Properties `co:"properties"`

	lastTick   time.Time
	offsetY    float64
	maxOffsetY float64
	velocity   float64

	onFinished func()
}

func (c *autoScrollComponent) OnCreate() {
	c.lastTick = time.Now()
	c.offsetY = 0.0

	data := co.GetOptionalData(c.Properties, defaultAutoScrollData)
	c.velocity = data.Velocity

	callbackData := co.GetOptionalCallbackData(c.Properties, defaultAutoScrollCallbackData)
	c.onFinished = callbackData.OnFinished
}

func (c *autoScrollComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties.LayoutData())
		co.WithData(std.ElementData{
			Essence: c,
			Layout:  c,
		})
		co.WithChildren(c.Properties.Children())
	})
}

func (c *autoScrollComponent) Apply(element *ui.Element) {
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
			Position: ui.NewPosition(0, -int(c.offsetY)),
			Size:     childSize,
		})
	}

	c.maxOffsetY = float64(maxInt(0, maxChildSize.Height-contentBounds.Height))

	element.SetIdealSize(ui.Size{
		Width:  maxChildSize.Width + element.Padding().Horizontal(),
		Height: maxChildSize.Height + element.Padding().Vertical(),
	})
}

func (c *autoScrollComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(c.lastTick).Seconds()
	c.lastTick = currentTime

	wasEnded := c.offsetY > c.maxOffsetY
	c.offsetY += c.velocity * elapsedSeconds
	isEnded := c.offsetY > c.maxOffsetY
	if isEnded && !wasEnded {
		c.onFinished()
	}

	// Relayout and redraw
	c.Apply(element)
	element.Invalidate()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
