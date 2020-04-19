package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type FixedTranslation struct {
	ecs.NilConstraint
	Entity   *ecs.Entity
	Position sprec.Vec3
}

func (c FixedTranslation) ApplyCorrectionImpulses() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(c.Entity)
	}
}

func (c FixedTranslation) ApplyCorrectionTranslations() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(c.Entity, result.Drift)
	}
}

func (c FixedTranslation) Calculate() FixedTranslationResult {
	transformComp := c.Entity.Transform
	deltaPosition := sprec.Vec3Diff(transformComp.Position, c.Position)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return FixedTranslationResult{
		Jacobian: ecs.SingleEntityJacobian{
			SlopeVelocity: sprec.NewVec3(
				normal.X,
				normal.Y,
				normal.Z,
			),
			SlopeAngularVelocity: sprec.ZeroVec3(),
		},
		Drift: deltaPosition.Length(),
	}
}

type FixedTranslationResult struct {
	Jacobian ecs.SingleEntityJacobian
	Drift    float32
}
