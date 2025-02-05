package data

import (
	"cmp"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadHomeData(engine *game.Engine, resourceSet *game.ResourceSet) async.Promise[*HomeData] {
	scenePromise := resourceSet.OpenModelByName("HomeScreen")
	vehiclePromise := resourceSet.OpenModelByName("Vehicle")

	promise := async.NewPromise[*HomeData]()
	go func() {
		var data HomeData
		err := cmp.Or(
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

type HomeData struct {
	Scene   *game.ModelDefinition
	Vehicle *game.ModelDefinition
}
