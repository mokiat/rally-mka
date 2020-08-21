package simulation

import (
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/input"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
)

func NewView(registry *resource.Registry, gfxWorker *async.Worker) *View {
	return &View{
		gameData: scene.NewData(registry, gfxWorker),
		stage:    scene.NewStage(gfxWorker),
	}
}

type View struct {
	gameData *scene.Data
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

func (v *View) Update(ctx game.UpdateContext) {
	if !ctx.Keyboard.IsPressed(input.KeyF) {
		v.stage.Update(ctx)
	}
}

func (v *View) Render(ctx game.RenderContext) {
	v.stage.Render(ctx)
}
