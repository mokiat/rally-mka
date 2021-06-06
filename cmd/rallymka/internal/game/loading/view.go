package loading

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/graphics"
)

func NewView(gfxEngine graphics.Engine) *View {
	return &View{
		gfxEngine: gfxEngine,
	}
}

type View struct {
	gfxEngine graphics.Engine
	gfxScene  graphics.Scene
	camera    graphics.Camera

	loadingAngle sprec.Angle
}

func (v *View) Load(window app.Window, cb func()) {
	v.gfxScene = v.gfxEngine.CreateScene()
	v.camera = v.gfxScene.CreateCamera()

	window.Schedule(func() error {
		cb()
		return nil
	})
}

func (v *View) Unload(window app.Window) {
	v.gfxScene.Delete()
}

func (v *View) Open(window app.Window) {}

func (v *View) Close(window app.Window) {}

func (v *View) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	return false
}

func (v *View) Update(window app.Window, elapsedSeconds float32) {
	v.loadingAngle += sprec.Degrees(elapsedSeconds * 180.0)
	cs := sprec.Abs(sprec.Cos(v.loadingAngle))
	sn := sprec.Abs(sprec.Sin(v.loadingAngle))
	v.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(cs, sn, 0.0))
}

func (v *View) Render(window app.Window, width, height int) {
	v.gfxScene.Render(graphics.NewViewport(0, 0, width, height), v.camera)
}
