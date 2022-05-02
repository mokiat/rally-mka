package main

import (
	"fmt"
	"log"

	glapp "github.com/mokiat/lacking-gl/app"
	glgame "github.com/mokiat/lacking-gl/game"
	glrender "github.com/mokiat/lacking-gl/render"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/rally-mka/internal"
	"github.com/mokiat/rally-mka/internal/game"
)

func main() {
	log.Println("running application")
	if err := runApplication(); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
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
	resourceLocator := ui.NewFileResourceLocator(".")
	uiCfg := ui.NewConfig(resourceLocator, renderAPI, glui.NewShaderCollection())
	uiController := ui.NewController(uiCfg, func(w *ui.Window) {
		internal.BootstrapApplication(w, gameController)
	})

	controller := app.NewLayeredController(gameController, uiController)
	return glapp.Run(cfg, controller)
}
