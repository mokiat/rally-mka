package ecscomp

import "github.com/mokiat/lacking/game/ecs"

func SetPlayerControl(entity *ecs.Entity, component *PlayerControl) {
	entity.SetComponent(PlayerControlComponentID, component)
}

func GetPlayerControl(entity *ecs.Entity) *PlayerControl {
	component := entity.Component(PlayerControlComponentID)
	if component == nil {
		return nil
	}
	return component.(*PlayerControl)
}

type PlayerControl struct{}
