package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type FixedPosition struct {
	ecs.NilConstraint
	Entity   *ecs.Entity
	Position sprec.Vec3
}

var _ ecs.Constraint = FixedPosition{}

func (c FixedPosition) ApplyCorrectionImpulses() {
	motionComp := c.Entity.Motion

	velocityCorrection := sprec.InverseVec3(motionComp.Velocity)
	motionComp.ApplyImpulse(sprec.Vec3Prod(velocityCorrection, motionComp.Mass))
}

func (c FixedPosition) ApplyCorrectionTranslations() {
	transformComp := c.Entity.Transform
	transformComp.Position = c.Position
}
