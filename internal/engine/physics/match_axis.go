package physics

import "github.com/mokiat/gomath/sprec"

type MatchAxisConstraint struct {
	NilConstraint
	FirstBody      *Body
	FirstBodyAxis  sprec.Vec3
	SecondBody     *Body
	SecondBodyAxis sprec.Vec3
}

func (c MatchAxisConstraint) ApplyImpulse() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(c.FirstBody, c.SecondBody)
	}
}

func (c MatchAxisConstraint) ApplyNudge() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(c.FirstBody, c.SecondBody, result.Drift)
	}
}

func (c MatchAxisConstraint) Calculate() MatchAxisConstraintResult {
	// FIXME: Does not handle when axis are pointing in opposite directions
	firstAxisWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAxis)
	secondAxisWS := sprec.QuatVec3Rotation(c.SecondBody.Orientation, c.SecondBodyAxis)
	cross := sprec.Vec3Cross(firstAxisWS, secondAxisWS)
	return MatchAxisConstraintResult{
		Jacobian: DoubleBodyJacobian{
			SlopeVelocityFirst:         sprec.ZeroVec3(),
			SlopeAngularVelocityFirst:  sprec.InverseVec3(cross),
			SlopeVelocitySecond:        sprec.ZeroVec3(),
			SlopeAngularVelocitySecond: cross,
		},
		Drift: cross.Length(),
	}
}

type MatchAxisConstraintResult struct {
	Jacobian DoubleBodyJacobian
	Drift    float32
}
