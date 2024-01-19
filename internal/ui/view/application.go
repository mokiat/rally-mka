package view

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

var Application = mvc.EventListener(co.Define(&applicationComponent{}))

type applicationComponent struct {
	co.BaseComponent

	appModel     *model.Application
	homeModel    *model.Home
	loadingModel *model.Loading
	playModel    *model.Play
}

func (c *applicationComponent) OnCreate() {
	eventBus := co.TypedValue[*mvc.EventBus](c.Scope())
	c.appModel = model.NewApplication(eventBus)
	c.homeModel = model.NewHome(eventBus)
	c.loadingModel = model.NewLoading(eventBus)
	c.playModel = model.NewPlay(eventBus)
}

func (c *applicationComponent) Render() co.Instance {
	return co.New(std.Switch, func() {
		co.WithData(std.SwitchData{
			ChildKey: c.appModel.ActiveView(),
		})

		co.WithChild(model.ViewNameIntro, co.New(IntroScreen, func() {
			co.WithData(IntroScreenData{
				AppModel:     c.appModel,
				Home:         c.homeModel,
				LoadingModel: c.loadingModel,
			})
		}))
		co.WithChild(model.ViewNameLoading, co.New(LoadingScreen, func() {
			co.WithData(LoadingScreenData{
				AppModel: c.appModel,
				Model:    c.loadingModel,
			})
		}))
		co.WithChild(model.ViewNameHome, co.New(HomeScreen, func() {
			co.WithData(HomeScreenData{
				AppModel: c.appModel,
				Loading:  c.loadingModel,
				Home:     c.homeModel,
				Play:     c.playModel,
			})
		}))
		co.WithChild(model.ViewNamePlay, co.New(PlayScreen, func() {
			co.WithData(PlayScreenData{
				AppModel: c.appModel,
				Play:     c.playModel,
			})
		}))
		co.WithChild(model.ViewNameLicenses, co.New(LicensesScreen, func() {
			co.WithData(LicensesScreenData{
				AppModel: c.appModel,
			})
		}))
		co.WithChild(model.ViewNameCredits, co.New(CreditsScreen, func() {
			co.WithData(CreditsScreenData{
				AppModel: c.appModel,
			})
		}))
	})
}

func (c *applicationComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.ActiveViewChangedEvent:
		c.Invalidate()
	}
}
