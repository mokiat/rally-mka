package physics

import "github.com/mokiat/gomath/sprec"

type FixedTranslationConstraint struct {
	NilConstraint
	Fixture sprec.Vec3
	Body    *Body
}

func (c FixedTranslationConstraint) ApplyImpulse() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.CorrectVelocity(c.Body)
	}
}

func (c FixedTranslationConstraint) ApplyNudge() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.CorrectPosition(c.Body, result.Drift)
	}
}

func (c FixedTranslationConstraint) Calculate() FixedTranslationConstraintResult {
	deltaPosition := sprec.Vec3Diff(c.Body.Position, c.Fixture)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return FixedTranslationConstraintResult{
		Jacobian: Jacobian{
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

type FixedTranslationConstraintResult struct {
	Jacobian Jacobian
	Drift    float32
}
