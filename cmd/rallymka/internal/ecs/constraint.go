package ecs

import "github.com/mokiat/gomath/sprec"

const driftCorrectionAmount = float32(0.01)          // TODO: Configurable?
const driftBaumgarteCorrectionAmount = float32(0.01) // TODO: Configurable?
const timeStep = float32(0.015)                      // TODO: Configurable

type SingleEntityJacobian struct {
	SlopeVelocity        sprec.Vec3
	SlopeAngularVelocity sprec.Vec3
}

func (j SingleEntityJacobian) Apply(entity *Entity) {
	motionComp := entity.Motion

	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocity, motionComp.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocity, motionComp.AngularVelocity))
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocity, j.SlopeVelocity)/motionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(motionComp.MomentOfInertia), j.SlopeAngularVelocity), j.SlopeAngularVelocity)
	lambda := lambdaUpper / lambdaLower

	motionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocity, lambda))
	motionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

func (j SingleEntityJacobian) ApplyNudge(entity *Entity, drift float32) {
	motionComp := entity.Motion

	lambdaUpper := -driftCorrectionAmount * drift
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocity, j.SlopeVelocity)/motionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(motionComp.MomentOfInertia), j.SlopeAngularVelocity), j.SlopeAngularVelocity)
	lambda := lambdaUpper / lambdaLower

	applyNudge(entity, sprec.Vec3Prod(j.SlopeVelocity, lambda))
	applyAngularNudge(entity, sprec.Vec3Prod(j.SlopeAngularVelocity, lambda))
}

type DoubleEntityJacobian struct {
	SlopeVelocityFirst         sprec.Vec3
	SlopeAngularVelocityFirst  sprec.Vec3
	SlopeVelocitySecond        sprec.Vec3
	SlopeAngularVelocitySecond sprec.Vec3
}

func (j DoubleEntityJacobian) Apply(firstEntity, secondEntity *Entity) {
	firstMotionComp := firstEntity.Motion
	secondMotionComp := secondEntity.Motion

	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocityFirst, firstMotionComp.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstMotionComp.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondMotionComp.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondMotionComp.AngularVelocity))
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleEntityJacobian) ApplySoft(firstEntity, secondEntity *Entity, force, gamma float32) {
	firstMotionComp := firstEntity.Motion
	secondMotionComp := secondEntity.Motion

	lambdaUpper := -(sprec.Vec3Dot(j.SlopeVelocityFirst, firstMotionComp.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocityFirst, firstMotionComp.AngularVelocity) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, secondMotionComp.Velocity) +
		sprec.Vec3Dot(j.SlopeAngularVelocitySecond, secondMotionComp.AngularVelocity) -
		force*gamma)
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleEntityJacobian) ApplyNew(firstEntity, secondEntity *Entity, velocityDrift float32) {
	firstMotionComp := firstEntity.Motion
	secondMotionComp := secondEntity.Motion

	lambdaUpper := -velocityDrift
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleEntityJacobian) ApplyBaumgarte(firstEntity, secondEntity *Entity, drift float32) {
	firstMotionComp := firstEntity.Motion
	secondMotionComp := secondEntity.Motion

	lambdaUpper := -driftBaumgarteCorrectionAmount * drift / timeStep
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	firstMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	firstMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	secondMotionComp.ApplyImpulse(sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	secondMotionComp.ApplyAngularImpulse(sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func (j DoubleEntityJacobian) ApplyNudge(firstEntity, secondEntity *Entity, drift float32) {
	firstMotionComp := firstEntity.Motion
	secondMotionComp := secondEntity.Motion

	lambdaUpper := -driftCorrectionAmount * drift
	lambdaLower := sprec.Vec3Dot(j.SlopeVelocityFirst, j.SlopeVelocityFirst)/firstMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), j.SlopeAngularVelocityFirst), j.SlopeAngularVelocityFirst) +
		sprec.Vec3Dot(j.SlopeVelocitySecond, j.SlopeVelocitySecond)/secondMotionComp.Mass +
		sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), j.SlopeAngularVelocitySecond), j.SlopeAngularVelocitySecond)
	lambda := lambdaUpper / lambdaLower

	applyNudge(firstEntity, sprec.Vec3Prod(j.SlopeVelocityFirst, lambda))
	applyAngularNudge(firstEntity, sprec.Vec3Prod(j.SlopeAngularVelocityFirst, lambda))
	applyNudge(secondEntity, sprec.Vec3Prod(j.SlopeVelocitySecond, lambda))
	applyAngularNudge(secondEntity, sprec.Vec3Prod(j.SlopeAngularVelocitySecond, lambda))
}

func applyNudge(entity *Entity, nudge sprec.Vec3) {
	transformComp := entity.Transform
	motionComp := entity.Motion
	transformComp.Translate(sprec.Vec3Quot(nudge, motionComp.Mass))
}

func applyAngularNudge(entity *Entity, angularNudge sprec.Vec3) {
	transformComp := entity.Transform
	motionComp := entity.Motion
	transformComp.Rotate(sprec.Mat3Vec3Prod(sprec.InverseMat3(motionComp.MomentOfInertia), angularNudge))
}

type Constraint interface {
	ApplyForces()
	ApplyCorrectionForces()
	ApplyCorrectionImpulses()
	ApplyCorrectionBaumgarte()
	ApplyCorrectionTranslations() // FIXME: Rename to ApplyCorrectionTransforms
}

type DebuggableConstraint interface {
	Error() float32
}

type RenderableConstraint interface {
	Lines() []DebugLine
}

var _ Constraint = NilConstraint{}

type NilConstraint struct{}

func (NilConstraint) ApplyForces() {}

func (NilConstraint) ApplyCorrectionForces() {}

func (NilConstraint) ApplyCorrectionImpulses() {}

func (NilConstraint) ApplyCorrectionBaumgarte() {}

func (NilConstraint) ApplyCorrectionTranslations() {}
