package simulation

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry, gfxWorker *graphics.Worker) *View {
	return &View{
		gameData: scene.NewData(registry, gfxWorker),
		camera:   ecs.NewCamera(),
		stage:    scene.NewStage(gfxWorker),
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

func (v *View) Update(ctx game.UpdateContext) {
	if !ctx.Keyboard.IsPressed(input.KeyF) {
		v.stage.Update(ctx, v.camera)
	}
}

func (v *View) Render(ctx game.RenderContext) {
	width := ctx.WindowSize.Width
	height := ctx.WindowSize.Height

	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	v.camera.SetProjectionMatrix(sprec.PerspectiveMat4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 300.0,
	))
	v.stage.Resize(width, height)

	v.stage.Render(ctx.GFXPipeline, v.camera)
}
