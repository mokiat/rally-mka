package data

import (
	"cmp"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadHomeData(engine *game.Engine, resourceSet *game.ResourceSet) async.Promise[*HomeData] {
	backgroundPromise := resourceSet.OpenModelByName("Home-Screen")
	scenePromise := resourceSet.OpenModelByName("HomeScreen")
	vehiclePromise := resourceSet.OpenModelByName("Vehicle")

	promise := async.NewPromise[*HomeData]()
	go func() {
		var data HomeData
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

type HomeData struct {
	Background *game.ModelDefinition
	Scene      *game.ModelDefinition
	Vehicle    *game.ModelDefinition
}
