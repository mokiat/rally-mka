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
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/action"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
	"github.com/x448/float16"
)

var HomeScreen = co.Define(&homeScreenComponent{})

type HomeScreenData struct {
	Loading *model.Loading
	Home    *model.Home
	Play    *model.Play
}

type homeScreenComponent struct {
	Scope      co.Scope       `co:"scope"`
	Data       HomeScreenData `co:"data"`
	Invalidate func()         `co:"invalidate"`

	engine      *game.Engine
	resourceSet *game.ResourceSet

	loadingModel *model.Loading
	homeModel    *model.Home
	playModel    *model.Play
	scene        *model.HomeScene

	showOptions bool
}

func (c *homeScreenComponent) OnCreate() {
	globalContext := co.TypedValue[global.Context](c.Scope)

	c.engine = globalContext.Engine
	c.resourceSet = globalContext.ResourceSet
	c.loadingModel = c.Data.Loading
	c.homeModel = c.Data.Home
	c.playModel = c.Data.Play

	// TODO: Figure out an alternative way for TypeComponents
	mvc.UseBinding(c.homeModel, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.HomeChange)
	})

	c.scene = c.homeModel.Scene()
	if c.scene == nil {
		c.scene = c.createScene()
		c.homeModel.SetScene(c.scene)
		c.onDayClicked()
	}
	c.engine.SetActiveScene(c.scene.Scene)
}

func (c *homeScreenComponent) OnDelete() {
	c.engine.SetActiveScene(nil)
}

