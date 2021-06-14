package main

import (
	"log"

	"github.com/mokiat/lacking/app"
	glfwapp "github.com/mokiat/lacking/framework/glfw/app"
	glgraphics "github.com/mokiat/lacking/framework/opengl/game/graphics"
	glui "github.com/mokiat/lacking/framework/opengl/ui"
	"github.com/mokiat/lacking/ui"
	_ "github.com/mokiat/lacking/ui/standard"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/intro"
)

func main() {
	cfg := glfwapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetCursorVisible(false)
	cfg.SetIcon("resources/icon.png")

	graphicsEngine := glgraphics.NewEngine()
	gameController := game.NewController(graphicsEngine)

	uiGLGraphics := glui.NewGraphics()
	uiController := ui.NewController(ui.FileResourceLocator{}, uiGLGraphics, intro.Config{
		GameController: gameController,
	})

	controller := app.NewLayeredController(gameController, uiController)

	log.Println("running application")
	if err := glfwapp.Run(cfg, controller); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}
