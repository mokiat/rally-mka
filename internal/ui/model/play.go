package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
)

var (
	PlayChange     = mvc.NewChange("play")
	PlayDataChange = mvc.SubChange(PlayChange, "data")
)

func newPlay() *Play {
	return &Play{
		Observable: mvc.NewObservable(),
	}
}

type Play struct {
	mvc.Observable
	sceneData game.Promise[*data.PlayData]
}

func (h *Play) Data() game.Promise[*data.PlayData] {
	return h.sceneData
}

func (h *Play) SetData(sceneData game.Promise[*data.PlayData]) {
	h.sceneData = sceneData
	h.SignalChange(PlayDataChange)
}
