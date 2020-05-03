package physics

import "github.com/mokiat/gomath/sprec"

const driftCorrectionAmount = float32(0.01) // TODO: Configurable?

type Jacobian struct {
	SlopeVelocity        sprec.Vec3
	SlopeAngularVelocity sprec.Vec3
}

func (j Jacobian) EffectiveVelocity(body *Body) float32 {
	return sprec.Vec3Dot(j.SlopeVelocity, body.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocity, body.AngularVelocity)
}

func (j Jacobian) InverseEffectiveMass(body *Body) float32 {
	return sprec.Vec3Dot(j.SlopeVelocity, j.SlopeVelocity)/body.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(body.MomentOfInertia), j.SlopeAngularVelocity), j.SlopeAngularVelocity)
}

func (j Jacobian) ApplyImpulse(body *Body, lambda float32) {
	body.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocity, lambda))
	body.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

func (j Jacobian) ApplyNudge(body *Body, lambda float32) {
	body.ApplyNudge(sprec.Vec3Prod(j.SlopeVelocity, lambda))
	body.ApplyAngularNudge(sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

func (j Jacobian) CorrectVelocity(body *Body) {
	lambda := -j.EffectiveVelocity(body) / j.InverseEffectiveMass(body)
	j.ApplyImpulse(body, lambda)
}

func (j Jacobian) CorrectPosition(body *Body, drift float32) {
	lambda := -driftCorrectionAmount * drift / j.InverseEffectiveMass(body)
	j.ApplyNudge(body, lambda)
}

type DoubleBodyJacobian struct {
	SlopeVelocityFirst         sprec.Vec3
	SlopeAngularVelocityFirst  sprec.Vec3
	SlopeVelocitySecond        sprec.Vec3
	SlopeAngularVelocitySecond sprec.Vec3
}

func (j DoubleBodyJacobian) EffectiveVelocity(firstBody, secondBody *Body) float32 {
	return sprec.Vec3Dot(j.SlopeVelocityFirst, firstBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstBody.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondBody.AngularVelocity)
}

func (j DoubleBodyJacobian) EffectiveMass(firstBody, secondBody *Body) float32 {
	inverseMass := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstBody.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondBody.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	return 1.0 / inverseMass
}

func (j DoubleBodyJacobian) ApplyLambda(firstBody, secondBody *Body, lambda float32) {
	firstBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleBodyJacobian) Apply(firstBody, secondBody *Body) {
	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocityFirst, firstBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstBody.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondBody.AngularVelocity))
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstBody.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondBody.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleBodyJacobian) ApplyNudge(firstBody, secondBody *Body, drift float32) {
	lambdaUpper := -driftCorrectionAmount * drift
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstBody.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondBody.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstBody.ApplyNudge(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstBody.ApplyAngularNudge(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondBody.ApplyNudge(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondBody.ApplyAngularNudge(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}
