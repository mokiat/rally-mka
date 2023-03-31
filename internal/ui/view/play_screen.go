package view

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
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
	Scope co.Scope       `co:"scope"`
	Data  PlayScreenData `co:"data"`

	controller *controller.PlayController

	rootElement *ui.Element
	exitMenu    co.Overlay
}

var _ ui.ElementKeyboardHandler = (*PlayScreenPresenter)(nil)

func (p *PlayScreenPresenter) OnCreate() {
	var context global.Context
	co.InjectContext(&context)

	// FIXME: This may actually panic if there is a third party
	// waiting / reading on this and it happens to match the Get call.
	playData, err := p.Data.Play.Data().Get()
	if err != nil {
		panic(fmt.Errorf("failed to get data: %w", err))
	}
	p.controller = controller.NewPlayController(co.Window(p.Scope).Window, context.Engine, playData)
	p.controller.Start()
}

func (p *PlayScreenPresenter) OnDelete() {
	defer p.controller.Stop()
}

func (p *PlayScreenPresenter) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	// TODO: Pass to controller
	switch event.Code {
	case ui.KeyCodeEscape:
		if event.Type == ui.KeyboardEventTypeKeyDown {
			p.controller.Pause()
			p.exitMenu = co.OpenOverlay(co.New(ExitMenu, func() {
				co.WithCallbackData(ExitMenuCallback{
					OnContinue: p.onContinue,
					OnHome:     p.onGoHome,
					OnExit:     p.onExit,
				})
			}))
		}
		return true
	default:
		return false
	}
}

func (p *PlayScreenPresenter) Render() co.Instance {
	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Reference: &p.rootElement,
			Essence:   p,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("dashboard", co.New(mat.Element, func() {
			co.WithData(mat.ElementData{
				Padding: ui.Spacing{
					Left:   20,
					Right:  20,
					Bottom: 20,
				},
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})

			co.WithLayoutData(mat.LayoutData{
				Left:   opt.V(0),
				Right:  opt.V(0),
				Bottom: opt.V(0),
			})

			co.WithChild("speedometer", co.New(widget.Speedometer, func() {
				co.WithData(widget.SpeedometerData{
					Source: p.controller,
				})

				co.WithLayoutData(mat.LayoutData{
					Left:   opt.V(0),
					Bottom: opt.V(0),
				})
			}))
		}))
	})
}

func (p *PlayScreenPresenter) onContinue() {
	p.exitMenu.Close()
	p.controller.Resume()
	co.Window(p.Scope).GrantFocus(p.rootElement)
}

func (p *PlayScreenPresenter) onGoHome() {
	// TODO: Go to loading and schedule data release
}

func (p *PlayScreenPresenter) onExit() {
	co.Window(p.Scope).Close()
}
