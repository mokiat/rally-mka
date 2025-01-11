package view

import (
	"fmt"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/debug/metric/metricui"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var PlayScreen = co.Define(&playScreenComponent{})

type PlayScreenData struct {
	AppModel *model.Application
	Play     *model.Play
}

type playScreenComponent struct {
	co.BaseComponent

	appModel *model.Application

	hideCursor bool
	controller *controller.PlayController

	debugVisible bool

	rootElement *ui.Element
	exitMenu    co.Overlay
}

var _ ui.ElementKeyboardHandler = (*playScreenComponent)(nil)
var _ ui.ElementMouseHandler = (*playScreenComponent)(nil)

func (c *playScreenComponent) OnCreate() {
	context := co.TypedValue[global.Context](c.Scope())
	screenData := co.GetData[PlayScreenData](c.Properties())
	c.appModel = screenData.AppModel
	playModel := screenData.Play

	// FIXME: This may actually panic if there is a third party
	// waiting / reading on this and it happens to match the Get call.
	playData, err := playModel.Data().Wait()
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
	defer co.Window(c.Scope()).SetCursorVisible(true)
}

func (c *playScreenComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return c.controller.OnMouseEvent(element, event)
}

func (c *playScreenComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {
	case ui.KeyCodeEscape:
		if event.Action == ui.KeyboardActionUp {
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
		if event.Action == ui.KeyboardActionDown {
			c.debugVisible = !c.debugVisible
			c.Invalidate()
		}
		return true
	case ui.KeyCodeEnter:
		if event.Action == ui.KeyboardActionDown {
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
			co.WithChild("flamegraph", co.New(metricui.FlameGraph, func() {
				co.WithData(metricui.FlameGraphData{
					UpdateInterval: time.Second,
				})
				co.WithLayoutData(layout.Data{
					Top:   opt.V(0),
					Left:  opt.V(0),
					Right: opt.V(0),
				})
			}))
		}

		co.WithChild("speedometer", co.New(widget.Speedometer, func() {
			co.WithLayoutData(layout.Data{
				Left:   opt.V(0),
				Bottom: opt.V(0),
			})
			co.WithData(widget.SpeedometerData{
				Source: c.controller,
			})
		}))

		co.WithChild("gearshifter", co.New(widget.GearShifter, func() {
			co.WithLayoutData(layout.Data{
				Right:  opt.V(0),
				Bottom: opt.V(0),
			})
			co.WithData(widget.GearShifterData{
				Source: c.controller,
			})
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
	c.appModel.SetActiveView(model.ViewNameHome)
}

func (c *playScreenComponent) onExit() {
	co.Window(c.Scope()).Close()
}
