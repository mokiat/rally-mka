package scene

import (
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/resource"
)

func NewData(registry *resource.Registry) *Data {
	return &Data{
		registry: registry,
	}
}

type Data struct {
	registry    *resource.Registry
	loadOutcome async.Outcome

	CarModel *resource.Model
	Level    *resource.Level
}

func (d *Data) Request() async.Outcome {
	d.loadOutcome = async.NewCompositeOutcome(
		// SUV: eaeb7483-7271-441f-a470-c0a8fa225161
		d.registry.LoadModel("eaeb7483-7271-441f-a470-c0a8fa225161").OnSuccess(resource.InjectModel(&d.CarModel)),
		d.registry.LoadLevel("884e6395-2300-47bb-9916-b80e3dc0e086").OnSuccess(resource.InjectLevel(&d.Level)), // Forest
		// d.registry.LoadLevel("acf21108-47ad-44ef-ba21-bf5473bfbaa0").OnSuccess(resource.InjectLevel(&d.Level)), // Highway
		// d.registry.LoadLevel("9ca25b5c-ffa0-4224-ad80-a3c4d67930b7").OnSuccess(resource.InjectLevel(&d.Level)), // Playground
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
