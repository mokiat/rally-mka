package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type Coilover struct {
	ecs.NilConstraint
	Target               *ecs.Entity
	TargetRelativeOffset sprec.Vec3
	Entity               *ecs.Entity
	Strength             float32
}

func (c Coilover) ApplyForces() {
	targetTransformComp := c.Target.Transform
	targetAnchorRelativePosition := targetTransformComp.Orientation.MulVec3(c.TargetRelativeOffset)

	targetMotionComp := c.Target.Motion
	entityMotionComp := c.Entity.Motion

	deltaVelocity := sprec.Vec3Diff(entityMotionComp.Velocity, sprec.Vec3Sum(targetMotionComp.Velocity, sprec.Vec3Cross(targetMotionComp.AngularVelocity, targetAnchorRelativePosition)))
	if deltaVelocity.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaVelocity)

	force := sprec.Vec3Prod(jacobian, deltaVelocity.Length()*c.Strength)
	targetMotionComp.ApplyOffsetForce(targetAnchorRelativePosition, sprec.InverseVec3(force))
	entityMotionComp.ApplyForce(sprec.InverseVec3(force))
}
