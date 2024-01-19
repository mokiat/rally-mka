package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	ViewNameIntro    ViewName = "intro"
	ViewNameHome     ViewName = "home"
	ViewNamePlay     ViewName = "play"
	ViewNameLoading  ViewName = "loading"
	ViewNameLicenses ViewName = "licenses"
	ViewNameCredits  ViewName = "credits"
)

type ViewName = string

func NewApplication(eventBus *mvc.EventBus) *Application {
	return &Application{
		eventBus: eventBus,

		activeView: ViewNameIntro,
	}
}

type Application struct {
	eventBus *mvc.EventBus

	activeView ViewName
}

func (a *Application) ActiveView() ViewName {
	return a.activeView
}

func (a *Application) SetActiveView(view ViewName) {
	if view != a.activeView {
		a.activeView = view
		a.eventBus.Notify(ActiveViewChangedEvent{})
	}
}

type ActiveViewChangedEvent struct{}
