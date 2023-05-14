package data

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func LoadHomeData(engine *game.Engine, resourceSet *game.ResourceSet) game.Promise[*HomeData] {
	scenePromise := resourceSet.OpenSceneByName("Home Screen")

	result := async.NewPromise[*HomeData]()
	go func() {
		var data HomeData
		err := firstErr(
			scenePromise.Inject(&data.Scene),
		)
		if err != nil {
			result.Fail(err)
		} else {
			result.Deliver(&data)
		}
	}()
	return game.SafePromise(result, engine)
}

type HomeData struct {
	Scene *game.SceneDefinition
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
