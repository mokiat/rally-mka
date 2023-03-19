package data

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadPlayData(engine *game.Engine, resourceSet *game.ResourceSet) game.Promise[*PlayData] {
	scenePromise := resourceSet.OpenSceneByName("Forest")
	vehiclePromise := resourceSet.OpenModelByName("SUV")

	result := async.NewPromise[*PlayData]()
	go func() {
		var data PlayData
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
	Scene   *game.SceneDefinition
	Vehicle *game.ModelDefinition
}
