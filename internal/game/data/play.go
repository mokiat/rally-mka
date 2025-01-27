package data

import (
	"cmp"
	"fmt"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
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

func LoadPlayData(engine *game.Engine, resourceSet *game.ResourceSet, lighting Lighting, input Input) async.Promise[*PlayData] {
	var backgroundName string
	switch lighting {
	case LightingDay:
		backgroundName = "Forest-Day"
	case LightingNight:
		backgroundName = "Forest-Night"
	default:
		panic(fmt.Errorf("unknown lighting mode %q", lighting))
	}

	backgroundPromise := resourceSet.OpenModelByName(backgroundName)
	scenePromise := resourceSet.OpenModelByName("Tiles")
	vehiclePromise := resourceSet.OpenModelByName("Vehicle")

	promise := async.NewPromise[*PlayData]()
	go func() {
		var data PlayData
		data.Lighting = lighting
		data.Input = input
		err := cmp.Or(
			backgroundPromise.Inject(&data.Background),
			scenePromise.Inject(&data.Scene),
			vehiclePromise.Inject(&data.Vehicle),
		)
		if err != nil {
			promise.Fail(err)
		} else {
			promise.Deliver(&data)
		}
	}()
	return promise
}

type PlayData struct {
	Background *game.ModelDefinition
	Scene      *game.ModelDefinition
	Vehicle    *game.ModelDefinition

	Lighting Lighting
	Input    Input
}
