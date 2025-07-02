package data

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadHomeData(engine *game.Engine, resourceSet *game.ResourceSet) async.Promise[*HomeData] {
	var data HomeData
	return async.InjectionPromise(async.JoinOperations(
		resourceSet.FetchResource("HomeScreen.dat", &data.Scene),
		resourceSet.FetchResource("Vehicle.dat", &data.Vehicle),
	), &data)
}

type HomeData struct {
	Scene   *game.ModelTemplate
	Vehicle *game.ModelTemplate
}
