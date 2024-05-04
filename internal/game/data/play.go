package data

import (
	"errors"
	"fmt"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

type Controller string

const (
	ControllerKeyboard Controller = "keyboard"
	ControllerMouse    Controller = "mouse"
	ControllerGamepad  Controller = "gamepad"
)

type Environment string

const (
	EnvironmentDay   Environment = "day"
	EnvironmentNight Environment = "night"
)

func LoadPlayData(engine *game.Engine, resourceSet *game.ResourceSet, environment Environment, controller Controller) async.Promise[*PlayData] {
	var backgroundName string
	switch environment {
	case EnvironmentDay:
		backgroundName = "Forest-Day"
	case EnvironmentNight:
		backgroundName = "Forest-Night"
	default:
		panic(fmt.Errorf("unknown environment %q", environment))
	}

	backgroundPromise := resourceSet.OpenModelByName(backgroundName)
	scenePromise := resourceSet.OpenModelByName("Forest")
	vehiclePromise := resourceSet.OpenModelByName("SUV")

	promise := async.NewPromise[*PlayData]()
	go func() {
		var data PlayData
		data.Environment = environment
		data.Controller = controller
		err := errors.Join(
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

	Environment Environment
	Controller  Controller
}
