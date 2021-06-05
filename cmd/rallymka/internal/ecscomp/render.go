package ecscomp

import (
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/render"
)

func SetRender(entity *ecs.Entity, component *Render) {
	entity.SetComponent(RenderComponentID, component)
}

func GetRender(entity *ecs.Entity) *Render {
	component := entity.Component(RenderComponentID)
	if component == nil {
		return nil
	}
	return component.(*Render)
}

type Render struct {
	Renderable *render.Renderable
}
