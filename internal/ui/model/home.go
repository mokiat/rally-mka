package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
)

var (
	HomeChange            = mvc.NewChange("home")
	HomeDataChange        = mvc.SubChange(HomeChange, "data")
	HomeSceneChange       = mvc.SubChange(HomeChange, "scene")
	HomeControllerChange  = mvc.SubChange(HomeChange, "controller")
	HomeEnvironmentChange = mvc.SubChange(HomeChange, "environment")
)

type Controller string

const (
	ControllerKeyboard Controller = "keyboard"
	ControllerMouse    Controller = "mouse"
	ControllerGamepad  Controller = "gamepad"
)

type Environment string

const (
	EnvironmentDay   Environment = "day"
	EnvironmentNight Environment = "night"
)

func newHome() *Home {
	return &Home{
		Observable:  mvc.NewObservable(),
		controller:  ControllerKeyboard,
		environment: EnvironmentDay,
	}
}

type Home struct {
	mvc.Observable
	sceneData   game.Promise[*data.HomeData]
	scene       *game.Scene
	controller  Controller
	environment Environment
}

func (h *Home) Data() game.Promise[*data.HomeData] {
	return h.sceneData
}

func (h *Home) SetData(sceneData game.Promise[*data.HomeData]) {
	h.sceneData = sceneData
	h.SignalChange(HomeDataChange)
}

func (h *Home) Scene() *game.Scene {
	return h.scene
}

func (h *Home) SetScene(scene *game.Scene) {
	h.scene = scene
	h.SignalChange(HomeSceneChange)
}

func (h *Home) Controller() Controller {
	return h.controller
}

func (h *Home) SetController(controller Controller) {
	h.controller = controller
	h.SignalChange(HomeControllerChange)
}

func (h *Home) Environment() Environment {
	return h.environment
}

func (h *Home) SetEnvironment(environment Environment) {
	h.environment = environment
	h.SignalChange(HomeEnvironmentChange)
}
