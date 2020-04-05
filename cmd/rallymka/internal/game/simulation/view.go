package simulation

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry, gfxWorker *graphics.Worker) *View {
	return &View{
		gameData: scene.NewData(registry, gfxWorker),
		camera:   ecs.NewCamera(),
		stage:    scene.NewStage(),
	}
}

type View struct {
	gameData *scene.Data
	camera   *ecs.Camera
	stage    *scene.Stage
}

func (v *View) Load() {
	v.gameData.Request()
}

func (v *View) Unload() {
	v.gameData.Dismiss()
}

func (v *View) IsAvailable() bool {
	return v.gameData.IsAvailable()
}

func (v *View) Open() {
	v.stage.Init(v.gameData, v.camera)
}

func (v *View) Close() {
	// TODO: Erase stage
}

func (v *View) Resize(width, height int) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	v.camera.SetProjectionMatrix(sprec.PerspectiveMat4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 300.0,
	))
	v.stage.Resize(width, height)
}

func (v *View) Update(elapsedTime time.Duration, actions input.ActionSet) {
	if !actions.FreezeFrame {
		v.stage.Update(elapsedTime, v.camera, ecs.CarInput{
			Forward:   actions.Forward,
			Backward:  actions.Backward,
			TurnLeft:  actions.Left,
			TurnRight: actions.Right,
			Handbrake: actions.Handbrake,
		})
	}
}

func (v *View) Render(pipeline *graphics.Pipeline) {
	v.stage.Render(pipeline, v.camera)
}
