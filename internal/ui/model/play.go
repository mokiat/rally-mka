package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/game/data"
)

func NewPlay(eventBus *mvc.EventBus) *Play {
	return &Play{}
}

type Play struct {
	sceneData game.Promise[*data.PlayData]
}

func (h *Play) Data() game.Promise[*data.PlayData] {
	return h.sceneData
}

func (h *Play) SetData(sceneData game.Promise[*data.PlayData]) {
	h.sceneData = sceneData
}
