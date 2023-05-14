package model

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
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

type HomeScene struct {
	Scene *game.Scene

	DaySkyColor         sprec.Vec3
	DayAmbientLight     *graphics.AmbientLight
	DayDirectionalLight *graphics.DirectionalLight

	NightSkyColor     sprec.Vec3
	NightAmbientLight *graphics.AmbientLight
	NightSpotLight    *graphics.SpotLight
}

func newHome() *Home {
	return &Home{
		Observable:  mvc.NewObservable(),
		controller:  data.ControllerKeyboard,
		environment: data.EnvironmentDay,
	}
}

type Home struct {
	mvc.Observable
	sceneData   game.Promise[*data.HomeData]
	scene       *HomeScene
	controller  data.Controller
	environment data.Environment
}

func (h *Home) Data() game.Promise[*data.HomeData] {
	return h.sceneData
}

func (h *Home) SetData(sceneData game.Promise[*data.HomeData]) {
	h.sceneData = sceneData
	h.SignalChange(HomeDataChange)
}

func (h *Home) Scene() *HomeScene {
	return h.scene
}

func (h *Home) SetScene(scene *HomeScene) {
	h.scene = scene
	h.SignalChange(HomeSceneChange)
}

func (h *Home) Controller() data.Controller {
	return h.controller
}

func (h *Home) SetController(controller data.Controller) {
	h.controller = controller
	h.SignalChange(HomeControllerChange)
}

func (h *Home) Environment() data.Environment {
	return h.environment
}

func (h *Home) SetEnvironment(environment data.Environment) {
	h.environment = environment
	h.SignalChange(HomeEnvironmentChange)
}
