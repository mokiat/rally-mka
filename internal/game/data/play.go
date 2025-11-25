package data

import (
	"fmt"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
	"github.com/mokiat/rally-mka/internal/game/level"
)

type Input string

const (
	InputKeyboard Input = "keyboard"
	InputMouse    Input = "mouse"
	InputGamepad  Input = "gamepad"
)

type Lighting string

const (
	LightingDay   Lighting = "day"
	LightingNight Lighting = "night"
)

func LoadPlayData(engine *game.Engine, resourceSet *game.ResourceSet, lighting Lighting, input Input, board *level.Board) async.Promise[*PlayData] {
	var backgroundName string
	switch lighting {
	case LightingDay:
		backgroundName = "PlayScreen-Day.dat"
	case LightingNight:
		backgroundName = "PlayScreen-Night.dat"
	default:
		panic(fmt.Errorf("unknown lighting mode %q", lighting))
	}

	data := PlayData{
		Lighting: lighting,
		Input:    input,
		Board:    board,
	}
	return async.InjectionPromise(async.JoinOperations(
		resourceSet.FetchResource(backgroundName, &data.Background),
		resourceSet.FetchResource("Tiles.dat", &data.Scene),
		resourceSet.FetchResource("Vehicle.dat", &data.Vehicle),
	), &data)
}

type PlayData struct {
	Background *game.ModelTemplate
	Scene      *game.ModelTemplate
	Vehicle    *game.ModelTemplate

	Lighting Lighting
	Input    Input
	Board    *level.Board
}
