package scene

import (
	"math/rand"
	"time"

	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/log"
)

const (
	modelIDSUV        = "eaeb7483-7271-441f-a470-c0a8fa225161"
	levelIDForest     = "884e6395-2300-47bb-9916-b80e3dc0e086"
	levelIDHighway    = "acf21108-47ad-44ef-ba21-bf5473bfbaa0"
	levelIDPlayground = "9ca25b5c-ffa0-4224-ad80-a3c4d67930b7"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewData(resourceSet *game.ResourceSet) *Data {
	return &Data{
		resourceSet: resourceSet,
	}
}

type Data struct {
	resourceSet *game.ResourceSet

	loadOutcome async.Outcome

	Level          *game.SceneDefinition
	SceneModel     *game.ModelDefinition
	CarModel       *game.ModelDefinition
	PlatformModel  *game.ModelDefinition
	CharacterModel *game.ModelDefinition
}

func (d *Data) Request() async.Outcome {
	var levelID string
	switch rand.Intn(2) {
	case 0:
		levelID = levelIDForest
	case 1:
		levelID = levelIDHighway
	default:
		levelID = levelIDPlayground
	}
	levelID = levelIDPlayground

	levelOutcome := async.NewOutcome()
	levelPromise := d.resourceSet.OpenScene(levelID)
	levelPromise.OnSuccess(func(model *game.SceneDefinition) {
		levelOutcome.Record(async.Result{
			Value: model,
		})
	})
	levelPromise.OnError(func(err error) {
		log.Info("Level loading failed: %v", err)
		levelOutcome.Record(async.Result{
			Err: err,
		})
	})

	sceneOutcome := async.NewOutcome()
	promise := d.resourceSet.OpenModel("0cb3f13f-3bfe-4913-afe1-9f7cd6bf4c15")
	promise.OnSuccess(func(model *game.ModelDefinition) {
		sceneOutcome.Record(async.Result{
			Value: model,
		})
	})
	promise.OnError(func(err error) {
		log.Info("Scene loading failed: %v", err)
		sceneOutcome.Record(async.Result{
			Err: err,
		})
	})

	carOutcome := async.NewOutcome()
	carPromise := d.resourceSet.OpenModel(modelIDSUV)
	carPromise.OnSuccess(func(model *game.ModelDefinition) {
		carOutcome.Record(async.Result{
			Value: model,
		})
	})
	carPromise.OnError(func(err error) {
		log.Info("SUV loading failed: %v", err)
		carOutcome.Record(async.Result{
			Err: err,
		})
	})

	platformOutcome := async.NewOutcome()
	platformPromise := d.resourceSet.OpenModel("59f0b554-b70c-497a-9771-13d1d3a2f644")
	platformPromise.OnSuccess(func(model *game.ModelDefinition) {
		platformOutcome.Record(async.Result{
			Value: model,
		})
	})
	platformPromise.OnError(func(err error) {
		log.Info("Platform loading failed: %v", err)
		platformOutcome.Record(async.Result{
			Err: err,
		})
	})

	characterOutcome := async.NewOutcome()
	characterPromise := d.resourceSet.OpenModel("847244e5-2797-479a-9ea6-e66875458da6")
	characterPromise.OnSuccess(func(model *game.ModelDefinition) {
		characterOutcome.Record(async.Result{
			Value: model,
		})
	})
	characterPromise.OnError(func(err error) {
		log.Info("Character loading failed: %v", err)
		characterOutcome.Record(async.Result{
			Err: err,
		})
	})

	d.loadOutcome = async.NewCompositeOutcome(
		levelOutcome.OnSuccess(func(value any) {
			d.Level = value.(*game.SceneDefinition)
		}),
		sceneOutcome.OnSuccess(func(value any) {
			d.SceneModel = value.(*game.ModelDefinition)
		}),
		carOutcome.OnSuccess(func(value any) {
			d.CarModel = value.(*game.ModelDefinition)
		}),
		platformOutcome.OnSuccess(func(value any) {
			d.PlatformModel = value.(*game.ModelDefinition)
		}),
		characterOutcome.OnSuccess(func(value any) {
			d.CharacterModel = value.(*game.ModelDefinition)
		}),
	)
	return d.loadOutcome
}

func (d *Data) Dismiss() {
	d.resourceSet.Delete()
}

func (d *Data) IsAvailable() bool {
	return d.loadOutcome.IsAvailable()
}
