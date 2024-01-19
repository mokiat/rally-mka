package data

import (
	"errors"

	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadHomeData(engine *game.Engine, resourceSet *game.ResourceSet) async.Promise[*HomeData] {
	scenePromise := resourceSet.OpenSceneByName("Home Screen")

	promise := async.NewPromise[*HomeData]()
	go func() {
		var data HomeData
		err := errors.Join(
			scenePromise.Inject(&data.Scene),
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
	Scene *game.SceneDefinition
}
