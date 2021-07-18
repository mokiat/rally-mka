package internal

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	t "github.com/mokiat/lacking/ui/template"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/intro"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/play"
)

func BootstrapApplication(window *ui.Window, gfxEngine graphics.Engine, gameController *game.Controller) {
	t.InitGlobalState(store.CreateApplicationState(gfxEngine, gameController))
	t.Initialize(window, t.New(Application, func() {}))
}

var Application = t.Connect(t.ShallowCached(t.Plain(func(props t.Properties) t.Instance {
	return t.New(mat.Switch, func() {
		t.WithData(props.Data())

		t.WithChild("intro", t.New(intro.View, func() {}))
		t.WithChild("play", t.New(play.View, func() {}))
	})
})), func(props t.Properties, state *t.ReducedState) (interface{}, interface{}) {
	var appState store.Application
	state.Inject(&appState)

	return mat.SwitchData{
		VisibleChildIndex: appState.MainViewIndex,
	}, nil
})
