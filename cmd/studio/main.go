package main

import (
	"os"

	"github.com/mokiat/lacking-studio/studio"
	"github.com/mokiat/lacking/debug/log"
)

func main() {
	if err := studio.Run(); err != nil {
		log.Error("Error: %v", err)
		os.Exit(1)
	}
}
