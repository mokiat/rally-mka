package physics

import "github.com/mokiat/gomath/sprec"

const driftCorrectionAmount = float32(0.01)          // TODO: Configurable?
const driftBaumgarteCorrectionAmount = float32(0.01) // TODO: Configurable?
const timeStep = float32(0.015)                      // TODO: Configurable

type SingleBodyJacobian struct {
	SlopeVelocity        sprec.Vec3
	SlopeAngularVelocity sprec.Vec3
}

func (j SingleBodyJacobian) Apply(body *Body) {
	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocity, body.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocity, body.AngularVelocity))
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocity, j.SlopeVelocity)/body.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(body.MomentOfInertia), j.SlopeAngularVelocity), j.SlopeAngularVelocity)
	lambda := lambdaUpper / lambdaLower

	body.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocity, lambda))
	body.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

func (j SingleBodyJacobian) ApplyNudge(body *Body, drift float32) {
	lambdaUpper := -driftCorrectionAmount * drift
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocity, j.SlopeVelocity)/body.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(body.MomentOfInertia), j.SlopeAngularVelocity), j.SlopeAngularVelocity)
	lambda := lambdaUpper / lambdaLower

	body.ApplyNudge(sprec.Vec3Prod(j.SlopeVelocity, lambda))
	body.ApplyAngularNudge(sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

type DoubleBodyJacobian struct {
	SlopeVelocityFirst         sprec.Vec3
	SlopeAngularVelocityFirst  sprec.Vec3
	SlopeVelocitySecond        sprec.Vec3
	SlopeAngularVelocitySecond sprec.Vec3
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

func (j DoubleBodyJacobian) ApplySoft(firstBody, secondBody *Body, force, gamma float32) {
	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocityFirst, firstBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstBody.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondBody.AngularVelocity) +
		force*gamma)
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

func (j DoubleBodyJacobian) ApplySoft2(firstBody, secondBody *Body, initialLambda, gamma float32) {
	// effectiveVelocity := -(sprec.Vec3Dot(j.SlopeVelocityFirst, firstBody.Velocity) +
	// 	sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstBody.AngularVelocity) +
	// 	sprec.Vec3Dot(j.SlopeVelocitySecond, secondBody.Velocity) +
	// 	sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondBody.AngularVelocity))
	// effectiveMass := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstBody.Mass +
	// 	sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstBody.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
	// 	sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondBody.Mass +
	// 	sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondBody.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := -j.EffectiveVelocity(firstBody, secondBody)/j.EffectiveMass(firstBody, secondBody) + initialLambda*gamma

	firstBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleBodyJacobian) ApplyImpulseBeta(firstBody, secondBody *Body, beta, drift float32) {
	timestep := float32(0.015) // FIXME
	lambda := (-j.EffectiveVelocity(firstBody, secondBody) - beta*drift/timestep) / j.EffectiveMass(firstBody, secondBody)
	j.applyLambda(firstBody, secondBody, lambda)
}

func (j DoubleBodyJacobian) ApplyImpulseBetaGamma(firstBody, secondBody *Body, beta, drift, gamma float32) {
	timestep := float32(0.015) // FIXME
	effMass := j.EffectiveMass(firstBody, secondBody)
	lambda := j.EffectiveVelocity(firstBody, secondBody)
	finalLambda := -gamma * (lambda + beta*drift/timestep) / effMass
	j.applyLambda(firstBody, secondBody, finalLambda)
}

func (j DoubleBodyJacobian) applyLambda(firstBody, secondBody *Body, lambda float32) {
	firstBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondBody.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondBody.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleBodyJacobian) EffectiveVelocity(firstBody, secondBody *Body) float32 {
	return sprec.Vec3Dot(j.SlopeVelocityFirst, firstBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstBody.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondBody.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondBody.AngularVelocity)
}

func (j DoubleBodyJacobian) EffectiveMass(firstBody, secondBody *Body) float32 {
	return sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstBody.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondBody.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondBody.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
}

func (j DoubleBodyJacobian) ApplyNew(firstBody, secondBody *Body, velocityDrift float32) {
	lambdaUpper := -velocityDrift
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

func (j DoubleBodyJacobian) ApplyBaumgarte(firstBody, secondBody *Body, drift float32) {
	lambdaUpper := -driftBaumgarteCorrectionAmount * drift / timeStep
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
