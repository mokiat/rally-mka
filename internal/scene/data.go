package scene

import (
	"math/rand"
	"time"

	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/resource"
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
	var levelID string
	switch rand.Intn(2) {
	case 0:
		levelID = levelIDForest
	case 1:
		levelID = levelIDHighway
	default:
		levelID = levelIDPlayground
	}
	d.loadOutcome = async.NewCompositeOutcome(
		d.registry.LoadModel(modelIDSUV).OnSuccess(resource.InjectModel(&d.CarModel)),
		d.registry.LoadLevel(levelID).OnSuccess(resource.InjectLevel(&d.Level)), // Forest
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
