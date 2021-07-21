package store

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
)

const (
	ViewIntro int = iota
	ViewPlay
)

func ApplicationReducer() (co.Reducer, interface{}) {
	return func(store *co.Store, action interface{}) interface{} {
			var value Application
			store.Inject(&value)

			switch action := action.(type) {
			case ChangeViewAction:
				value.MainViewIndex = action.ViewIndex
			case SetGameDataAction:
				value.GameData = action.GameData
			}
			return value
		}, Application{
			MainViewIndex: ViewIntro,
		}
}

type Application struct {
	MainViewIndex int
	GameData      *scene.Data
}

type ChangeViewAction struct {
	ViewIndex int
}

type SetGameDataAction struct {
	GameData *scene.Data
}
