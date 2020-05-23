package main

import (
	"log"
	"time"

	"github.com/mokiat/lacking/game"
	rallygame "github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
)

func main() {
	log.Println("game started")
	app := game.NewApp(game.AppConfig{
		WindowTitle:        "RallyMKA",
		WindowWidth:        1024,
		WindowHeight:       576,
		WindowHideCursor:   true,
		WindowVSync:        true,
		UpdateLoopInterval: 16 * time.Millisecond,
	})
	controller := rallygame.NewController()
	if err := app.Run(controller); err != nil {
		log.Fatalf("game crashed: %v", err)
	}
	log.Println("game closed")
}
