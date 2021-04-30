package scene

import (
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/resource"
)

func NewData(registry *resource.Registry, gfxWorker *async.Worker) *Data {
	return &Data{
		registry:  registry,
		gfxWorker: gfxWorker,
	}
}

type Data struct {
	registry    *resource.Registry
	gfxWorker   *async.Worker
	loadOutcome async.Outcome

	CarModel *resource.Model
	Level    *resource.Level
}

func (d *Data) Request() async.Outcome {
	d.loadOutcome = async.NewCompositeOutcome(
		d.registry.LoadModel("suv").OnSuccess(resource.InjectModel(&d.CarModel)),
		d.registry.LoadLevel("forest").OnSuccess(resource.InjectLevel(&d.Level)),
	)
	return d.loadOutcome
}

func (d *Data) Dismiss() {
	d.registry.UnloadModel(d.CarModel)
	d.registry.UnloadLevel(d.Level)
}

func (d *Data) IsAvailable() bool {
	return d.loadOutcome.IsAvailable()
}
