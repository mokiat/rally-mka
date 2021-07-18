package store

import (
	"github.com/mokiat/lacking/game/graphics"
	t "github.com/mokiat/lacking/ui/template"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
)

const (
	ViewIntro int = iota
	ViewPlay
)

func CreateApplicationState(gfxEngine graphics.Engine, gameController *game.Controller) *t.ReducedState {
	return t.NewReducedState(func(state *t.ReducedState, action interface{}) interface{} {
		if state == nil {
			return Application{
				GFXEngine:      gfxEngine,
				GameController: gameController,
				MainViewIndex:  ViewIntro,
			}
		}
		var appState Application
		state.Inject(&appState)
		switch action := action.(type) {
		case ChangeViewAction:
			appState.MainViewIndex = action.ViewIndex
		case SetGameDataAction:
			appState.GameData = action.GameData
		}
		return appState
	})
}

type Application struct {
	GFXEngine      graphics.Engine
	GameController *game.Controller
	MainViewIndex  int
	GameData       *scene.Data
}

type ChangeViewAction struct {
	ViewIndex int
}

type SetGameDataAction struct {
	GameData *scene.Data
}
