package main

import (
	"log"
	"runtime"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	log.Println("starting application")

	app := game.Application{}
	if err := app.Run(); err != nil {
		log.Fatalf("application crashed: %v", err)
	}

	log.Println("application closed")
}
