package data

import (
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

func LoadPlayData(engine *game.Engine, resourceSet *game.ResourceSet, environment Environment, controller Controller) game.Promise[*PlayData] {
	var sceneName string
	switch environment {
	case EnvironmentDay:
		sceneName = "Forest-Day"
	case EnvironmentNight:
		sceneName = "Forest-Night"
	default:
		panic(fmt.Errorf("unknown environment %q", environment))
	}

	scenePromise := resourceSet.OpenSceneByName(sceneName)
	vehiclePromise := resourceSet.OpenModelByName("SUV")

	result := async.NewPromise[*PlayData]()
	go func() {
		var data PlayData
		data.Environment = environment
		data.Controller = controller
		err := firstErr(
			scenePromise.Inject(&data.Scene),
			vehiclePromise.Inject(&data.Vehicle),
		)
		if err != nil {
			result.Fail(err)
		} else {
			result.Deliver(&data)
		}
	}()
	return game.SafePromise(result, engine)
}

type PlayData struct {
	Scene       *game.SceneDefinition
	Vehicle     *game.ModelDefinition
	Environment Environment
	Controller  Controller
}
