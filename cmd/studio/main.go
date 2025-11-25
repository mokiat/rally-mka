package main

import (
	"log/slog"
	"os"

	"github.com/mokiat/lacking-studio/studio"
	"github.com/mokiat/lacking/game/asset/conv"
	"github.com/mokiat/lacking/game/asset/dsl"
)

var _ = dsl.Use(
	conv.NewModelConverter(),
)

func main() {
	if err := studio.Run(); err != nil {
		slog.Error("Crashed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
