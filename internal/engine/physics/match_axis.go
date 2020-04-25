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
	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAxis)
	secondRadiusWS := sprec.QuatVec3Rotation(c.SecondBody.Orientation, c.SecondBodyAxis)
	deltaPosition := sprec.Vec3Diff(secondRadiusWS, firstRadiusWS)
	normal := sprec.BasisZVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}
	return MatchAxisConstraintResult{
		Jacobian: DoubleBodyJacobian{
			SlopeVelocityFirst: sprec.ZeroVec3(),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normal.Z*firstRadiusWS.Y - normal.Y*firstRadiusWS.Z),
				-(normal.X*firstRadiusWS.Z - normal.Z*firstRadiusWS.X),
				-(normal.Y*firstRadiusWS.X - normal.X*firstRadiusWS.Y),
			),
			SlopeVelocitySecond: sprec.ZeroVec3(),
			SlopeAngularVelocitySecond: sprec.NewVec3(
				(normal.Z*secondRadiusWS.Y - normal.Y*secondRadiusWS.Z),
				(normal.X*secondRadiusWS.Z - normal.Z*secondRadiusWS.X),
				(normal.Y*secondRadiusWS.X - normal.X*secondRadiusWS.Y),
			),
		},
		Drift: deltaPosition.Length(),
	}
}

type MatchAxisConstraintResult struct {
	Jacobian DoubleBodyJacobian
	Drift    float32
}
