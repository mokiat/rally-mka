//go:build js

package main

import (
	"fmt"

	jsapp "github.com/mokiat/lacking-js/app"
	jsgame "github.com/mokiat/lacking-js/game"
	jsui "github.com/mokiat/lacking-js/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	gameui "github.com/mokiat/rally-mka/internal/ui"
	"github.com/mokiat/rally-mka/resources"
)

func runApplication() error {
	registry, err := asset.NewWebRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}
	resourceLocator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))
	gameController := game.NewController(registry, jsgame.NewShaderCollection(), nil) // FIXME
	uiController := ui.NewController(resourceLocator, jsui.NewShaderCollection(), func(w *ui.Window) {
		gameui.BootstrapApplication(w, gameController)
	})

	cfg := jsapp.NewConfig("screen")
	cfg.AddGLExtension("EXT_color_buffer_float")
	cfg.SetFullscreen(false)
	cfg.SetAudioEnabled(false)
	return jsapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
