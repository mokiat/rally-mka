package view

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

var Application = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data = co.GetData[*model.Application](props)
	)

	mvc.UseBinding(data, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ApplicationActiveViewChange)
	})

	return co.New(mat.Switch, func() {
		co.WithData(mat.SwitchData{
			ChildKey: data.ActiveView(),
		})

		co.WithChild(model.ViewNameIntro, co.New(IntroScreen, func() {
			co.WithData(IntroScreenData{
				Home:         data.Home(),
				LoadingModel: data.Loading(),
			})
		}))
		co.WithChild(model.ViewNameLoading, co.New(LoadingScreen, func() {
			co.WithData(LoadingScreenData{
				Model: data.Loading(),
			})
		}))
		co.WithChild(model.ViewNameHome, co.New(HomeScreen, func() {
			co.WithData(HomeScreenData{
				Loading: data.Loading(),
				Home:    data.Home(),
				Play:    data.Play(),
			})
		}))
		co.WithChild(model.ViewNamePlay, co.New(PlayScreen, func() {
			co.WithData(PlayScreenData{
				Play: data.Play(),
			})
		}))
		co.WithChild(model.ViewNameLicenses, co.New(LicensesScreen, nil))
		co.WithChild(model.ViewNameCredits, co.New(CreditsScreen, nil))
	})
})
