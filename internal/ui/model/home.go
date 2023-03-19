package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
)

var (
	HomeChange      = mvc.NewChange("home")
	HomeDataChange  = mvc.SubChange(HomeChange, "data")
	HomeSceneChange = mvc.SubChange(HomeChange, "scene")
)

func newHome() *Home {
	return &Home{
		Observable: mvc.NewObservable(),
	}
}

type Home struct {
	mvc.Observable
	sceneData game.Promise[*data.HomeData]
	scene     *game.Scene
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
