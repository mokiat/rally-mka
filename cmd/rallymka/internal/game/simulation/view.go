package simulation

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry) *View {
	return &View{
		gameData: scene.NewData(registry),
		camera:   scene.NewCamera(),
		stage:    scene.NewStage(),
	}
}

type View struct {
	gameData *scene.Data
	camera   *scene.Camera
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
	v.stage.Init(v.gameData)
}

func (v *View) Close() {
	// TODO: Erase stage
}

func (v *View) Resize(width, height int) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	v.camera.SetProjectionMatrix(math.PerspectiveMat4x4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 300.0,
	))
}

func (v *View) Update(elapsedSeconds float32, actions input.ActionSet) {
	if !actions.FreezeFrame {
		v.stage.Update(elapsedSeconds, v.camera, scene.CarInput{
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
