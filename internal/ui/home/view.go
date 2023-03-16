package home

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

var View = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	onContinueClicked := func() {
		log.Info("Continue")
	}

	onNewGameClicked := func() {
		log.Info("New Game")
		mvc.Dispatch(scope, action.ChangeView{
			ViewName: model.ViewNamePlay,
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
			BackgroundColor: opt.V(ui.Black()),
			Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("background-picture", co.New(mat.Picture, func() {
			co.WithData(mat.PictureData{
				Image: co.OpenImage(scope, "ui/images/background.png"),
				Mode:  mat.ImageModeCover,
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(250),
				Right:  opt.V(0),
			})
		}))

		co.WithChild("button-holder", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Left:           opt.V(100),
				VerticalCenter: opt.V(0),
				Width:          opt.V(300),
				Height:         opt.V(200),
			})

			buttonWidth := opt.V(130)

			co.WithChild("continue-button", co.New(widget.HomeButton, func() {
				co.WithData(widget.HomeButtonData{
					Text: "Continue",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(0),
					Left:   opt.V(0),
					Width:  buttonWidth,
					Height: opt.V(30),
				})
				co.WithCallbackData(widget.HomeButtonCallbackData{
					ClickListener: onContinueClicked,
				})
			}))

			co.WithChild("new-game-button", co.New(widget.HomeButton, func() {
				co.WithData(widget.HomeButtonData{
					Text: "New Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(50),
					Left:   opt.V(0),
					Width:  buttonWidth,
					Height: opt.V(30),
				})
				co.WithCallbackData(widget.HomeButtonCallbackData{
					ClickListener: onNewGameClicked,
				})
			}))

			co.WithChild("load-game-button", co.New(widget.HomeButton, func() {
				co.WithData(widget.HomeButtonData{
					Text: "Load Game",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(100),
					Left:   opt.V(0),
					Width:  buttonWidth,
					Height: opt.V(30),
				})
				co.WithCallbackData(widget.HomeButtonCallbackData{
					ClickListener: onLoadGameClicked,
				})
			}))

			co.WithChild("options-button", co.New(widget.HomeButton, func() {
				co.WithData(widget.HomeButtonData{
					Text: "Options",
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(150),
					Left:   opt.V(0),
					Width:  buttonWidth,
					Height: opt.V(30),
				})
				co.WithCallbackData(widget.HomeButtonCallbackData{
					ClickListener: onOptionsClicked,
				})
			}))
		}))
	})
})
