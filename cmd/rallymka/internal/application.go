package internal

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/global"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/intro"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ui/play"
)

func BootstrapApplication(window *ui.Window, gfxEngine graphics.Engine, gameController *game.Controller) {
	co.Initialize(window, co.New(co.StoreProvider, func() {
		co.WithData(co.StoreProviderData{
			Entries: []co.StoreProviderEntry{
				co.NewStoreProviderEntry(store.ApplicationReducer()),
			},
		})

		co.WithChild("app", co.New(Application, func() {
			co.WithContext(global.Context{
				GFXEngine:      gfxEngine,
				GameController: gameController,
			})
		}))
	}))
}

type ApplicationData = mat.SwitchData

var Application = co.Connect(co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	return co.New(mat.Switch, func() {
		co.WithData(props.Data())

		co.WithChild("intro", co.New(intro.View, func() {}))
		co.WithChild("play", co.New(play.View, func() {}))
	})

})), co.ConnectMapping{
	Data: func(props co.Properties) interface{} {
		var appStore store.Application
		co.InjectStore(&appStore)

		return ApplicationData{
			VisibleChildIndex: appStore.MainViewIndex,
		}
	},
})
