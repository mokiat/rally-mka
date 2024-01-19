package model

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
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

func NewHome(eventBus *mvc.EventBus) *Home {
	return &Home{
		eventBus:    eventBus,
		controller:  data.ControllerKeyboard,
		environment: data.EnvironmentDay,
	}
}

type Home struct {
	eventBus    *mvc.EventBus
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
}

func (h *Home) Scene() *HomeScene {
	return h.scene
}

func (h *Home) SetScene(scene *HomeScene) {
	h.scene = scene
}

func (h *Home) Controller() data.Controller {
	return h.controller
}

func (h *Home) SetController(controller data.Controller) {
	if controller != h.controller {
		h.controller = controller
		h.eventBus.Notify(ControllerChangedEvent{})
	}
}

func (h *Home) Environment() data.Environment {
	return h.environment
}

func (h *Home) SetEnvironment(environment data.Environment) {
	if environment != h.environment {
		h.environment = environment
		h.eventBus.Notify(EnvironmentChangedEvent{})
	}
}

type ControllerChangedEvent struct{}

type EnvironmentChangedEvent struct{}
