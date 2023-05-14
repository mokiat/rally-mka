package model

import "github.com/mokiat/lacking/ui/mvc"

const (
	ViewNameIntro    ViewName = "intro"
	ViewNameHome     ViewName = "home"
	ViewNamePlay     ViewName = "play"
	ViewNameLoading  ViewName = "loading"
	ViewNameLicenses ViewName = "licenses"
	ViewNameCredits  ViewName = "credits"
)

type ViewName = string

var (
	ApplicationChange           = mvc.NewChange("application")
	ApplicationActiveViewChange = mvc.SubChange(ApplicationChange, "active_view")
)

func NewApplication() *Application {
	return &Application{
		Observable: mvc.NewObservable(),
		loading:    newLoading(),
		home:       newHome(),
		play:       newPlay(),
		activeView: ViewNameIntro,
	}
}

type Application struct {
	mvc.Observable
	loading    *Loading
	home       *Home
	play       *Play
	activeView ViewName
}

func (a *Application) Loading() *Loading {
	return a.loading
}

func (a *Application) Home() *Home {
	return a.home
}

func (a *Application) Play() *Play {
	return a.play
}

func (a *Application) ActiveView() ViewName {
	return a.activeView
}

func (a *Application) SetActiveView(view ViewName) {
	a.activeView = view
	a.SignalChange(ApplicationActiveViewChange)
}
