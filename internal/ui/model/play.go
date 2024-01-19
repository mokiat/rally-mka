package model

import (
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/util/async"
	"github.com/mokiat/rally-mka/internal/game/data"
)

func NewPlay(eventBus *mvc.EventBus) *Play {
	return &Play{}
}

type Play struct {
	sceneData async.Promise[*data.PlayData]
}

func (h *Play) Data() async.Promise[*data.PlayData] {
	return h.sceneData
}

func (h *Play) SetData(sceneData async.Promise[*data.PlayData]) {
	h.sceneData = sceneData
}
