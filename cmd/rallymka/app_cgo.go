//go:build !js

package main

import (
	"fmt"

	glapp "github.com/mokiat/lacking-native/app"
	glgame "github.com/mokiat/lacking-native/game"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	gameui "github.com/mokiat/rally-mka/internal/ui"
	"github.com/mokiat/rally-mka/resources"
)

func runApplication() error {
	registryStorage, err := asset.NewFSStorage("./assets")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	registryFormatter := asset.NewBlobFormatter()

	registry, err := asset.NewRegistry(registryStorage, registryFormatter)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))

	gameController := game.NewController(registry, glgame.NewShaderCollection(), glgame.NewShaderBuilder())
	uiController := ui.NewController(locator, glui.NewShaderCollection(), func(w *ui.Window) {
		gameui.BootstrapApplication(w, gameController)
	})

	cfg := glapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetFullscreen(false)
	cfg.SetMaximized(true)
	cfg.SetMinSize(1024, 576)
	cfg.SetVSync(true)
	cfg.SetIcon("ui/images/icon.png")
	cfg.SetLocator(locator)
	cfg.SetAudioEnabled(false)
	return glapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
