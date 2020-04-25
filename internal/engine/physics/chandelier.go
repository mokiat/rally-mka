package physics

import "github.com/mokiat/gomath/sprec"

type ChandelierConstraint struct {
	NilConstraint
	Fixture    sprec.Vec3
	Body       *Body
	BodyAnchor sprec.Vec3
	Length     float32
}

func (c ChandelierConstraint) ApplyImpulse() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(c.Body)
	}
}

func (c ChandelierConstraint) ApplyNudge() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(c.Body, result.Drift)
	}
}

func (c ChandelierConstraint) Calculate() ChandelierConstraintResult {
	anchorWS := sprec.Vec3Sum(c.Body.Position, sprec.QuatVec3Rotation(c.Body.Orientation, c.BodyAnchor))
	radiusWS := sprec.Vec3Diff(anchorWS, c.Body.Position)
	deltaPosition := sprec.Vec3Diff(anchorWS, c.Fixture)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return ChandelierConstraintResult{
		Jacobian: SingleBodyJacobian{
			SlopeVelocity: sprec.NewVec3(
				normal.X,
				normal.Y,
				normal.Z,
			),
			SlopeAngularVelocity: sprec.NewVec3(
				normal.Z*radiusWS.Y-normal.Y*radiusWS.Z,
				normal.X*radiusWS.Z-normal.Z*radiusWS.X,
				normal.Y*radiusWS.X-normal.X*radiusWS.Y,
			),
		},
		Drift: deltaPosition.Length() - c.Length,
	}
}

type ChandelierConstraintResult struct {
	Jacobian SingleBodyJacobian
	Drift    float32
}
