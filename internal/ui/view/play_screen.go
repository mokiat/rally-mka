package view

import (
	"fmt"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/lacking/util/metrics"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var PlayScreen = co.Define(&playScreenComponent{})

type PlayScreenData struct {
	Play *model.Play
}

type playScreenComponent struct {
	co.BaseComponent

	hideCursor bool
	controller *controller.PlayController

	debugVisible       bool
	debugRegions       []metrics.RegionStat
	debugRegionsTicker *time.Ticker
	debugRegionsStop   chan struct{}

	rootElement *ui.Element
	exitMenu    co.Overlay
}

var _ ui.ElementKeyboardHandler = (*playScreenComponent)(nil)
var _ ui.ElementMouseHandler = (*playScreenComponent)(nil)

func (c *playScreenComponent) OnCreate() {
	context := co.TypedValue[global.Context](c.Scope())

	// FIXME: This is ugly and complicated. Come up with a better API
	// than what Go provides that is integrated into component library and
	// handles everything (cleanup, thread scheduling, etc).
	c.debugRegionsTicker = time.NewTicker(time.Second)
	c.debugRegionsStop = make(chan struct{})
	go func() {
		for {
			select {
			case <-c.debugRegionsTicker.C:
				co.Schedule(c.Scope(), func() {
					c.debugRegions = metrics.RegionStats()
					c.Invalidate()
				})
			case <-c.debugRegionsStop:
				return
			}
		}
	}()

	screenData := co.GetData[PlayScreenData](c.Properties())

	// FIXME: This may actually panic if there is a third party
	// waiting / reading on this and it happens to match the Get call.
	playData, err := screenData.Play.Data().Get()
	if err != nil {
		panic(fmt.Errorf("failed to get data: %w", err))
	}
	c.controller = controller.NewPlayController(co.Window(c.Scope()).Window, context.Engine, playData)
	c.controller.Start(playData.Environment, playData.Controller)

	c.hideCursor = playData.Controller != data.ControllerMouse
	co.Window(c.Scope()).SetCursorVisible(!c.hideCursor)
}

func (c *playScreenComponent) OnDelete() {
	defer c.controller.Stop()
	defer c.debugRegionsTicker.Stop()
	defer close(c.debugRegionsStop)
	defer co.Window(c.Scope()).SetCursorVisible(true)
}

func (c *playScreenComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.controller.OnMouseEvent(element, event)
}

func (c *playScreenComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {
	case ui.KeyCodeEscape:
		if event.Type == ui.KeyboardEventTypeKeyUp {
			c.controller.Pause()
			co.Window(c.Scope()).SetCursorVisible(true)
			c.exitMenu = co.OpenOverlay(c.Scope(), co.New(ExitMenu, func() {
				co.WithCallbackData(ExitMenuCallback{
					OnContinue: c.onContinue,
					OnHome:     c.onGoHome,
					OnExit:     c.onExit,
				})
			}))
		}
		return true
	case ui.KeyCodeTab:
		if event.Type == ui.KeyboardEventTypeKeyDown {
			c.debugVisible = !c.debugVisible
			c.Invalidate()
		}
		return true
	case ui.KeyCodeEnter:
		if event.Type == ui.KeyboardEventTypeKeyDown {
			c.controller.ToggleCamera()
		}
		return true
	default:
		return c.controller.OnKeyboardEvent(event)
	}
}

func (c *playScreenComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Reference: &c.rootElement,
			Essence:   c,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    layout.Anchor(),
		})

		if c.debugVisible {
			co.WithChild("regions", co.New(widget.RegionBlock, func() {
				co.WithData(widget.RegionBlockData{
					Regions: c.debugRegions,
				})
				co.WithLayoutData(layout.Data{
					Top:   opt.V(0),
					Left:  opt.V(0),
					Right: opt.V(0),
				})
			}))
		}

		co.WithChild("dashboard", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				Left:   opt.V(0),
				Right:  opt.V(0),
				Bottom: opt.V(0),
			})
			co.WithData(std.ElementData{
				Padding: ui.Spacing{
					Left:   20,
					Right:  20,
					Bottom: 20,
				},
				Layout: layout.Anchor(),
			})

			co.WithChild("speedometer", co.New(widget.Speedometer, func() {
				co.WithLayoutData(layout.Data{
					Left:   opt.V(20),
					Bottom: opt.V(0),
				})
				co.WithData(widget.SpeedometerData{
					MaxVelocity: 200.0,
					Source:      c.controller,
				})
			}))

			co.WithChild("gearshifter", co.New(widget.GearShifter, func() {
				co.WithLayoutData(layout.Data{
					Right:  opt.V(20),
					Bottom: opt.V(0),
				})
				co.WithData(widget.GearShifterData{
					Source: c.controller,
				})
			}))
		}))
	})
}

func (c *playScreenComponent) onContinue() {
	c.exitMenu.Close()
	c.controller.Resume()
	co.Window(c.Scope()).GrantFocus(c.rootElement)
	co.Window(c.Scope()).SetCursorVisible(!c.hideCursor)
}

func (c *playScreenComponent) onGoHome() {
	c.exitMenu.Close()
	mvc.Dispatch(c.Scope(), action.ChangeView{
		ViewName: model.ViewNameHome,
	})
}

func (c *playScreenComponent) onExit() {
	co.Window(c.Scope()).Close()
}
