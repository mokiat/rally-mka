package simulation

import (
	"time"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
)

func NewView(registry *resource.Registry, gfxWorker *async.Worker) *View {
	return &View{
		gfxWorker: gfxWorker,
		gameData:  scene.NewData(registry, gfxWorker),
		stage:     scene.NewStage(),
	}
}

type View struct {
	gfxWorker *async.Worker
	gameData  *scene.Data
	stage     *scene.Stage

	freezeFrame bool
}

func (v *View) Load(window app.Window, cb func()) {
	loadOutcome := async.NewOutcome()
	go func() {
		v.gameData.Request().Wait()
		v.stage.Init(v.gfxWorker, v.gameData)
		loadOutcome.Record(async.Result{})
	}()
	go func() {
		loadOutcome.Wait()
		window.Schedule(func() error {
			cb()
			return nil
		})
	}()
}

func (v *View) Unload(window app.Window) {
	v.gameData.Dismiss()
}

func (v *View) Open(window app.Window) {}

func (v *View) Close(window app.Window) {}

func (v *View) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	if event.Code == app.KeyCodeF {
		switch event.Type {
		case app.KeyboardEventTypeKeyDown:
			v.freezeFrame = true
			return true
		case app.KeyboardEventTypeKeyUp:
			v.freezeFrame = false
			return true
		}
	}
	return v.stage.OnKeyboardEvent(event)
}

func (v *View) Update(window app.Window, elapsedTime time.Duration) {
	if !v.freezeFrame {
		v.stage.Update(window, elapsedTime)
	}
}

func (v *View) Render(window app.Window, width, height int, pipeline *graphics.Pipeline) {
	v.stage.Render(width, height, pipeline)
}