func (c *homeScreenComponent) Render() co.Instance {
	controller := c.homeModel.Controller()
	environment := c.homeModel.Environment()

	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Layout: layout.Anchor(),
		})

		co.WithChild("pane", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				Top:    opt.V(0),
				Bottom: opt.V(0),
				Left:   opt.V(0),
				Width:  opt.V(320),
			})
			co.WithData(std.ContainerData{
				BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 192)),
				Layout:          layout.Anchor(),
			})

			co.WithChild("holder", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					Left:           opt.V(75),
					VerticalCenter: opt.V(0),
				})
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentAlignment: layout.HorizontalAlignmentLeft,
						ContentSpacing:   15,
					}),
				})

				if c.showOptions {
					co.WithChild("start-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Start",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onStartClicked,
						})
					}))

					co.WithChild("back-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Back",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onBackClicked,
						})
					}))
				} else {
					co.WithChild("play-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Play",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onPlayClicked,
						})
					}))

					co.WithChild("licenses-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Licenses",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onLicensesClicked,
						})
					}))

					co.WithChild("credits-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Credits",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onCreditsClicked,
						})
					}))

					co.WithChild("exit-button", co.New(widget.Button, func() {
						co.WithData(widget.ButtonData{
							Text: "Exit",
						})
						co.WithCallbackData(widget.ButtonCallbackData{
							OnClick: c.onExitClicked,
						})
					}))
				}
			}))
		}))

		if c.showOptions {
			co.WithChild("options", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					Top:    opt.V(0),
					Bottom: opt.V(0),
					Left:   opt.V(320),
					Right:  opt.V(0),
				})
				co.WithData(std.ContainerData{
					BackgroundColor: opt.V(ui.RGBA(0, 0, 0, 128)),
					Layout:          layout.Anchor(),
				})

				co.WithChild("options-pane", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						HorizontalCenter: opt.V(0),
						VerticalCenter:   opt.V(0),
					})
					co.WithData(std.ElementData{
						Layout: layout.Vertical(layout.VerticalSettings{
							ContentAlignment: layout.HorizontalAlignmentCenter,
							ContentSpacing:   20,
						}),
					})

					co.WithChild("controller-toggles", co.New(std.Element, func() {
						co.WithData(std.ElementData{
							Layout: layout.Horizontal(layout.HorizontalSettings{
								ContentAlignment: layout.VerticalAlignmentCenter,
								ContentSpacing:   40,
							}),
						})

						co.WithChild("keyboard", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Keyboard",
								Selected: controller == data.ControllerKeyboard,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								OnToggle: c.onKeyboardClicked,
							})
						}))

						co.WithChild("mouse", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Mouse",
								Selected: controller == data.ControllerMouse,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								OnToggle: c.onMouseClicked,
							})
						}))

						co.WithChild("gamepad", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Gamepad",
								Selected: controller == data.ControllerGamepad,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								OnToggle: c.onGamepadClicked,
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

					co.WithChild("controller-image", co.New(std.Picture, func() {
						co.WithLayoutData(layout.Data{
							Width:  opt.V(600),
							Height: opt.V(300),
						})
						co.WithData(std.PictureData{
							BackgroundColor: opt.V(ui.RGBA(0x00, 0x00, 0x00, 0x9A)),
							Image:           co.OpenImage(c.Scope, imageURL),
							ImageColor:      opt.V(ui.White()),
							Mode:            std.ImageModeStretch,
						})
					}))

					co.WithChild("controller-text", co.New(std.Label, func() {
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope, "ui:///roboto-regular.ttf"),
							FontSize:  opt.V(float32(24.0)),
							FontColor: opt.V(ui.White()),
							Text:      c.controllerDescription(controller),
						})
					}))

					co.WithChild("separator", co.New(widget.Separator, func() {
						co.WithLayoutData(layout.Data{
							Width:  opt.V(600),
							Height: opt.V(4),
						})
					}))

					co.WithChild("environment-toggles", co.New(std.Element, func() {
						co.WithData(std.ElementData{
							Layout: layout.Horizontal(layout.HorizontalSettings{
								ContentAlignment: layout.VerticalAlignmentCenter,
								ContentSpacing:   40,
							}),
						})

						co.WithChild("day", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Day",
								Selected: environment == data.EnvironmentDay,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								OnToggle: c.onDayClicked,
							})
						}))

						co.WithChild("night", co.New(widget.Toggle, func() {
							co.WithData(widget.ToggleData{
								Text:     "Night",
								Selected: environment == data.EnvironmentNight,
							})
							co.WithCallbackData(widget.ToggleCallbackData{
								OnToggle: c.onNightClicked,
							})
						}))
					}))

					co.WithChild("environment-text", co.New(std.Label, func() {
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope, "ui:///roboto-regular.ttf"),
							FontSize:  opt.V(float32(24.0)),
							FontColor: opt.V(ui.White()),
							Text:      c.environmentDescription(environment),
						})
					}))
				}))
			}))
		}
	})
}

