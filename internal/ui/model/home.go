package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/rally-mka/internal/game/data"
)

func NewHomeModel() *HomeModel {
	return &HomeModel{
		mode:     HomeScreenModeEntry,
		input:    data.InputKeyboard,
		lighting: data.LightingDay,
		level:    data.Levels[0],
	}
}

type HomeModel struct {
	sceneData *data.HomeData
	mode      HomeScreenMode
	input     data.Input
	lighting  data.Lighting
	level     data.Level
	scene     *HomeScene
}

func (h *HomeModel) Data() *data.HomeData {
	return h.sceneData
}

func (h *HomeModel) SetData(sceneData *data.HomeData) {
	h.sceneData = sceneData
}

func (h *HomeModel) Scene() *HomeScene {
	return h.scene
}

func (h *HomeModel) SetScene(scene *HomeScene) {
	h.scene = scene
}

func (h *HomeModel) Mode() HomeScreenMode {
	return h.mode
}

func (h *HomeModel) SetMode(mode HomeScreenMode) {
	h.mode = mode
}

func (h *HomeModel) Input() data.Input {
	return h.input
}

func (h *HomeModel) SetInput(input data.Input) {
	h.input = input
}

func (h *HomeModel) Lighting() data.Lighting {
	return h.lighting
}

func (h *HomeModel) SetLighting(lighting data.Lighting) {
	h.lighting = lighting
}

func (h *HomeModel) Level() data.Level {
	return h.level
}

func (h *HomeModel) SetLevel(level data.Level) {
	h.level = level
}

type HomeScene struct {
	Scene *game.Scene

	DaySky              *graphics.Sky
	DayAmbientLight     *graphics.AmbientLight
	DayDirectionalLight *graphics.DirectionalLight

	NightSky          *graphics.Sky
	NightAmbientLight *graphics.AmbientLight
}

type HomeScreenMode uint8

const (
	HomeScreenModeEntry HomeScreenMode = iota
	HomeScreenModeLighting
	HomeScreenModeControls
	HomeScreenModeLevel
)
