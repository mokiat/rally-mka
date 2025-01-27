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

	appModel     *model.ApplicationModel
	errorModel   *model.ErrorModel
	loadingModel *model.LoadingModel
	homeModel    *model.HomeModel
	playModel    *model.PlayModel
}

func (c *applicationComponent) OnCreate() {
	eventBus := co.TypedValue[*mvc.EventBus](c.Scope())
	c.appModel = model.NewApplicationModel(eventBus)
	c.errorModel = model.NewErrorModel()
	c.loadingModel = model.NewLoadingModel()
	c.homeModel = model.NewHomeModel()
	c.playModel = model.NewPlayModel()
}

func (c *applicationComponent) Render() co.Instance {
	return co.New(std.Switch, func() {
		co.WithData(std.SwitchData{
			ChildKey: c.appModel.ActiveView(),
		})

		co.WithChild(model.ViewNameIntro, co.New(IntroScreen, func() {
			co.WithData(IntroScreenData{
				AppModel:     c.appModel,
				ErrorMdoel:   c.errorModel,
				LoadingModel: c.loadingModel,
				HomeModel:    c.homeModel,
			})
		}))
		co.WithChild(model.ViewNameError, co.New(ErrorScreen, func() {
			co.WithData(ErrorScreenData{
				ErrorModel: c.errorModel,
			})
		}))
		co.WithChild(model.ViewNameLoading, co.New(LoadingScreen, func() {
			co.WithData(LoadingScreenData{
				AppModel:     c.appModel,
				LoadingModel: c.loadingModel,
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
		co.WithChild(model.ViewNameHome, co.New(HomeScreen, func() {
			co.WithData(HomeScreenData{
				AppModel:     c.appModel,
				ErrorModel:   c.errorModel,
				LoadingModel: c.loadingModel,
				HomeModel:    c.homeModel,
				PlayModel:    c.playModel,
			})
		}))
		co.WithChild(model.ViewNamePlay, co.New(PlayScreen, func() {
			co.WithData(PlayScreenData{
				AppModel:  c.appModel,
				PlayModel: c.playModel,
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
