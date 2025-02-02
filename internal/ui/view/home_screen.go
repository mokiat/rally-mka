package view

import (
	"fmt"
	"time"

	"github.com/mokiat/gblob"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/game/data"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/widget"
	"github.com/x448/float16"
)

const (
	buttonAppearAfter     = time.Second + 100*time.Millisecond
	buttonAppearIncrement = 100 * time.Millisecond
)

var HomeScreen = co.Define(&homeScreenComponent{})

type HomeScreenData struct {
	AppModel     *model.ApplicationModel
	ErrorModel   *model.ErrorModel
	LoadingModel *model.LoadingModel
	HomeModel    *model.HomeModel
	PlayModel    *model.PlayModel
}

type homeScreenComponent struct {
	co.BaseComponent

	engine      *game.Engine
	resourceSet *game.ResourceSet

	appModel     *model.ApplicationModel
	errorModel   *model.ErrorModel
	loadingModel *model.LoadingModel
	homeModel    *model.HomeModel
	playModel    *model.PlayModel
	scene        *model.HomeScene
}

func (c *homeScreenComponent) OnCreate() {
	globalContext := co.TypedValue[global.Context](c.Scope())
	c.engine = globalContext.Engine
	c.resourceSet = globalContext.ResourceSet

	data := co.GetData[HomeScreenData](c.Properties())
	c.appModel = data.AppModel
	c.errorModel = data.ErrorModel
	c.loadingModel = data.LoadingModel
	c.homeModel = data.HomeModel
	c.playModel = data.PlayModel

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
	mode := c.homeModel.Mode()
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

				switch mode {
				case model.HomeScreenModeEntry:
					c.withEntryModeMenu()
				case model.HomeScreenModeLighting:
					c.withLightingModeMenu()
				case model.HomeScreenModeControls:
					c.withControlsModeMenu()
				case model.HomeScreenModeLevel:
					c.withLevelModeMenu()
				}
			}))
		}))

		switch mode {
		case model.HomeScreenModeEntry:
			// Nothing to show.
		case model.HomeScreenModeLighting:
			// Nothing to show.
		case model.HomeScreenModeControls:
			c.withControlsModeContent()
		case model.HomeScreenModeLevel:
			c.withLevelModeContent()
		}
	})
}

func (c *homeScreenComponent) withEntryModeMenu() {
	co.WithChild("entry-play-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Play",
			AppearAfter: buttonAppearAfter,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onPlayClicked,
		})
	}))

	co.WithChild("entry-licenses-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Licenses",
			AppearAfter: buttonAppearAfter + buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onLicensesClicked,
		})
	}))

	co.WithChild("entry-credits-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Credits",
			AppearAfter: buttonAppearAfter + 2*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onCreditsClicked,
		})
	}))

	co.WithChild("entry-exit-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Exit",
			AppearAfter: buttonAppearAfter + 3*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onExitClicked,
		})
	}))
}

func (c *homeScreenComponent) withLightingModeMenu() {
	environment := c.homeModel.Lighting()

	co.WithChild("lighting-day-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Day",
			Selected:    environment == data.LightingDay,
			AppearAfter: buttonAppearAfter,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onDayClicked,
		})
	}))

	co.WithChild("lighting-night-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Night",
			Selected:    environment == data.LightingNight,
			AppearAfter: buttonAppearAfter + buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onNightClicked,
		})
	}))

	co.WithChild("lighting-padding", co.New(std.Spacing, func() {
		co.WithData(std.SpacingData{
			Size: ui.NewSize(10, 32),
		})
	}))

	co.WithChild("lighting-next-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Confirm",
			AppearAfter: buttonAppearAfter + 2*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onNextClicked,
		})
	}))

	co.WithChild("lighting-back-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Back",
			AppearAfter: buttonAppearAfter + 3*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onBackClicked,
		})
	}))
}

func (c *homeScreenComponent) withControlsModeMenu() {
	controller := c.homeModel.Input()

	co.WithChild("controls-keyboard-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Keyboard",
			Selected:    controller == data.InputKeyboard,
			AppearAfter: buttonAppearAfter,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onKeyboardClicked,
		})
	}))

	co.WithChild("controls-mouse-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Mouse",
			Selected:    controller == data.InputMouse,
			AppearAfter: buttonAppearAfter + buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onMouseClicked,
		})
	}))

	co.WithChild("controls-gamepad-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Gamepad",
			Selected:    controller == data.InputGamepad,
			AppearAfter: buttonAppearAfter + 2*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onGamepadClicked,
		})
	}))

	co.WithChild("controls-padding", co.New(std.Spacing, func() {
		co.WithData(std.SpacingData{
			Size: ui.NewSize(10, 32),
		})
	}))

	co.WithChild("controls-next-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Confirm",
			AppearAfter: buttonAppearAfter + 3*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onNextClicked,
		})
	}))

	co.WithChild("controls-back-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Back",
			AppearAfter: buttonAppearAfter + 4*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onBackClicked,
		})
	}))
}

