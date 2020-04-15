package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type CopyRotation struct {
	ecs.NilConstraint
	Target *ecs.Entity
	Entity *ecs.Entity
}

func (t CopyRotation) ApplyCorrectionImpulses() {
	targetMotionComp := t.Target.Motion
	entityMotionComp := t.Entity.Motion

	targetAngularVelocity := targetMotionComp.AngularVelocity
	entityAngularVelocity := entityMotionComp.AngularVelocity
	deltaAngularVelocity := sprec.Vec3Diff(entityAngularVelocity, targetAngularVelocity)

	targetEffectiveMass := sprec.Mat3Vec3Prod(targetMotionComp.MomentOfInertia, deltaAngularVelocity).Length()
	entityEffectiveMass := sprec.Mat3Vec3Prod(entityMotionComp.MomentOfInertia, deltaAngularVelocity).Length()
	totalMass := targetEffectiveMass * entityEffectiveMass / (targetEffectiveMass + entityEffectiveMass)

	impulse := sprec.Vec3Prod(deltaAngularVelocity, totalMass)
	targetMotionComp.ApplyAngularImpulse(impulse)
	entityMotionComp.ApplyAngularImpulse(sprec.InverseVec3(impulse))
}

func (t CopyRotation) ApplyCorrectionTranslations() {
	// TODO
}
