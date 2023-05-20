package view

import (
	"fmt"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/metrics"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var PlayScreen = co.DefineType(&PlayScreenPresenter{})

type PlayScreenData struct {
	Play *model.Play
}

type PlayScreenPresenter struct {
	Scope      co.Scope       `co:"scope"`
	Data       PlayScreenData `co:"data"`
	Invalidate func()         `co:"invalidate"`

	hideCursor bool
	controller *controller.PlayController

	debugVisible       bool
	debugRegions       []metrics.RegionStat
	debugRegionsTicker *time.Ticker
	debugRegionsStop   chan struct{}

	rootElement *ui.Element
	exitMenu    co.Overlay
}

var _ ui.ElementKeyboardHandler = (*PlayScreenPresenter)(nil)
var _ ui.ElementMouseHandler = (*PlayScreenPresenter)(nil)

func (p *PlayScreenPresenter) OnCreate() {
	var context global.Context
	co.InjectContext(&context)

	// FIXME: This is ugly and complicated. Come up with a better API
	// than what Go provides that is integrated into component library and
	// handles everything (cleanup, thread scheduling, etc).
	p.debugRegionsTicker = time.NewTicker(time.Second)
	p.debugRegionsStop = make(chan struct{})
	go func() {
		for {
			select {
			case <-p.debugRegionsTicker.C:
				co.Schedule(func() {
					p.debugRegions = metrics.RegionStats()
					p.Invalidate()
				})
			case <-p.debugRegionsStop:
				return
			}
		}
	}()

	// FIXME: This may actually panic if there is a third party
	// waiting / reading on this and it happens to match the Get call.
	playData, err := p.Data.Play.Data().Get()
	if err != nil {
		panic(fmt.Errorf("failed to get data: %w", err))
	}
	p.controller = controller.NewPlayController(co.Window(p.Scope).Window, context.Engine, playData)
	p.controller.Start(playData.Environment, playData.Controller)

	p.hideCursor = playData.Controller != data.ControllerMouse
	co.Window(p.Scope).SetCursorVisible(!p.hideCursor)
}

func (p *PlayScreenPresenter) OnDelete() {
	defer p.controller.Stop()
	defer p.debugRegionsTicker.Stop()
	defer close(p.debugRegionsStop)
	defer co.Window(p.Scope).SetCursorVisible(true)
}

func (p *PlayScreenPresenter) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	return p.controller.OnMouseEvent(element, event)
}

func (p *PlayScreenPresenter) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {
	case ui.KeyCodeEscape:
		if event.Type == ui.KeyboardEventTypeKeyUp {
			p.controller.Pause()
			co.Window(p.Scope).SetCursorVisible(true)
			p.exitMenu = co.OpenOverlay(p.Scope, co.New(ExitMenu, func() {
				co.WithCallbackData(ExitMenuCallback{
					OnContinue: p.onContinue,
					OnHome:     p.onGoHome,
					OnExit:     p.onExit,
				})
			}))
		}
		return true
	case ui.KeyCodeTab:
		if event.Type == ui.KeyboardEventTypeKeyDown {
			p.debugVisible = !p.debugVisible
			p.Invalidate()
		}
		return true
	case ui.KeyCodeEnter:
		if event.Type == ui.KeyboardEventTypeKeyDown {
			p.controller.ToggleCamera()
		}
		return true
	default:
		return p.controller.OnKeyboardEvent(event)
	}
}

func (p *PlayScreenPresenter) Render() co.Instance {
	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Reference: &p.rootElement,
			Essence:   p,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    layout.Anchor(),
		})

		if p.debugVisible {
			co.WithChild("regions", co.New(widget.RegionBlock, func() {
				co.WithData(widget.RegionBlockData{
					Regions: p.debugRegions,
				})
				co.WithLayoutData(layout.Data{
					Top:   opt.V(0),
					Left:  opt.V(0),
					Right: opt.V(0),
				})
			}))
		}

		co.WithChild("dashboard", co.New(mat.Element, func() {
			co.WithLayoutData(layout.Data{
				Left:   opt.V(0),
				Right:  opt.V(0),
				Bottom: opt.V(0),
			})
			co.WithData(mat.ElementData{
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
					Source: p.controller,
				})
			}))

			co.WithChild("gearshifter", co.New(widget.GearShifter, func() {
				co.WithLayoutData(layout.Data{
					Right:  opt.V(20),
					Bottom: opt.V(0),
				})
				co.WithData(widget.GearShifterData{
					Source: p.controller,
				})
			}))
		}))
	})
}

func (p *PlayScreenPresenter) onContinue() {
	p.exitMenu.Close()
	p.controller.Resume()
	co.Window(p.Scope).GrantFocus(p.rootElement)
	co.Window(p.Scope).SetCursorVisible(!p.hideCursor)
}

func (p *PlayScreenPresenter) onGoHome() {
	p.exitMenu.Close()
	mvc.Dispatch(p.Scope, action.ChangeView{
		ViewName: model.ViewNameHome,
	})
}

func (p *PlayScreenPresenter) onExit() {
	co.Window(p.Scope).Close()
}
