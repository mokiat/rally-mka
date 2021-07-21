package intro

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/global"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
)

var View = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var context global.Context
	co.InjectContext(&context)

	var logoImg ui.Image
	co.UseState(func() interface{} {
		return co.OpenImage("resources/ui/intro/logo.png")
	}).Inject(&logoImg)

	co.Once(func() {
		gameData := scene.NewData(
			context.GameController.Registry(),
			context.GameController.GFXWorker(),
		)
		gameData.Request().OnSuccess(func(interface{}) {
			co.Schedule(func() {
				co.Dispatch(store.SetGameDataAction{
					GameData: gameData,
				})
				co.Dispatch(store.ChangeViewAction{
					ViewIndex: store.ViewPlay,
				})
			})
		})
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("logo-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				BackgroundColor: ui.Transparent(),
				Image:           logoImg,
				Mode:            mat.ImageModeFit,
			})
			co.WithLayoutData(mat.AnchorLayoutData{
				Width:                    optional.NewInt(512),
				Height:                   optional.NewInt(128),
				HorizontalCenter:         optional.NewInt(0),
				HorizontalCenterRelation: mat.RelationCenter,
				VerticalCenter:           optional.NewInt(0),
				VerticalCenterRelation:   mat.RelationCenter,
			})
		}))
	})
}))
