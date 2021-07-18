package intro

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
	t "github.com/mokiat/lacking/ui/template"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
)

type ViewData struct {
	GFXEngine      graphics.Engine
	GameController *game.Controller
}

var View = t.Connect(t.ShallowCached(t.Plain(func(props t.Properties) t.Instance {
	var (
		data    ViewData
		logoImg ui.Image
	)
	props.InjectData(&data)

	t.UseState(func() interface{} {
		return t.OpenImage("resources/ui/intro/logo.png")
	}).Inject(&logoImg)

	t.Once(func() {
		gameData := scene.NewData(data.GameController.Registry(), data.GameController.GFXWorker())
		gameData.Request().OnSuccess(func(interface{}) {
			t.Window().Schedule(func() error {
				t.Dispatch(store.SetGameDataAction{
					GameData: gameData,
				})
				t.Dispatch(store.ChangeViewAction{
					ViewIndex: store.ViewPlay,
				})
				return nil
			})
		})
	})

	return t.New(mat.Container, func() {
		t.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		t.WithChild("logo-picture", t.New(mat.Picture, func() {
			t.WithData(mat.PictureData{
				BackgroundColor: ui.Transparent(),
				Image:           logoImg,
				Mode:            mat.ImageModeFit,
			})
			t.WithLayoutData(mat.AnchorLayoutData{
				Width:                    optional.NewInt(512),
				Height:                   optional.NewInt(128),
				HorizontalCenter:         optional.NewInt(0),
				HorizontalCenterRelation: mat.RelationCenter,
				VerticalCenter:           optional.NewInt(0),
				VerticalCenterRelation:   mat.RelationCenter,
			})
		}))
	})
})), func(props t.Properties, state *t.ReducedState) (interface{}, interface{}) {
	var appState store.Application
	state.Inject(&appState)

	return ViewData{
		GFXEngine:      appState.GFXEngine,
		GameController: appState.GameController,
	}, nil
})
