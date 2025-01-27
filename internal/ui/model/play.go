package model

import (
	"github.com/mokiat/rally-mka/internal/game/data"
)

func NewPlayModel() *PlayModel {
	return &PlayModel{}
}

type PlayModel struct {
	sceneData *data.PlayData
}

func (h *PlayModel) Data() *data.PlayData {
	return h.sceneData
}

func (h *PlayModel) SetData(sceneData *data.PlayData) {
	h.sceneData = sceneData
}
