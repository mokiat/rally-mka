package ecscomp

import (
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/physics"
)

type Physics struct {
	Body *physics.Body
}

func SetPhysics(entity *ecs.Entity, component *Physics) {
	entity.SetComponent(PhysicsComponentID, component)
}

func GetPhysics(entity *ecs.Entity) *Physics {
	component := entity.Component(PhysicsComponentID)
	if component == nil {
		return nil
	}
	return component.(*Physics)
}