func (c *homeScreenComponent) withLevelModeMenu() {
	appearAfter := buttonAppearAfter
	selectedLevel := c.homeModel.Level()

	for i, level := range data.Levels {
		co.WithChild(fmt.Sprintf("level-%d", i), co.New(widget.Button, func() {
			co.WithData(widget.ButtonData{
				Text:        level.Name,
				Selected:    level == selectedLevel,
				AppearAfter: appearAfter,
			})
			co.WithCallbackData(widget.ButtonCallbackData{
				OnClick: func() {
					c.onLevelClicked(level)
				},
			})
		}))
		appearAfter += buttonAppearIncrement
	}

	co.WithChild("level-padding", co.New(std.Spacing, func() {
		co.WithData(std.SpacingData{
			Size: ui.NewSize(10, 32),
		})
	}))

	co.WithChild("level-start-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Start",
			AppearAfter: appearAfter + buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onStartClicked,
		})
	}))

	co.WithChild("level-back-button", co.New(widget.Button, func() {
		co.WithData(widget.ButtonData{
			Text:        "Back",
			AppearAfter: appearAfter + 2*buttonAppearIncrement,
		})
		co.WithCallbackData(widget.ButtonCallbackData{
			OnClick: c.onBackClicked,
		})
	}))
}

func (c *homeScreenComponent) withControlsModeContent() {
	controller := c.homeModel.Input()

	var controllerImage string
	switch controller {
	case data.InputKeyboard:
		controllerImage = "ui/images/keyboard.png"
	case data.InputMouse:
		controllerImage = "ui/images/mouse.png"
	case data.InputGamepad:
		controllerImage = "ui/images/gamepad.png"
	}
	image := co.OpenImage(c.Scope(), controllerImage)

	co.WithChild("panel", co.New(std.Container, func() {
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

		co.WithChild("centered-pane", co.New(std.Element, func() {
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

			co.WithChild("controller-image", co.New(std.Picture, func() {
				co.WithLayoutData(layout.Data{
					Width:  opt.V(image.Size().Width / 2),
					Height: opt.V(image.Size().Height / 2),
				})
				co.WithData(std.PictureData{
					Image:      image,
					ImageColor: opt.V(ui.White()),
					Mode:       std.ImageModeStretch,
				})
			}))

			co.WithChild("controller-text", co.New(std.Label, func() {
				co.WithData(std.LabelData{
					Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
					FontSize:  opt.V(float32(24.0)),
					FontColor: opt.V(ui.White()),
					Text:      c.controllerDescription(controller),
				})
			}))
		}))
	}))
}

