package view

import (
	"sync"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
)

type HomeScreenData struct {
	Loading *model.Loading
	Home    *model.Home
	Play    *model.Play
}

var HomeScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		globalContext = co.GetContext[global.Context]()
		screenData    = co.GetData[HomeScreenData](props)

		engine       = globalContext.Engine
		loadingModel = screenData.Loading
		homeModel    = screenData.Home
		playModel    = screenData.Play
	)

	co.Once(func() {
		if scene := homeModel.Scene(); scene != nil {
			engine.SetActiveScene(homeModel.Scene())
			return
		}

		promise := homeModel.Data()
		sceneData, err := promise.Get()
		if err != nil {
			log.Error("ERROR: %v", err) // TODO: Go to error screen
			return
		}

		scene := engine.CreateScene()
		scene.Initialize(sceneData.Scene)

		camera := scene.Graphics().CreateCamera()
		camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
		camera.SetFoV(sprec.Degrees(66))
		camera.SetAutoExposure(true)
		camera.SetExposure(0.01)
		camera.SetAutoFocus(false)
		camera.SetAutoExposureSpeed(0.1)
		scene.Graphics().SetActiveCamera(camera)

		sunLight := scene.Graphics().CreateDirectionalLight(graphics.DirectionalLightInfo{
			EmitColor: dprec.NewVec3(0.5, 0.5, 0.3),
			EmitRange: 16000, // FIXME
		})

		lightNode := game.NewNode()
		lightNode.SetPosition(dprec.NewVec3(-100.0, 100.0, 0.0))
		lightNode.SetRotation(dprec.QuatProd(
			dprec.RotationQuat(dprec.Degrees(-90), dprec.BasisYVec3()),
			dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
		))
		// FIXME: This should work out of the box for directional lights
		lightNode.UseTransformation(func(parent, current dprec.Mat4) dprec.Mat4 {
			// Remove parent's rotation
			parent.M11 = 1.0
			parent.M12 = 0.0
			parent.M13 = 0.0
			parent.M21 = 0.0
			parent.M22 = 1.0
			parent.M23 = 0.0
			parent.M31 = 0.0
			parent.M32 = 0.0
			parent.M33 = 1.0
			return dprec.Mat4Prod(parent, current)
		})
		lightNode.SetDirectionalLight(sunLight)

		sceneModel := scene.FindModel("Content")
		// TODO: Remove manual attachment, once this is configurable from
		// the packing.
		scene.Root().AppendChild(sceneModel.Root())

		if cameraNode := scene.Root().FindNode("Camera"); cameraNode != nil {
			cameraNode.SetCamera(camera)
			cameraNode.AppendChild(lightNode)
		}

		if animation := sceneModel.FindAnimation("Action"); animation != nil {
			// playback := scene.PlayAnimation(animation) // TODO
			startTime := animation.StartTime()
			endTime := animation.EndTime()
			go func() {
				var animMU sync.Mutex
				animTime := startTime
				for range time.Tick(10 * time.Millisecond) {
					animMU.Lock()
					animTime += (10 * time.Millisecond).Seconds() * 0.3
					if animTime >= endTime {
						animTime -= (endTime - startTime)
					}
					animMU.Unlock()
					co.Schedule(func() {
						animMU.Lock()
						animation.Apply(animTime)
						animMU.Unlock()
					})
				}
			}()
		}

		homeModel.SetScene(scene)
	})

	co.Defer(func() {
		engine.SetActiveScene(nil)
	})

	onPlayClicked := func() {
		// TODO: Show an intermediate configuration window, where the user
		// can select controller type and assitance.

		co.Once(func() {
			resourceSet := engine.CreateResourceSet()
			promise := data.LoadPlayData(engine, resourceSet)
			playModel.SetData(promise)

			loadingModel.SetPromise(promise)
			loadingModel.SetNextViewName(model.ViewNamePlay)
			mvc.Dispatch(scope, action.ChangeView{
				ViewName: model.ViewNameLoading,
			})
		})
	}

	onLicensesClicked := func() {
		mvc.Dispatch(scope, action.ChangeView{
			ViewName: model.ViewNameLicenses,
		})
	}

	onCreditsClicked := func() {
		mvc.Dispatch(scope, action.ChangeView{
			ViewName: model.ViewNameCredits,
		})
	}

	onExitClicked := func() {
		co.Window(scope).Close()
	}

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})

		co.WithChild("pane", co.New(mat.Container, func() {
			co.WithData(mat.ContainerData{
				BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 192)),
				Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
			})
			co.WithLayoutData(mat.LayoutData{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(0),
				Width:  opt.V(320),
			})

			co.WithChild("holder", co.New(mat.Element, func() {
				co.WithData(mat.ElementData{
					Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
						ContentAlignment: mat.AlignmentLeft,
						ContentSpacing:   15,
					}),
				})
				co.WithLayoutData(mat.LayoutData{
					Left:           opt.V(75),
					VerticalCenter: opt.V(0),
				})

				co.WithChild("play-button", co.New(widget.HomeButton, func() {
					co.WithData(widget.HomeButtonData{
						Text: "Play",
					})
					co.WithCallbackData(widget.HomeButtonCallbackData{
						ClickListener: onPlayClicked,
					})
				}))

				co.WithChild("licenses-button", co.New(widget.HomeButton, func() {
					co.WithData(widget.HomeButtonData{
						Text: "Licenses",
					})
					co.WithCallbackData(widget.HomeButtonCallbackData{
						ClickListener: onLicensesClicked,
					})
				}))

				co.WithChild("credits-button", co.New(widget.HomeButton, func() {
					co.WithData(widget.HomeButtonData{
						Text: "Credits",
					})
					co.WithCallbackData(widget.HomeButtonCallbackData{
						ClickListener: onCreditsClicked,
					})
				}))

				co.WithChild("exit-button", co.New(widget.HomeButton, func() {
					co.WithData(widget.HomeButtonData{
						Text: "Exit",
					})
					co.WithCallbackData(widget.HomeButtonCallbackData{
						ClickListener: onExitClicked,
					})
				}))
			}))
		}))
	})
})
