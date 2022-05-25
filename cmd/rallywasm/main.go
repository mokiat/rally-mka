//go:build js && wasm

package main

import (
	"os"

	jsapp "github.com/mokiat/lacking-js/app"
	jsgame "github.com/mokiat/lacking-js/game"
	jsrender "github.com/mokiat/lacking-js/render"
	jsui "github.com/mokiat/lacking-js/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/resource"
	"github.com/mokiat/rally-mka/internal"
	"github.com/mokiat/rally-mka/internal/game"
	"github.com/mokiat/rally-mka/resources"
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
	registry := asset.NewWebRegistry(".")
	resourceLocator := mat.WrappedResourceLocator(resource.NewFSLocator(resources.UI))
	renderAPI := jsrender.NewAPI()
	graphicsEngine := graphics.NewEngine(renderAPI, jsgame.NewShaderCollection())
	gameController := game.NewController(registry, graphicsEngine)
	uiCfg := ui.NewConfig(resourceLocator, renderAPI, jsui.NewShaderCollection())
	uiController := ui.NewController(uiCfg, func(w *ui.Window) {
		internal.BootstrapApplication(w, gameController)
	})

	cfg := jsapp.NewConfig("screen")
	cfg.AddGLExtension("EXT_color_buffer_float")
	cfg.AddGLExtension("EXT_float_blend")

	controller := app.NewLayeredController(
		gameController,
		uiController,
	)
	return jsapp.Run(cfg, controller)
}