func (c *homeScreenComponent) withLevelModeContent() {
	level := c.homeModel.Level()
	co.WithChild("panel", co.New(std.Element, func() {
		co.WithLayoutData(layout.Data{
			Top:    opt.V(0),
			Bottom: opt.V(0),
			Left:   opt.V(320),
			Right:  opt.V(0),
		})
		co.WithData(std.ElementData{
			Layout: layout.Anchor(),
		})

		co.WithChild("level-"+level.Name, co.New(widget.Level, func() {
			co.WithLayoutData(layout.Data{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
			co.WithData(widget.LevelData{
				Board: level.Board,
			})
		}))
	}))
}

func (c *homeScreenComponent) createScene() *model.HomeScene {
	result := &model.HomeScene{}

	sceneData := c.homeModel.Data()

	scene := c.engine.CreateScene()
	backgroundModel := scene.CreateModel(game.ModelInfo{
		Name:       "Background",
		Definition: sceneData.Background,
		IsDynamic:  true,
	})

	daySkyNode := backgroundModel.FindNode("Sky-Day")
	result.DaySky = (daySkyNode.Target().(game.SkyNodeTarget)).Sky

	nightSkyNode := backgroundModel.FindNode("Sky-Night")
	result.NightSky = (nightSkyNode.Target().(game.SkyNodeTarget)).Sky

	sceneModel := scene.CreateModel(game.ModelInfo{
		Name:       "HomeScreen",
		Definition: sceneData.Scene,
		IsDynamic:  false,
	})
	scene.Root().AppendChild(sceneModel.Root())
	result.Scene = scene

	camera := c.createCamera(scene.Graphics())
	scene.Graphics().SetActiveCamera(camera)

	scene.CreateModel(game.ModelInfo{
		Name:       "Vehicle",
		Definition: sceneData.Vehicle,
		Position:   opt.V(dprec.NewVec3(0.0, 0.0, 0.4)),
		IsDynamic:  false,
	})

	// TODO: Move these to the scene data
	result.DayAmbientLight = c.createDayAmbientLight(scene.Graphics())
	result.DayDirectionalLight = scene.Graphics().CreateDirectionalLight(graphics.DirectionalLightInfo{
		Position: dprec.NewVec3(-100.0, 100.0, 0.0),
		Rotation: dprec.QuatProd(
			dprec.RotationQuat(dprec.Degrees(-90), dprec.BasisYVec3()),
			dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
		),
		EmitColor:  dprec.NewVec3(1, 1, 0.6),
		CastShadow: true,
	})

	result.NightAmbientLight = c.createNightAmbientLight(scene.Graphics())

	if cameraNode := scene.Root().FindNode("Camera"); cameraNode != nil {
		cameraNode.SetTarget(game.CameraNodeTarget{
			Camera: camera,
		})
	}

	const animationName = "Action"
	if animation := sceneModel.FindAnimation(animationName); animation != nil {
		playback := animation.Playback()
		playback.SetLoop(true)
		sceneModel.BindAnimationSource(playback)
		scene.PlayAnimationTree(playback)
	}
	return result
}

func (c *homeScreenComponent) createCamera(scene *graphics.Scene) *graphics.Camera {
	result := scene.CreateCamera()
	result.SetFoVMode(graphics.FoVModeHorizontalPlus)
	result.SetFoV(sprec.Degrees(30))
	result.SetAutoExposure(false)
	result.SetExposure(1.0)
	result.SetAutoFocus(false)
	result.SetAutoExposureSpeed(0.1)
	result.SetCascadeDistances([]float32{16.0})
	return result
}

func (c *homeScreenComponent) createDayAmbientLight(scene *graphics.Scene) *graphics.AmbientLight {
	reflectionData := gblob.LittleEndianBlock(make([]byte, 4*2))
	reflectionData.SetUint16(0, float16.Fromfloat32(2.0).Bits())
	reflectionData.SetUint16(2, float16.Fromfloat32(2.5).Bits())
	reflectionData.SetUint16(4, float16.Fromfloat32(3.0).Bits())
	reflectionData.SetUint16(6, float16.Fromfloat32(1.0).Bits())

	reflectionTexture := c.engine.Graphics().API().CreateColorTextureCube(render.ColorTextureCubeInfo{
		GenerateMipmaps: false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA16F,
		MipmapLayers: []render.MipmapCubeLayer{
			{
				Dimension:      1,
				FrontSideData:  reflectionData,
				BackSideData:   reflectionData,
				LeftSideData:   reflectionData,
				RightSideData:  reflectionData,
				TopSideData:    reflectionData,
				BottomSideData: reflectionData,
			},
		},
	})

	refractionTexture := c.engine.Graphics().API().CreateColorTextureCube(render.ColorTextureCubeInfo{
		GenerateMipmaps: false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA16F,
		MipmapLayers: []render.MipmapCubeLayer{
			{
				Dimension:      1,
				FrontSideData:  reflectionData,
				BackSideData:   reflectionData,
				LeftSideData:   reflectionData,
				RightSideData:  reflectionData,
				TopSideData:    reflectionData,
				BottomSideData: reflectionData,
			},
		},
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

	reflectionTexture := c.engine.Graphics().API().CreateColorTextureCube(render.ColorTextureCubeInfo{
		GenerateMipmaps: false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA16F,
		MipmapLayers: []render.MipmapCubeLayer{
			{
				Dimension:      1,
				FrontSideData:  reflectionData,
				BackSideData:   reflectionData,
				LeftSideData:   reflectionData,
				RightSideData:  reflectionData,
				TopSideData:    reflectionData,
				BottomSideData: reflectionData,
			},
		},
	})

	refractionTexture := c.engine.Graphics().API().CreateColorTextureCube(render.ColorTextureCubeInfo{
		GenerateMipmaps: false,
		GammaCorrection: false,
		Format:          render.DataFormatRGBA16F,
		MipmapLayers: []render.MipmapCubeLayer{
			{
				Dimension:      1,
				FrontSideData:  reflectionData,
				BackSideData:   reflectionData,
				LeftSideData:   reflectionData,
				RightSideData:  reflectionData,
				TopSideData:    reflectionData,
				BottomSideData: reflectionData,
			},
		},
	})

	return scene.CreateAmbientLight(graphics.AmbientLightInfo{
		ReflectionTexture: reflectionTexture,
		RefractionTexture: refractionTexture,
		Position:          dprec.ZeroVec3(),
		InnerRadius:       5000,
		OuterRadius:       5000,
	})
}

func (c *homeScreenComponent) controllerDescription(controller data.Input) string {
	switch controller {
	case data.InputKeyboard:
		return "Keyboard: Uses assists. Provides an average challenge."
	case data.InputMouse:
		return "Mouse: Just point and drive. Good for a casual drive."
	case data.InputGamepad:
		return "Gamepad: No assists. Requires significant skills to control."
	default:
		return ""
	}
}

func (c *homeScreenComponent) onKeyboardClicked() {
	c.homeModel.SetInput(data.InputKeyboard)
	c.Invalidate()
}

func (c *homeScreenComponent) onMouseClicked() {
	c.homeModel.SetInput(data.InputMouse)
	c.Invalidate()
}

func (c *homeScreenComponent) onGamepadClicked() {
	c.homeModel.SetInput(data.InputGamepad)
	c.Invalidate()
}

func (c *homeScreenComponent) onDayClicked() {
	c.homeModel.SetLighting(data.LightingDay)

	// Disable night lighting
	c.scene.NightSky.SetActive(false)
	c.scene.NightAmbientLight.SetActive(false)

	// Enable day lighting
	c.scene.DaySky.SetActive(true)
	c.scene.DayAmbientLight.SetActive(true)
	c.scene.DayDirectionalLight.SetActive(true)

	c.Invalidate()
}

func (c *homeScreenComponent) onNightClicked() {
	c.homeModel.SetLighting(data.LightingNight)

	// Disable day lighting
	c.scene.DaySky.SetActive(false)
	c.scene.DayAmbientLight.SetActive(false)
	c.scene.DayDirectionalLight.SetActive(false)

	// Enable night lighting
	c.scene.NightSky.SetActive(true)
	c.scene.NightAmbientLight.SetActive(true)

	c.Invalidate()
}

func (c *homeScreenComponent) onLevelClicked(level data.Level) {
	c.homeModel.SetLevel(level)
	c.Invalidate()
}

func (c *homeScreenComponent) onPlayClicked() {
	c.homeModel.SetMode(model.HomeScreenModeLighting)
	c.Invalidate()
}

func (c *homeScreenComponent) onLicensesClicked() {
	c.appModel.SetActiveView(model.ViewNameLicenses)
}

func (c *homeScreenComponent) onCreditsClicked() {
	c.appModel.SetActiveView(model.ViewNameCredits)
}

func (c *homeScreenComponent) onExitClicked() {
	co.Window(c.Scope()).Close()
}

func (c *homeScreenComponent) onNextClicked() {
	switch c.homeModel.Mode() {
	case model.HomeScreenModeLighting:
		c.homeModel.SetMode(model.HomeScreenModeControls)
	case model.HomeScreenModeControls:
		c.homeModel.SetMode(model.HomeScreenModeLevel)
	}
	c.Invalidate()
}

func (c *homeScreenComponent) onBackClicked() {
	switch c.homeModel.Mode() {
	case model.HomeScreenModeLighting:
		c.homeModel.SetMode(model.HomeScreenModeEntry)
	case model.HomeScreenModeControls:
		c.homeModel.SetMode(model.HomeScreenModeLighting)
	case model.HomeScreenModeLevel:
		c.homeModel.SetMode(model.HomeScreenModeControls)
	}
	c.Invalidate()
}

func (c *homeScreenComponent) onStartClicked() {
	promise := model.NewLoadingPromise(
		co.Window(c.Scope()),
		data.LoadPlayData(c.engine, c.resourceSet, c.homeModel.Lighting(), c.homeModel.Input(), c.homeModel.Level().Board),
		c.playModel.SetData,
		c.errorModel.SetError,
	)
	c.loadingModel.SetState(model.LoadingState{
		Promise:         promise,
		SuccessViewName: model.ViewNamePlay,
		ErrorViewName:   model.ViewNameError,
	})
	c.appModel.SetActiveView(model.ViewNameLoading)
}
