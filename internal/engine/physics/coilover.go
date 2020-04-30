package physics

import (
	"github.com/mokiat/gomath/sprec"
)

type CoiloverConstraint struct {
	NilConstraint
	FirstBody       *Body
	FirstBodyAnchor sprec.Vec3
	SecondBody      *Body
	Length          float32
	Frequency       float32
	DampingRatio    float32
}

func (c CoiloverConstraint) ApplyImpulse() {
	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAnchor)
	firstAnchorWS := sprec.Vec3Sum(c.FirstBody.Position, firstRadiusWS)
	secondAnchorWS := c.SecondBody.Position
	deltaPosition := sprec.Vec3Diff(secondAnchorWS, firstAnchorWS)
	drift := deltaPosition.Length()
	normal := sprec.BasisXVec3()
	if drift > 0.0001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	jacobian := DoubleBodyJacobian{
		SlopeVelocityFirst: sprec.NewVec3(
			-normal.X,
			-normal.Y,
			-normal.Z,
		),
		SlopeAngularVelocityFirst: sprec.NewVec3(
			-(normal.Z*firstRadiusWS.Y - normal.Y*firstRadiusWS.Z),
			-(normal.X*firstRadiusWS.Z - normal.Z*firstRadiusWS.X),
			-(normal.Y*firstRadiusWS.X - normal.X*firstRadiusWS.Y),
		),
		SlopeVelocitySecond: sprec.NewVec3(
			normal.X,
			normal.Y,
			normal.Z,
		),
		SlopeAngularVelocitySecond: sprec.ZeroVec3(),
	}

	beta := float32(0.65)
	gamma := float32(0.003)
	// beta := float32(0.1)
	// gamma := float32(0.0001)
	jacobian.ApplyImpulseBetaGamma(c.FirstBody, c.SecondBody, beta, drift, gamma)
}
