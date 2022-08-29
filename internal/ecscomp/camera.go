package ecscomp

import (
	"github.com/mokiat/gomath/dprec"
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
	AnchorPosition dprec.Vec3
	AnchorDistance float64
	CameraDistance float64
	Camera         *graphics.Camera
}