func (c *homeScreenComponent) createScene() *model.HomeScene {
	result := &model.HomeScene{}

	promise := c.homeModel.Data()
	sceneData, err := promise.Get()
	if err != nil {
		log.Error("ERROR: %v", err) // TODO: Go to error screen
		return nil
	}

	scene := c.engine.CreateScene()
	scene.Initialize(sceneData.Scene)
	result.Scene = scene

	camera := c.createCamera(scene.Graphics())
	scene.Graphics().SetActiveCamera(camera)

	result.DaySkyColor = sprec.NewVec3(20.0, 25.0, 30.0)
	result.DayAmbientLight = c.createDayAmbientLight(scene.Graphics())
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
	result.NightAmbientLight = c.createNightAmbientLight(scene.Graphics())
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

func (c *homeScreenComponent) createCamera(scene *graphics.Scene) *graphics.Camera {
	result := scene.CreateCamera()
	result.SetFoVMode(graphics.FoVModeHorizontalPlus)
	result.SetFoV(sprec.Degrees(66))
	result.SetAutoExposure(false)
	result.SetExposure(0.1)
	result.SetAutoFocus(false)
	result.SetAutoExposureSpeed(0.1)
	return result
}

func (c *homeScreenComponent) createDayAmbientLight(scene *graphics.Scene) *graphics.AmbientLight {
	reflectionData := gblob.LittleEndianBlock(make([]byte, 4*2))
	reflectionData.SetUint16(0, float16.Fromfloat32(20.0).Bits())
	reflectionData.SetUint16(2, float16.Fromfloat32(25.0).Bits())
	reflectionData.SetUint16(4, float16.Fromfloat32(30.0).Bits())
	reflectionData.SetUint16(6, float16.Fromfloat32(1.0).Bits())

	reflectionTexture := c.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
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

	refractionTexture := c.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
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

func (c *homeScreenComponent) createNightAmbientLight(scene *graphics.Scene) *graphics.AmbientLight {
	reflectionData := gblob.LittleEndianBlock(make([]byte, 4*2))
	reflectionData.SetUint16(0, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(2, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(4, float16.Fromfloat32(0.1).Bits())
	reflectionData.SetUint16(6, float16.Fromfloat32(1.0).Bits())

	reflectionTexture := c.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
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

	refractionTexture := c.engine.Graphics().CreateCubeTexture(graphics.CubeTextureDefinition{
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

func (c *homeScreenComponent) controllerDescription(controller data.Controller) string {
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

func (c *homeScreenComponent) environmentDescription(environment data.Environment) string {
	switch environment {
	case data.EnvironmentDay:
		return "Day: A good starting point to learn the track."
	case data.EnvironmentNight:
		return "Night: Can be relaxing if you already know the track."
	default:
		return ""
	}
}

func (c *homeScreenComponent) onKeyboardClicked() {
	c.homeModel.SetController(data.ControllerKeyboard)
}

func (c *homeScreenComponent) onMouseClicked() {
	c.homeModel.SetController(data.ControllerMouse)
}

func (c *homeScreenComponent) onGamepadClicked() {
	c.homeModel.SetController(data.ControllerGamepad)
}

func (c *homeScreenComponent) onDayClicked() {
	c.homeModel.SetEnvironment(data.EnvironmentDay)

	// Disable night lighting
	c.scene.NightAmbientLight.SetActive(false)
	c.scene.NightSpotLight.SetActive(false)

	// Enable day lighting
	c.scene.Scene.Graphics().Sky().SetBackgroundColor(c.scene.DaySkyColor)
	c.scene.DayAmbientLight.SetActive(true)
	c.scene.DayDirectionalLight.SetActive(true)
}

func (c *homeScreenComponent) onNightClicked() {
	c.homeModel.SetEnvironment(data.EnvironmentNight)

	// Disable day lighting
	c.scene.DayAmbientLight.SetActive(false)
	c.scene.DayDirectionalLight.SetActive(false)

	// Enable night lighting
	c.scene.Scene.Graphics().Sky().SetBackgroundColor(c.scene.NightSkyColor)
	c.scene.NightAmbientLight.SetActive(true)
	c.scene.NightSpotLight.SetActive(true)
}

func (c *homeScreenComponent) onStartClicked() {
	promise := data.LoadPlayData(c.engine, c.resourceSet, c.homeModel.Environment(), c.homeModel.Controller())
	c.playModel.SetData(promise)

	c.loadingModel.SetPromise(promise)
	c.loadingModel.SetNextViewName(model.ViewNamePlay)
	mvc.Dispatch(c.Scope, action.ChangeView{
		ViewName: model.ViewNameLoading,
	})
}

func (c *homeScreenComponent) onBackClicked() {
	// TODO: Add a `Property` concept instead of manual Invalidate.
	c.showOptions = false
	c.Invalidate()
}

func (c *homeScreenComponent) onPlayClicked() {
	c.showOptions = true
	c.Invalidate()
}

func (c *homeScreenComponent) onLicensesClicked() {
	mvc.Dispatch(c.Scope, action.ChangeView{
		ViewName: model.ViewNameLicenses,
	})
}

func (c *homeScreenComponent) onCreditsClicked() {
	mvc.Dispatch(c.Scope, action.ChangeView{
		ViewName: model.ViewNameCredits,
	})
}

func (c *homeScreenComponent) onExitClicked() {
	co.Window(c.Scope).Close()
}
