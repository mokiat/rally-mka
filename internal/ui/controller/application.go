package controller

import (
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

func NewApplication(appModel *model.Application) *Application {
	return &Application{
		appModel: appModel,
	}
}

type Application struct {
	appModel *model.Application
}

func (a *Application) Reduce(act mvc.Action) bool {
	switch act := act.(type) {
	case action.ChangeView:
		a.appModel.SetActiveView(act.ViewName)
		return true
	case action.SetGameData:
		a.appModel.SetGameData(act.GameData)
		return true
	default:
		return false
	}
}
