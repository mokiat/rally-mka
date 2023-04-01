package view

import (
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

var HomeScreen = co.DefineType(&HomeScreenPresenter{})

type HomeScreenData struct {
	Loading *model.Loading
	Home    *model.Home
	Play    *model.Play
}

type HomeScreenPresenter struct {
	Scope      co.Scope       `co:"scope"`
	Data       HomeScreenData `co:"data"`
	Invalidate func()         `co:"invalidate"`

	engine       *game.Engine
	loadingModel *model.Loading
	homeModel    *model.Home
	playModel    *model.Play

	showOptions bool
}

func (p *HomeScreenPresenter) OnCreate() {
	var globalContext global.Context
	co.InjectContext(&globalContext)

	p.engine = globalContext.Engine
	p.loadingModel = p.Data.Loading
	p.homeModel = p.Data.Home
	p.playModel = p.Data.Play

	// TODO: Figure out an alternative way for TypeComponents
	mvc.UseBinding(p.homeModel, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.HomeChange)
	})

	if p.homeModel.Scene() == nil {
		p.homeModel.SetScene(p.createScene())
	}
	scene := p.homeModel.Scene()
	p.engine.SetActiveScene(scene)
}

func (p *HomeScreenPresenter) OnDelete() {
	p.engine.SetActiveScene(nil)
}

func (p *HomeScreenPresenter) Render() co.Instance {
	controller := p.homeModel.Controller()
	environment := p.homeModel.Environment()

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

				if p.showOptions {
					co.WithChild("start-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Start",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onStartClicked,
						})
					}))

					co.WithChild("back-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Back",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onBackClicked,
						})
					}))
				} else {
					co.WithChild("play-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Play",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onPlayClicked,
						})
					}))

					co.WithChild("licenses-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Licenses",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onLicensesClicked,
						})
					}))

					co.WithChild("credits-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Credits",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onCreditsClicked,
						})
					}))

					co.WithChild("exit-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Exit",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							ClickListener: p.onExitClicked,
						})
					}))
				}
			}))
		}))

		if p.showOptions {
			co.WithChild("options", co.New(mat.Container, func() {
				co.WithData(mat.ContainerData{
					BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 128)),
					Layout:          mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
				})
				co.WithLayoutData(mat.LayoutData{
					Top:    opt.V(0),
					Bottom: opt.V(0),
					Left:   opt.V(320),
					Right:  opt.V(0),
				})

				co.WithChild("options-pane", co.New(mat.Element, func() {
					co.WithData(mat.ElementData{
						Layout: mat.NewVerticalLayout(mat.VerticalLayoutSettings{
							ContentAlignment: mat.AlignmentCenter,
							ContentSpacing:   20,
						}),
					})
					co.WithLayoutData(mat.LayoutData{
						HorizontalCenter: opt.V(0),
						VerticalCenter:   opt.V(0),
					})

					co.WithChild("controller-toggles", co.New(mat.Element, func() {
						co.WithData(mat.ElementData{
							Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
								ContentAlignment: mat.AlignmentCenter,
								ContentSpacing:   40,
							}),
						})

						co.WithChild("keyboard", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Keyboard",
								Selected: controller == model.ControllerKeyboard,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onKeyboardClicked,
							})
						}))

						co.WithChild("mouse", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Mouse",
								Selected: controller == model.ControllerMouse,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onMouseClicked,
							})
						}))

						co.WithChild("gamepad", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Gamepad",
								Selected: controller == model.ControllerGamepad,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onGamepadClicked,
							})
						}))
					}))

					co.WithChild("controller-image", co.New(mat.Container, func() {
						co.WithData(mat.ContainerData{
							BackgroundColor: opt.V(ui.Black()),
						})
						co.WithLayoutData(mat.LayoutData{
							Width:  opt.V(600),
							Height: opt.V(300),
						})
					}))

					co.WithChild("controller-text", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:          co.OpenFont(p.Scope, "mat:///roboto-regular.ttf"),
							FontSize:      opt.V(float32(28.0)),
							FontColor:     opt.V(ui.White()),
							TextAlignment: mat.AlignmentCenter,
							Text:          p.controllerDescription(controller),
						})
					}))

					co.WithChild("separator", co.New(widget.Separator, func() {
						co.WithLayoutData(mat.LayoutData{
							Width:  opt.V(600),
							Height: opt.V(4),
						})
					}))

					co.WithChild("environment-toggles", co.New(mat.Element, func() {
						co.WithData(mat.ElementData{
							Layout: mat.NewHorizontalLayout(mat.HorizontalLayoutSettings{
								ContentAlignment: mat.AlignmentCenter,
								ContentSpacing:   40,
							}),
						})

						co.WithChild("day", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Day",
								Selected: environment == model.EnvironmentDay,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onDayClicked,
							})
						}))

						co.WithChild("night", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Night",
								Selected: environment == model.EnvironmentNight,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onNightClicked,
							})
						}))
					}))

					co.WithChild("environment-text", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:          co.OpenFont(p.Scope, "mat:///roboto-regular.ttf"),
							FontSize:      opt.V(float32(28.0)),
							FontColor:     opt.V(ui.White()),
							TextAlignment: mat.AlignmentCenter,
							Text:          p.environmentDescription(environment),
						})
					}))
				}))
			}))
		}
	})
}

