//go:build js

package main

import (
	"fmt"

	jsapp "github.com/mokiat/lacking-js/app"
	jsgame "github.com/mokiat/lacking-js/game"
	jsrender "github.com/mokiat/lacking-js/render"
	jsui "github.com/mokiat/lacking-js/ui"
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
	registry, err := asset.NewWebRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}
	resourceLocator := mat.WrappedResourceLocator(resource.NewFSLocator(resources.UI))
	renderAPI := jsrender.NewAPI()
	gameController := game.NewController(registry, renderAPI, jsgame.NewShaderCollection())
	uiCfg := ui.NewConfig(resourceLocator, renderAPI, jsui.NewShaderCollection())
	uiController := ui.NewController(uiCfg, func(w *ui.Window) {
		gameui.BootstrapApplication(w, gameController)
	})

	cfg := jsapp.NewConfig("screen")
	cfg.AddGLExtension("EXT_color_buffer_float")
	cfg.SetFullscreen(false)
	return jsapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
