package view

import (
	"github.com/mokiat/gblob"
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
	"github.com/x448/float16"
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

	engine *game.Engine

	loadingModel *model.Loading
	homeModel    *model.Home
	playModel    *model.Play
	scene        *model.HomeScene

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

	p.scene = p.homeModel.Scene()
	if p.scene == nil {
		p.scene = p.createScene()
		p.homeModel.SetScene(p.scene)
		p.onDayClicked()
	}
	p.engine.SetActiveScene(p.scene.Scene)
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
								Selected: controller == data.ControllerKeyboard,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onKeyboardClicked,
							})
						}))

						co.WithChild("mouse", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Mouse",
								Selected: controller == data.ControllerMouse,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onMouseClicked,
							})
						}))

						co.WithChild("gamepad", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Gamepad",
								Selected: controller == data.ControllerGamepad,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onGamepadClicked,
							})
						}))
					}))

					var imageURL string
					switch controller {
					case data.ControllerKeyboard:
						imageURL = "ui/images/keyboard.png"
					case data.ControllerMouse:
						imageURL = "ui/images/mouse.png"
					case data.ControllerGamepad:
						imageURL = "ui/images/gamepad.png"
					}

					co.WithChild("controller-image", co.New(mat.Picture, func() {
						co.WithData(mat.PictureData{
							BackgroundColor: opt.V(ui.RGBA(0x00, 0x00, 0x00, 0x9A)),
							Image:           co.OpenImage(p.Scope, imageURL),
							ImageColor:      opt.V(ui.White()),
							Mode:            mat.ImageModeStretch,
						})
						co.WithLayoutData(mat.LayoutData{
							Width:  opt.V(600),
							Height: opt.V(300),
						})
					}))

					co.WithChild("controller-text", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:          co.OpenFont(p.Scope, "mat:///roboto-regular.ttf"),
							FontSize:      opt.V(float32(24.0)),
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
								Selected: environment == data.EnvironmentDay,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onDayClicked,
							})
						}))

						co.WithChild("night", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Night",
								Selected: environment == data.EnvironmentNight,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								ClickListener: p.onNightClicked,
							})
						}))
					}))

					co.WithChild("environment-text", co.New(mat.Label, func() {
						co.WithData(mat.LabelData{
							Font:          co.OpenFont(p.Scope, "mat:///roboto-regular.ttf"),
							FontSize:      opt.V(float32(24.0)),
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

func (p *HomeScreenPresenter) createScene() *model.HomeScene {
	result := &model.HomeScene{}

	promise := p.homeModel.Data()
	sceneData, err := promise.Get()
	if err != nil {
		log.Error("ERROR: %v", err) // TODO: Go to error screen
		return nil
	}

	scene := p.engine.CreateScene()
	scene.Initialize(sceneData.Scene)
	result.Scene = scene

	camera := p.createCamera(scene.Graphics())
	scene.Graphics().SetActiveCamera(camera)

	result.DaySkyColor = sprec.NewVec3(20.0, 25.0, 30.0)
	result.DayAmbientLight = p.createDayAmbientLight(scene.Graphics())
	result.DayDirectionalLight = scene.Graphics().CreateDirectionalLight(graphics.DirectionalLightInfo{
		EmitColor: dprec.NewVec3(10, 10, 6),
		EmitRange: 16000, // FIXME
	})
	dayDirectionalLightNode := game.NewNode()
	dayDirectionalLightNode.SetPosition(dprec.NewVec3(-100.0, 100.0, 0.0))
	dayDirectionalLightNode.SetRotation(dprec.QuatProd(
		dprec.RotationQuat(dprec.Degrees(-90), dprec.BasisYVec3()),
		dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
	))
	dayDirectionalLightNode.UseTransformation(game.IgnoreParentRotation)
	dayDirectionalLightNode.SetAttachable(result.DayDirectionalLight)

	result.NightSkyColor = sprec.NewVec3(0.01, 0.01, 0.01)
	result.NightAmbientLight = p.createNightAmbientLight(scene.Graphics())
	result.NightSpotLight = scene.Graphics().CreateSpotLight(graphics.SpotLightInfo{
		EmitColor:          dprec.NewVec3(5000.0, 5000.0, 7500.0),
		EmitOuterConeAngle: dprec.Degrees(50),
		EmitInnerConeAngle: dprec.Degrees(20),
		EmitRange:          1000,
	})
	nightSpotLightNode := game.NewNode()
	nightSpotLightNode.SetPosition(dprec.NewVec3(0.0, 0.0, 0.0))
	nightSpotLightNode.SetRotation(dprec.RotationQuat(dprec.Degrees(0), dprec.BasisXVec3()))
	nightSpotLightNode.SetAttachable(result.NightSpotLight)

	sceneModel := scene.FindModel("Content")
	// TODO: Remove manual attachment, once this is configurable from
	// the packing.
	scene.Root().AppendChild(sceneModel.Root())

	if cameraNode := scene.Root().FindNode("Camera"); cameraNode != nil {
		cameraNode.SetAttachable(camera)
		cameraNode.AppendChild(dayDirectionalLightNode)
		cameraNode.AppendChild(nightSpotLightNode)
	}

	const animationName = "Action"
	if animation := sceneModel.FindAnimation(animationName); animation != nil {
		playback := scene.PlayAnimation(animation)
		playback.SetLoop(true)
		playback.SetSpeed(0.3)
	}
	return result
}

func (p *HomeScreenPresenter) createCamera(scene *graphics.Scene) *graphics.Camera {
	result := scene.CreateCamera()
	result.SetFoVMode(graphics.FoVModeHorizontalPlus)
	result.SetFoV(sprec.Degrees(66))
	result.SetAutoExposure(false)
	result.SetExposure(0.1)
	result.SetAutoFocus(false)
	result.SetAutoExposureSpeed(0.1)
	return result
}

func (p *HomeScreenPresenter) createDayAmbientLight(scene *graphics.Scene) *graphics.AmbientLight {
	reflectionData := gblob.LittleEndianBlock(make([]byte, 4*2))
	reflectionData.SetUint16(0, float16.Fromfloat32(20.0).Bits())
	reflectionData.SetUint16(2, float16.Fromfloat32(25.0).Bits())
	reflectionData.SetUint16(4, float16.Fromfloat32(30.0).Bits())
	reflectionData.SetUint16(6, float16.Fromfloat32(1.0).Bits())

	reflectionTexture := p.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
		Dimension:       1,
		Filtering:       graphics.FilterNearest,
		InternalFormat:  graphics.InternalFormatRGBA16F,
		DataFormat:      graphics.DataFormatRGBA16F,
		GammaCorrection: false,
		GenerateMipmaps: false,
		FrontSideData:   reflectionData,
		BackSideData:    reflectionData,
		LeftSideData:    reflectionData,
		RightSideData:   reflectionData,
		TopSideData:     reflectionData,
		BottomSideData:  reflectionData,
	})

	refractionTexture := p.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
		Dimension:       1,
		Filtering:       graphics.FilterNearest,
		InternalFormat:  graphics.InternalFormatRGBA16F,
		DataFormat:      graphics.DataFormatRGBA16F,
		GammaCorrection: false,
		GenerateMipmaps: false,
		FrontSideData:   reflectionData,
		BackSideData:    reflectionData,
		LeftSideData:    reflectionData,
		RightSideData:   reflectionData,
		TopSideData:     reflectionData,
		BottomSideData:  reflectionData,
	})

	return scene.CreateAmbientLight(graphics.AmbientLightInfo{
		ReflectionTexture: reflectionTexture,
		RefractionTexture: refractionTexture,
		Position:          dprec.ZeroVec3(),
		InnerRadius:       5000,
		OuterRadius:       5000,
	})
}

func (p *HomeScreenPresenter) createNightAmbientLight(scene *graphics.Scene) *graphics.AmbientLight {
	reflectionData := gblob.LittleEndianBlock(make([]byte, 4*2))
	reflectionData.SetUint16(0, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(2, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(4, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(6, float16.Fromfloat32(1.0).Bits())

	reflectionTexture := p.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
		Dimension:       1,
		Filtering:       graphics.FilterNearest,
		InternalFormat:  graphics.InternalFormatRGBA16F,
		DataFormat:      graphics.DataFormatRGBA16F,
		GammaCorrection: false,
		GenerateMipmaps: false,
		FrontSideData:   reflectionData,
		BackSideData:    reflectionData,
		LeftSideData:    reflectionData,
		RightSideData:   reflectionData,
		TopSideData:     reflectionData,
		BottomSideData:  reflectionData,
	})

	refractionTexture := p.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
		Dimension:       1,
		Filtering:       graphics.FilterNearest,
		InternalFormat:  graphics.InternalFormatRGBA16F,
		DataFormat:      graphics.DataFormatRGBA16F,
		GammaCorrection: false,
		GenerateMipmaps: false,
		FrontSideData:   reflectionData,
		BackSideData:    reflectionData,
		LeftSideData:    reflectionData,
		RightSideData:   reflectionData,
		TopSideData:     reflectionData,
		BottomSideData:  reflectionData,
	})

	return scene.CreateAmbientLight(graphics.AmbientLightInfo{
		ReflectionTexture: reflectionTexture,
		RefractionTexture: refractionTexture,
		Position:          dprec.ZeroVec3(),
		InnerRadius:       5000,
		OuterRadius:       5000,
	})
}

