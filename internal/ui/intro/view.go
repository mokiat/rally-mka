package intro

import (
	"fmt"

	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
	"github.com/mokiat/rally-mka/internal/global"
	"github.com/mokiat/rally-mka/internal/scene"
	"github.com/mokiat/rally-mka/internal/store"
)

var View = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	context := co.GetContext[global.Context]()

	co.Once(func() {
		co.Window().SetCursorVisible(false)
	})

	co.Defer(func() {
		co.Window().SetCursorVisible(true)
	})

	co.Once(func() {
		gameData := scene.NewData(
			context.GameController.Registry(),
		)
		gameData.Request().OnSuccess(func(interface{}) {
			co.Schedule(func() {
				co.Dispatch(store.SetGameDataAction{
					GameData: gameData,
				})
				co.Dispatch(store.ChangeViewAction{
					ViewIndex: store.ViewHome,
				})
			})
		}).OnError(func(err error) {
			panic(fmt.Errorf("failed to load assets: %w", err))
		})
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("logo-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				BackgroundColor: optional.Value(ui.Transparent()),
				Image:           co.OpenImage("resources/ui/images/logo.png"),
				Mode:            mat.ImageModeFit,
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            optional.Value(512),
				Height:           optional.Value(128),
				HorizontalCenter: optional.Value(0),
				VerticalCenter:   optional.Value(0),
			})
		}))
	})
}))
