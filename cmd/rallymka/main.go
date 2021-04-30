package main

import (
	"log"

	glfwapp "github.com/mokiat/lacking/glfw/app"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
)

func main() {
	cfg := glfwapp.NewConfig("Rally MKA", 1024, 576)
	cfg.SetVSync(true)
	cfg.SetCursorVisible(false)
	cfg.SetIcon("resources/icon.png")

	controller := game.NewController()
	log.Println("running application")
	if err := glfwapp.Run(cfg, controller); err != nil {
		log.Fatalf("application error: %v", err)
	}
	log.Println("application closed")
}
