package intro

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/global"
	"github.com/mokiat/rally-mka/internal/scene"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

var View = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	context := co.GetContext[global.Context]()

	co.Once(func() {
		co.Window(scope).SetCursorVisible(false)
	})

	co.Defer(func() {
		co.Window(scope).SetCursorVisible(true)
	})

	co.Once(func() {
		gameEngine := context.Engine
		gameData := scene.NewData(
			gameEngine,
			gameEngine.CreateResourceSet(),
		)
		dataRequest := gameData.Request()
		dataRequest.OnSuccess(func(struct{}) {
			mvc.Dispatch(scope, action.SetGameData{
				GameData: gameData,
			})
			mvc.Dispatch(scope, action.ChangeView{
				ViewName: model.ViewNameHome,
			})
		})
		dataRequest.OnError(func(err error) {
			panic(fmt.Errorf("failed to load assets: %w", err))
		})
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: opt.V(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("logo-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				BackgroundColor: opt.V(ui.Transparent()),
				Image:           co.OpenImage(scope, "ui/images/logo.png"),
				Mode:            mat.ImageModeFit,
			})
			co.WithLayoutData(mat.LayoutData{
				Width:            opt.V(512),
				Height:           opt.V(128),
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
		}))
	})
})
