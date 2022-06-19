package main

import (
	"fmt"
	"os"

	glapp "github.com/mokiat/lacking-gl/app"
	glgame "github.com/mokiat/lacking-gl/game"
	glrender "github.com/mokiat/lacking-gl/render"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/resource"
	"github.com/mokiat/rally-mka/internal/game"
	gameui "github.com/mokiat/rally-mka/internal/ui"
)

func main() {
	log.Info("Started")
	if err := runApplication(); err != nil {
		log.Error("Crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Stopped")
}

func runApplication() error {
	cfg := glapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetIcon("resources/icon.png")

	registry, err := asset.NewDirRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	renderAPI := glrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, glgame.NewShaderCollection())
	gameController := game.NewController(registry, graphicsEngine)
	resourceLocator := mat.WrappedResourceLocator(resource.NewFileLocator("./resources"))
	uiCfg := ui.NewConfig(resourceLocator, renderAPI, glui.NewShaderCollection())
	uiController := ui.NewController(uiCfg, func(w *ui.Window) {
		gameui.BootstrapApplication(w, gameController)
	})

	controller := app.NewLayeredController(
		gameController,
		uiController,
	)
	return glapp.Run(cfg, controller)
}
