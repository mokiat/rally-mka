package ecscomp

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
)

func SetCameraStand(entity *ecs.Entity, component *CameraStand) {
	entity.SetComponent(CameraStandComponentID, component)
}

func GetCameraStand(entity *ecs.Entity) *CameraStand {
	component := entity.Component(CameraStandComponentID)
	if component == nil {
		return nil
	}
	return component.(*CameraStand)
}

type CameraStand struct {
	Target         *game.Node
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *graphics.Camera
}
