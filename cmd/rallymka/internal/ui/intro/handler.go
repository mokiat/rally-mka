package intro

import (
	"fmt"

	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/play"
)

type Config struct {
	GameController *game.Controller
}

func (c Config) SetupView(view *ui.View) error {
	template, err := view.Context().OpenTemplate("resources/ui/intro/view.xml")
	if err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}
	rootControl, err := view.Context().InstantiateTemplate(template, nil)
	if err != nil {
		return fmt.Errorf("failed to instantiate template: %w", err)
	}
	view.SetRoot(rootControl)
	view.SetHandler(&Handler{
		gameController: c.GameController,
	})
	return nil
}

type Handler struct {
	gameController *game.Controller
}

func (h *Handler) OnCreate(view *ui.View) {
	gameData := scene.NewData(h.gameController.Registry(), h.gameController.GFXWorker())
	gameData.Request().OnSuccess(func(interface{}) {
		view.Context().Schedule(func() {
			view.Window().OpenView(ui.ViewModeReplace, play.Config{
				GameController: h.gameController,
				GameData:       gameData,
			})
		})
	})
}

func (h *Handler) OnShow(view *ui.View) {}

func (h *Handler) OnHide(view *ui.View) {}

func (h *Handler) OnDestroy(view *ui.View) {}
