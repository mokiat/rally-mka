package model

import (
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/scene"
)

var (
	ChangeApplication = mvc.NewChange("application")
	ChangeActiveView  = mvc.SubChange(ChangeApplication, "active_view")
	ChangeGameData    = mvc.SubChange(ChangeApplication, "game_data")
)

const (
	ViewNameIntro ViewName = "intro"
	ViewNameHome  ViewName = "home"
	ViewNamePlay  ViewName = "play"
)

type ViewName string

func NewApplication() *Application {
	return &Application{
		Observable: mvc.NewObservable(),
		activeView: ViewNameIntro,
	}
}

type Application struct {
	mvc.Observable
	activeView ViewName
	gameData   *scene.Data
}

func (a *Application) ActiveView() ViewName {
	return a.activeView
}

func (a *Application) SetActiveView(view ViewName) {
	a.activeView = view
	a.SignalChange(ChangeActiveView)
}

func (a *Application) GameData() *scene.Data {
	return a.gameData
}

func (a *Application) SetGameData(data *scene.Data) {
	a.gameData = data
	a.SignalChange(ChangeGameData)
}
