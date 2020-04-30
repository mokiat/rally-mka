package physics

import (
	"github.com/mokiat/gomath/sprec"
)

type DamperConstraint struct {
	NilConstraint
	FirstBody       *Body
	FirstBodyAnchor sprec.Vec3
	SecondBody      *Body
	Strength        float32
	initialLambda   float32
}

// func (c DamperConstraint) ApplyForce() {
// 	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAnchor)
// 	deltaVelocity := sprec.Vec3Diff(c.SecondBody.Velocity, sprec.Vec3Sum(c.FirstBody.Velocity, sprec.Vec3Cross(c.FirstBody.AngularVelocity, firstRadiusWS)))
// 	if deltaVelocity.SqrLength() < 0.00001 {
// 		return
// 	}
// 	jacobian := sprec.UnitVec3(deltaVelocity)

// 	force := sprec.Vec3Prod(jacobian, deltaVelocity.Length()*c.Strength)
// 	c.FirstBody.ApplyOffsetForce(firstRadiusWS, sprec.InverseVec3(force))
// 	c.SecondBody.ApplyForce(sprec.InverseVec3(force))
// }

// func (c *DamperConstraint) ApplyForce() {
// 	c.initialLambda = 0.0
// }

// func (c *DamperConstraint) ApplyImpulse() {
// 	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAnchor)
// 	deltaVelocity := sprec.Vec3Diff(c.SecondBody.Velocity, sprec.Vec3Sum(c.FirstBody.Velocity, sprec.Vec3Cross(c.FirstBody.AngularVelocity, firstRadiusWS)))
// 	if deltaVelocity.SqrLength() < 0.00001 {
// 		return
// 	}
// 	normal := sprec.UnitVec3(deltaVelocity)

// 	jacobian := DoubleBodyJacobian{
// 		SlopeVelocityFirst: sprec.NewVec3(
// 			-normal.X,
// 			-normal.Y,
// 			-normal.Z,
// 		),
// 		SlopeAngularVelocityFirst: sprec.NewVec3(
// 			-(normal.Z*firstRadiusWS.Y - normal.Y*firstRadiusWS.Z),
// 			-(normal.X*firstRadiusWS.Z - normal.Z*firstRadiusWS.X),
// 			-(normal.Y*firstRadiusWS.X - normal.X*firstRadiusWS.Y),
// 		),
// 		SlopeVelocitySecond: sprec.NewVec3(
// 			normal.X,
// 			normal.Y,
// 			normal.Z,
// 		),
// 		SlopeAngularVelocitySecond: sprec.ZeroVec3(),
// 	}
// 	if sprec.Abs(c.initialLambda) < 0.00001 {
// 		timeStep := float32(0.015)
// 		c.initialLambda = jacobian.EffectiveVelocity(c.FirstBody, c.SecondBody) / jacobian.EffectiveMass(c.FirstBody, c.SecondBody) / timeStep
// 		fmt.Printf("initial lambda: %f\n", c.initialLambda)
// 	}

// 	force := sprec.Vec3Dot(normal, deltaVelocity) * c.Strength
// 	jacobian.ApplySoft(c.FirstBody, c.SecondBody, force, -1.111)
// 	// jacobian.ApplySoft2(c.FirstBody, c.SecondBody, c.initialLambda, 0.995)
// 	// jacobian.ApplySoft2(c.FirstBody, c.SecondBody, c.initialLambda, -0.0005)

// 	// c.FirstBody.ApplyOffsetForce(firstRadiusWS, sprec.InverseVec3(force))
// 	// c.SecondBody.ApplyForce(sprec.InverseVec3(force))
// }
