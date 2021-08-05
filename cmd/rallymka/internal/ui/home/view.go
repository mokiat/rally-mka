package home

import (
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/optional"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/global"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
)

var View = co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	co.OpenFontCollection("resources/ui/fonts/roboto.ttc")

	onContinueClicked := func() {
		log.Info("Continue")
	}

	onNewGameClicked := func() {
		log.Info("New Game")
		co.Dispatch(store.ChangeViewAction{
			ViewIndex: store.ViewPlay,
		})
	}

	onLoadGameClicked := func() {
		log.Info("Load Game")
	}

	onOptionsClicked := func() {
		log.Info("Options")
	}

	var context global.Context
	co.InjectContext(&context)

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.NewColor(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("background-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				Image: co.OpenImage("resources/ui/images/background.png"),
				Mode:  mat.ImageModeCover,
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    optional.NewInt(0),
				Bottom: optional.NewInt(0),
				Left:   optional.NewInt(300),
				Right:  optional.NewInt(0),
			})
		}))

		co.WithChild("button-holder", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Left:           optional.NewInt(100),
				VerticalCenter: optional.NewInt(0),
				Width:          optional.NewInt(300),
				Height:         optional.NewInt(200),
			})

			buttonPadding := ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    2,
				Bottom: 2,
			}
			buttonWidth := optional.NewInt(130)

			co.WithChild("continue-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.GetFont("roboto", "bold"),
					FontSize:      optional.NewInt(26),
					FontColor:     optional.NewColor(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Continue",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.NewInt(0),
					Left:   optional.NewInt(0),
					Width:  buttonWidth,
					Height: optional.NewInt(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onContinueClicked,
				})
			}))

			co.WithChild("new-game-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.GetFont("roboto", "bold"),
					FontSize:      optional.NewInt(26),
					FontColor:     optional.NewColor(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "New Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.NewInt(50),
					Left:   optional.NewInt(0),
					Width:  buttonWidth,
					Height: optional.NewInt(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onNewGameClicked,
				})
			}))

			co.WithChild("load-game-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.GetFont("roboto", "bold"),
					FontSize:      optional.NewInt(26),
					FontColor:     optional.NewColor(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Load Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.NewInt(100),
					Left:   optional.NewInt(0),
					Width:  buttonWidth,
					Height: optional.NewInt(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onLoadGameClicked,
				})
			}))

			co.WithChild("options-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.GetFont("roboto", "bold"),
					FontSize:      optional.NewInt(26),
					FontColor:     optional.NewColor(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Options",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.NewInt(150),
					Left:   optional.NewInt(0),
					Width:  buttonWidth,
					Height: optional.NewInt(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onOptionsClicked,
				})
			}))
		}))
	})
}))
