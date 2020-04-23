package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type Spring struct {
	ecs.NilConstraint
	Target               *ecs.Entity
	TargetRelativeOffset sprec.Vec3
	Entity               *ecs.Entity
	Length               float32
	Stiffness            float32
}

func (s Spring) ApplyForces() {
	targetTransformComp := s.Target.Transform
	entityTransformComp := s.Entity.Transform

	targetAnchorRelativePosition := sprec.QuatVec3Rotation(targetTransformComp.Orientation, s.TargetRelativeOffset)
	targetAnchorPosition := sprec.Vec3Sum(targetTransformComp.Position, targetAnchorRelativePosition)
	entityPosition := entityTransformComp.Position

	deltaPosition := sprec.Vec3Diff(entityPosition, targetAnchorPosition)
	extendDistance := deltaPosition.Length() - s.Length
	if sprec.Abs(extendDistance) < 0.0001 {
		return
	}
	if deltaPosition.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaPosition)

	targetMotionComp := s.Target.Motion
	entityMotionComp := s.Entity.Motion

	force := sprec.Vec3Prod(jacobian, extendDistance*s.Stiffness)
	targetMotionComp.ApplyOffsetForce(targetAnchorRelativePosition, force)
	entityMotionComp.ApplyForce(sprec.InverseVec3(force))
}
