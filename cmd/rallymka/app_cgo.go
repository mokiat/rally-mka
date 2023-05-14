//go:build !js

package main

import (
	"fmt"

	glapp "github.com/mokiat/lacking-gl/app"
	glgame "github.com/mokiat/lacking-gl/game"
	glrender "github.com/mokiat/lacking-gl/render"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/resource"
	gameui "github.com/mokiat/rally-mka/internal/ui"
	"github.com/mokiat/rally-mka/resources"
)

func runApplication() error {
	registry, err := asset.NewDirRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}
	locator := mat.WrappedResourceLocator(resource.NewFSLocator(resources.UI))

	renderAPI := glrender.NewAPI()
	gameController := game.NewController(registry, renderAPI, glgame.NewShaderCollection())
	uiCfg := ui.NewConfig(locator, renderAPI, glui.NewShaderCollection())
	uiController := ui.NewController(uiCfg, func(w *ui.Window) {
		gameui.BootstrapApplication(w, gameController)
	})

	cfg := glapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetFullscreen(true)
	cfg.SetMaximized(false)
	cfg.SetMinSize(1024, 576)
	cfg.SetVSync(true)
	cfg.SetIcon("ui/images/icon.png")
	cfg.SetMaximized(true)
	cfg.SetLocator(locator)
	return glapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
