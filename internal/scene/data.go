package scene

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

func NewData(engine *game.Engine, resourceSet *game.ResourceSet) *Data {
	return &Data{
		engine:      engine,
		resourceSet: resourceSet,
	}
}

type Data struct {
	engine      *game.Engine
	resourceSet *game.ResourceSet

	loadPromise async.Promise[struct{}]

	Level    *game.SceneDefinition
	CarModel *game.ModelDefinition
}

func (d *Data) Request() game.Promise[struct{}] {
	levelPromise := d.resourceSet.OpenSceneByName("Forest")
	carPromise := d.resourceSet.OpenModelByName("SUV")

	d.loadPromise = async.NewPromise[struct{}]()
	go func() {
		err := firstErr(
			levelPromise.Inject(&d.Level),
			carPromise.Inject(&d.CarModel),
		)
		if err != nil {
			d.loadPromise.Fail(err)
		} else {
			d.loadPromise.Deliver(struct{}{})
		}
	}()
	return game.SafePromise(d.loadPromise, d.engine)
}

func (d *Data) Dismiss() {
	d.resourceSet.Delete()
}

func (d *Data) IsAvailable() bool {
	return d.loadPromise.Ready()
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
