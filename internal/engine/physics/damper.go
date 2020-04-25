package physics

import "github.com/mokiat/gomath/sprec"

type DamperConstraint struct {
	NilConstraint
	FirstBody       *Body
	FirstBodyAnchor sprec.Vec3
	SecondBody      *Body
	Strength        float32
}

func (c DamperConstraint) ApplyForce() {
	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAnchor)
	deltaVelocity := sprec.Vec3Diff(c.SecondBody.Velocity, sprec.Vec3Sum(c.FirstBody.Velocity, sprec.Vec3Cross(c.FirstBody.AngularVelocity, firstRadiusWS)))
	if deltaVelocity.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaVelocity)

	force := sprec.Vec3Prod(jacobian, deltaVelocity.Length()*c.Strength)
	c.FirstBody.ApplyOffsetForce(firstRadiusWS, sprec.InverseVec3(force))
	c.SecondBody.ApplyForce(sprec.InverseVec3(force))
}
