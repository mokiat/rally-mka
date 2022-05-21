package home

import (
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
	"github.com/mokiat/rally-mka/internal/store"
)

var View = co.Define(func(props co.Properties) co.Instance {
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

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			BackgroundColor: optional.Value(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("background-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				Image: co.OpenImage("ui/images/background.png"),
				Mode:  mat.ImageModeCover,
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    optional.Value(0),
				Bottom: optional.Value(0),
				Left:   optional.Value(250),
				Right:  optional.Value(0),
			})
		}))

		co.WithChild("button-holder", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Left:           optional.Value(100),
				VerticalCenter: optional.Value(0),
				Width:          optional.Value(300),
				Height:         optional.Value(200),
			})

			buttonPadding := ui.Spacing{
				Left:   5,
				Right:  5,
				Top:    2,
				Bottom: 2,
			}
			buttonWidth := optional.Value(130)

			co.WithChild("continue-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.OpenFont("mat:///roboto-bold.ttf"),
					FontSize:      optional.Value(float32(26)),
					FontColor:     optional.Value(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Continue",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.Value(0),
					Left:   optional.Value(0),
					Width:  buttonWidth,
					Height: optional.Value(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onContinueClicked,
				})
			}))

			co.WithChild("new-game-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.OpenFont("mat:///roboto-bold.ttf"),
					FontSize:      optional.Value(float32(26)),
					FontColor:     optional.Value(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "New Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.Value(50),
					Left:   optional.Value(0),
					Width:  buttonWidth,
					Height: optional.Value(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onNewGameClicked,
				})
			}))

			co.WithChild("load-game-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.OpenFont("mat:///roboto-bold.ttf"),
					FontSize:      optional.Value(float32(26)),
					FontColor:     optional.Value(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Load Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.Value(100),
					Left:   optional.Value(0),
					Width:  buttonWidth,
					Height: optional.Value(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onLoadGameClicked,
				})
			}))

			co.WithChild("options-button", co.New(mat.Button, func() {
				co.WithData(mat.ButtonData{
					Padding:       buttonPadding,
					Font:          co.OpenFont("mat:///roboto-bold.ttf"),
					FontSize:      optional.Value(float32(26)),
					FontColor:     optional.Value(ui.White()),
					FontAlignment: mat.AlignmentLeft,
					Text:          "Options",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    optional.Value(150),
					Left:   optional.Value(0),
					Width:  buttonWidth,
					Height: optional.Value(30),
				})
				co.WithCallbackData(mat.ButtonCallbackData{
					ClickListener: onOptionsClicked,
				})
			}))
		}))
	})
})
