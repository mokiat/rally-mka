package main

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"time"

	"github.com/mokiat/rally-mka/internal/game/level"
)

func main() {
	if err := runApp(); err != nil {
		slog.Error("Crashed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func runApp() error {
	random := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

	generator := level.NewGenerator(level.GeneratorConfig{
		Random: random,
		Size:   9,
	})

	worstDuration := time.Duration(0)
	var board *level.Board
	for range 100 {
		startTime := time.Now()
		board = generator.Generate()
		elapsedTime := time.Since(startTime)
		if elapsedTime > worstDuration {
			worstDuration = elapsedTime
			slog.Info("Generated board", slog.String("duration", elapsedTime.String()))
			data, err := level.SerializeBoard(board)
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println(string(data))
			fmt.Println()
		}
	}

	return nil
}
