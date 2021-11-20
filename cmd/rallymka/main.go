package main

import (
	"fmt"
	"log"

	"github.com/mokiat/lacking/app"
	glfwapp "github.com/mokiat/lacking/framework/glfw/app"
	glgraphics "github.com/mokiat/lacking/framework/opengl/game/graphics"
	glui "github.com/mokiat/lacking/framework/opengl/ui"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
)

func main() {
	log.Println("running application")
	if err := runApplication(); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}

func runApplication() error {
	cfg := glfwapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetIcon("resources/icon.png")

	registry, err := asset.NewDirRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	graphicsEngine := glgraphics.NewEngine()
	gameController := game.NewController(registry, graphicsEngine)

	resourceLocator := ui.NewFileResourceLocator(".")
	uiGLGraphics := glui.NewGraphics()
	uiController := ui.NewController(resourceLocator, uiGLGraphics, func(w *ui.Window) {
		internal.BootstrapApplication(w, graphicsEngine, gameController)
	})

	controller := app.NewLayeredController(gameController, uiController)
	return glfwapp.Run(cfg, controller)
}
