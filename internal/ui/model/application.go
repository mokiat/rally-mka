package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	ViewNameIntro    ViewName = "intro"
	ViewNameError    ViewName = "error"
	ViewNameHome     ViewName = "home"
	ViewNamePlay     ViewName = "play"
	ViewNameLoading  ViewName = "loading"
	ViewNameLicenses ViewName = "licenses"
	ViewNameCredits  ViewName = "credits"
)

type ViewName = string

func NewApplicationModel(eventBus *mvc.EventBus) *ApplicationModel {
	return &ApplicationModel{
		eventBus: eventBus,

		activeView: ViewNameIntro,
	}
}

type ApplicationModel struct {
	eventBus *mvc.EventBus

	activeView ViewName
}

func (a *ApplicationModel) ActiveView() ViewName {
	return a.activeView
}

func (a *ApplicationModel) SetActiveView(view ViewName) {
	if view != a.activeView {
		a.activeView = view
		a.eventBus.Notify(ActiveViewChangedEvent{})
	}
}

type ActiveViewChangedEvent struct{}