func (p *HomeScreenPresenter) createScene() *game.Scene {
	promise := p.homeModel.Data()
	sceneData, err := promise.Get()
	if err != nil {
		log.Error("ERROR: %v", err) // TODO: Go to error screen
		return nil
	}

	scene := p.engine.CreateScene()
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

	const animationName = "Action"
	// const animationName = "GroundAction"
	if animation := sceneModel.FindAnimation(animationName); animation != nil {
		playback := scene.PlayAnimation(animation)
		playback.SetLoop(true)
		playback.SetSpeed(0.3)
	}

	return scene
}

func (p *HomeScreenPresenter) controllerDescription(controller model.Controller) string {
	switch controller {
	case model.ControllerKeyboard:
		return "The day drive is a good starting point for new players."
	case model.ControllerMouse:
		return "The night drive is a bit more taxing on the computer."
	case model.ControllerGamepad:
		return "The night drive is a bit more taxing on the computer."
	default:
		return ""
	}
}

func (p *HomeScreenPresenter) environmentDescription(environment model.Environment) string {
	switch environment {
	case model.EnvironmentDay:
		return "The day drive is a good starting point for new players."
	case model.EnvironmentNight:
		return "The night drive is a bit more taxing on the computer."
	default:
		return ""
	}
}

func (p *HomeScreenPresenter) onKeyboardClicked() {
	p.homeModel.SetController(model.ControllerKeyboard)
}

func (p *HomeScreenPresenter) onMouseClicked() {
	p.homeModel.SetController(model.ControllerMouse)
}

func (p *HomeScreenPresenter) onGamepadClicked() {
	p.homeModel.SetController(model.ControllerGamepad)
}

func (p *HomeScreenPresenter) onDayClicked() {
	p.homeModel.SetEnvironment(model.EnvironmentDay)
}

func (p *HomeScreenPresenter) onNightClicked() {
	p.homeModel.SetEnvironment(model.EnvironmentNight)
}

func (p *HomeScreenPresenter) onStartClicked() {
	resourceSet := p.engine.CreateResourceSet()
	promise := data.LoadPlayData(p.engine, resourceSet)
	p.playModel.SetData(promise)

	p.loadingModel.SetPromise(promise)
	p.loadingModel.SetNextViewName(model.ViewNamePlay)
	mvc.Dispatch(p.Scope, action.ChangeView{
		ViewName: model.ViewNameLoading,
	})
}

func (p *HomeScreenPresenter) onBackClicked() {
	// TODO: Add a `Property` concept instead of manual Invalidate.
	p.showOptions = false
	p.Invalidate()
}

func (p *HomeScreenPresenter) onPlayClicked() {
	p.showOptions = true
	p.Invalidate()
}

func (p *HomeScreenPresenter) onLicensesClicked() {
	mvc.Dispatch(p.Scope, action.ChangeView{
		ViewName: model.ViewNameLicenses,
	})
}

func (p *HomeScreenPresenter) onCreditsClicked() {
	mvc.Dispatch(p.Scope, action.ChangeView{
		ViewName: model.ViewNameCredits,
	})
}

func (p *HomeScreenPresenter) onExitClicked() {
	co.Window(p.Scope).Close()
}
