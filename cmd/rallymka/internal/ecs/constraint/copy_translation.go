package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type CopyTranslation struct {
	ecs.NilConstraint
	Target         *ecs.Entity
	Entity         *ecs.Entity
	RelativeOffset sprec.Vec3
	SkipX          bool
	SkipY          bool
	SkipZ          bool
}

func (t CopyTranslation) ApplyCorrectionImpulses() {
	targetTransformComp := t.Target.Transform
	entityTransformComp := t.Entity.Transform

	targetAnchorRelativePosition := targetTransformComp.Orientation.MulVec3(t.RelativeOffset)
	targetAnchorPosition := sprec.Vec3Sum(targetTransformComp.Position, targetAnchorRelativePosition)
	deltaPosition := sprec.Vec3Diff(entityTransformComp.Position, targetAnchorPosition)
	if deltaPosition.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaPosition)
	if t.SkipY {
		jacobian = sprec.Vec3Diff(jacobian, sprec.Vec3Prod(targetTransformComp.Orientation.OrientationY(), sprec.Vec3Dot(jacobian, targetTransformComp.Orientation.OrientationY())))
	}

	targetMotionComp := t.Target.Motion
	entityMotionComp := t.Entity.Motion

	targetAnchorVelocity := sprec.Vec3Sum(targetMotionComp.Velocity, sprec.Vec3Cross(targetMotionComp.AngularVelocity, targetAnchorRelativePosition))
	entityVelocity := entityMotionComp.Velocity
	deltaVelocity := sprec.Vec3Diff(entityVelocity, targetAnchorVelocity)

	targetEffectiveMass := 1.0 / ((1.0 / targetMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(targetMotionComp.MomentOfInertia), sprec.Vec3Cross(targetAnchorRelativePosition, jacobian)), sprec.Vec3Cross(targetAnchorRelativePosition, jacobian)))
	entityEffectiveMass := entityMotionComp.Mass
	totalMass := targetEffectiveMass * entityEffectiveMass / (targetEffectiveMass + entityEffectiveMass)

	impulseStrength := totalMass * sprec.Vec3Dot(jacobian, deltaVelocity)
	impulse := sprec.Vec3Prod(jacobian, impulseStrength)
	targetMotionComp.ApplyOffsetImpulse(targetAnchorRelativePosition, impulse)
	entityMotionComp.ApplyImpulse(sprec.InverseVec3(impulse))
}

func (t CopyTranslation) ApplyCorrectionTranslations() {
	targetTransformComp := t.Target.Transform
	entityTransformComp := t.Entity.Transform

	targetAnchorRelativePosition := targetTransformComp.Orientation.MulVec3(t.RelativeOffset)
	targetAnchorPosition := sprec.Vec3Sum(targetTransformComp.Position, targetAnchorRelativePosition)
	deltaPosition := sprec.Vec3Diff(entityTransformComp.Position, targetAnchorPosition)
	if deltaPosition.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaPosition)
	if t.SkipY {
		jacobian = sprec.Vec3Diff(jacobian, sprec.Vec3Prod(targetTransformComp.Orientation.OrientationY(), sprec.Vec3Dot(jacobian, targetTransformComp.Orientation.OrientationY())))
	}

	targetMotionComp := t.Target.Motion
	entityMotionComp := t.Entity.Motion

	targetEffectiveMass := 1.0 / ((1.0 / targetMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(targetMotionComp.MomentOfInertia), sprec.Vec3Cross(targetAnchorRelativePosition, jacobian)), sprec.Vec3Cross(targetAnchorRelativePosition, jacobian)))
	entityEffectiveMass := entityMotionComp.Mass
	totalMass := targetEffectiveMass * entityEffectiveMass / (targetEffectiveMass + entityEffectiveMass)

	nudgeStrength := totalMass * sprec.Vec3Dot(jacobian, deltaPosition)
	nudge := sprec.Vec3Prod(jacobian, nudgeStrength)

	targetTransformComp.Translate(sprec.Vec3Quot(nudge, targetMotionComp.Mass))
	entityTransformComp.Translate(sprec.InverseVec3(sprec.Vec3Quot(nudge, entityMotionComp.Mass)))

	targetTransformComp.Rotate(sprec.Mat3Vec3Prod(sprec.InverseMat3(targetMotionComp.MomentOfInertia), sprec.Vec3Cross(targetAnchorRelativePosition, nudge)))
}
