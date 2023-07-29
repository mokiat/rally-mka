package view

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

var Application = mvc.Wrap(co.Define(&applicationComponent{}))

type applicationComponent struct {
	co.BaseComponent

	homeModel    *model.Home
	loadingModel *model.Loading
	playModel    *model.Play
	activeView   string
}

func (c *applicationComponent) OnUpsert() {
	data := co.GetData[*model.Application](c.Properties())
	c.homeModel = data.Home()
	c.loadingModel = data.Loading()
	c.playModel = data.Play()
	c.activeView = data.ActiveView()

	mvc.UseBinding(c.Scope(), data, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ApplicationActiveViewChange)
	})
}

func (c *applicationComponent) Render() co.Instance {
	return co.New(std.Switch, func() {
		co.WithData(std.SwitchData{
			ChildKey: c.activeView,
		})

		co.WithChild(model.ViewNameIntro, co.New(IntroScreen, func() {
			co.WithData(IntroScreenData{
				Home:         c.homeModel,
				LoadingModel: c.loadingModel,
			})
		}))
		co.WithChild(model.ViewNameLoading, co.New(LoadingScreen, func() {
			co.WithData(LoadingScreenData{
				Model: c.loadingModel,
			})
		}))
		co.WithChild(model.ViewNameHome, co.New(HomeScreen, func() {
			co.WithData(HomeScreenData{
				Loading: c.loadingModel,
				Home:    c.homeModel,
				Play:    c.playModel,
			})
		}))
		co.WithChild(model.ViewNamePlay, co.New(PlayScreen, func() {
			co.WithData(PlayScreenData{
				Play: c.playModel,
			})
		}))
		co.WithChild(model.ViewNameLicenses, co.New(LicensesScreen, nil))
		co.WithChild(model.ViewNameCredits, co.New(CreditsScreen, nil))
	})
}