func (p *HomeScreenPresenter) controllerDescription(controller data.Controller) string {
	switch controller {
	case data.ControllerKeyboard:
		return "Keyboard: Uses assists. Provides an average challenge."
	case data.ControllerMouse:
		return "Mouse: Just point and drive. Good for a casual play."
	case data.ControllerGamepad:
		return "Gamepad: No assists. Requires significant skills to control."
	default:
		return ""
	}
}

func (p *HomeScreenPresenter) environmentDescription(environment data.Environment) string {
	switch environment {
	case data.EnvironmentDay:
		return "Day: A good starting point to learn the track."
	case data.EnvironmentNight:
		return "Night: Can be relaxing if you already know the track."
	default:
		return ""
	}
}

func (p *HomeScreenPresenter) onKeyboardClicked() {
	p.homeModel.SetController(data.ControllerKeyboard)
}

func (p *HomeScreenPresenter) onMouseClicked() {
	p.homeModel.SetController(data.ControllerMouse)
}

func (p *HomeScreenPresenter) onGamepadClicked() {
	p.homeModel.SetController(data.ControllerGamepad)
}

func (p *HomeScreenPresenter) onDayClicked() {
	p.homeModel.SetEnvironment(data.EnvironmentDay)

	// Disable night lighting
	p.scene.NightAmbientLight.SetActive(false)
	p.scene.NightSpotLight.SetActive(false)

	// Enable day lighting
	p.scene.Scene.Graphics().Sky().SetBackgroundColor(p.scene.DaySkyColor)
	p.scene.DayAmbientLight.SetActive(true)
	p.scene.DayDirectionalLight.SetActive(true)
}

func (p *HomeScreenPresenter) onNightClicked() {
	p.homeModel.SetEnvironment(data.EnvironmentNight)

	// Disable day lighting
	p.scene.DayAmbientLight.SetActive(false)
	p.scene.DayDirectionalLight.SetActive(false)

	// Enable night lighting
	p.scene.Scene.Graphics().Sky().SetBackgroundColor(p.scene.NightSkyColor)
	p.scene.NightAmbientLight.SetActive(true)
	p.scene.NightSpotLight.SetActive(true)
}

func (p *HomeScreenPresenter) onStartClicked() {
	resourceSet := p.engine.CreateResourceSet()
	promise := data.LoadPlayData(p.engine, resourceSet, p.homeModel.Environment(), p.homeModel.Controller())
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
