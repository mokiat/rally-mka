package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

type CollisionComponent struct {
	RestitutionCoef float32
	CollisionShape  interface{}
	Group           int
}

type BoxShape struct {
	MinX float32
	MaxX float32
	MinY float32
	MaxY float32
	MinZ float32
	MaxZ float32
}

type CylinderShape struct {
	Length float32
	Radius float32
}

type MeshShape struct {
	Mesh *collision.Mesh
}

type GroundCollisionConstraint struct {
	NilConstraint
	Entity           *Entity
	OriginalPosition sprec.Vec3
	Normal           sprec.Vec3
	ContactPoint     sprec.Vec3
	Depth            float32
}

func (c GroundCollisionConstraint) ApplyForces() {
	transformComp := c.Entity.Transform
	motionComp := c.Entity.Motion
	relativeContactPosition := sprec.Vec3Diff(c.ContactPoint, transformComp.Position)
	contactVelocity := sprec.Vec3Sum(motionComp.Velocity, sprec.Vec3Cross(motionComp.AngularVelocity, relativeContactPosition))

	lateralVelocity := sprec.Vec3Diff(contactVelocity, sprec.Vec3Prod(c.Normal, sprec.Vec3Dot(contactVelocity, c.Normal)))
	maxFriction := float32(100.0)
	if lateralVelocity.Length() > maxFriction {
		lateralVelocity = sprec.ResizedVec3(lateralVelocity, maxFriction)
	}

	c.Entity.Motion.ApplyOffsetImpulse(relativeContactPosition, sprec.Vec3Prod(lateralVelocity, -motionComp.Mass/10))
}

// func (c GroundCollisionConstraint) ApplyCorrectionForces() {
// 	transformComp := c.Entity.Transform
// 	motionComp := c.Entity.Motion

// 	relativeContactPosition := sprec.Vec3Diff(c.ContactPoint, transformComp.Position)
// 	contactAcceleration := sprec.Vec3Sum(motionComp.Acceleration, sprec.Vec3Cross(motionComp.AngularAcceleration, relativeContactPosition))

// 	normalAcceleration := sprec.Vec3Dot(c.Normal, contactAcceleration)
// 	if normalAcceleration > 0 {
// 		return // moving away from ground
// 	}

// 	totalMass := 1.0 / ((1.0 / motionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(motionComp.MomentOfInertia), sprec.Vec3Cross(relativeContactPosition, c.Normal)), sprec.Vec3Cross(relativeContactPosition, c.Normal)))
// 	forceStrength := totalMass * sprec.Vec3Dot(c.Normal, contactAcceleration)
// 	motionComp.ApplyOffsetForce(relativeContactPosition, sprec.Vec3Prod(c.Normal, -forceStrength))
// }

func (c GroundCollisionConstraint) ApplyCorrectionImpulses() {
	transformComp := c.Entity.Transform
	motionComp := c.Entity.Motion
	collisionComp := c.Entity.Collision

	relativeContactPosition := sprec.Vec3Diff(c.ContactPoint, transformComp.Position)
	contactVelocity := sprec.Vec3Sum(motionComp.Velocity, sprec.Vec3Cross(motionComp.AngularVelocity, relativeContactPosition))

	jacobian := sprec.InverseVec3(c.Normal)
	normalVelocity := sprec.Vec3Dot(c.Normal, contactVelocity)
	if normalVelocity > 0 {
		return // moving away from ground
	}

	// restitutionClamp := float32(1.0)
	// if sprec.Abs(normalVelocity) < 2 {
	// 	restitutionClamp = 0.1
	// }
	// if sprec.Abs(normalVelocity) < 1 {
	// 	restitutionClamp = 0.05
	// }
	// if sprec.Abs(normalVelocity) < 0.5 {
	// 	restitutionClamp = 0.0
	// }

	restitutionClamp := float32(0.0) // TODO: Delete, use above one

	totalMass := (1 + collisionComp.RestitutionCoef*restitutionClamp) / ((1.0 / motionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(motionComp.MomentOfInertia), sprec.Vec3Cross(relativeContactPosition, jacobian)), sprec.Vec3Cross(relativeContactPosition, jacobian)))
	impulseStrength := totalMass*sprec.Vec3Dot(jacobian, contactVelocity) + totalMass*c.Depth // FIXME
	motionComp.ApplyOffsetImpulse(relativeContactPosition, sprec.InverseVec3(sprec.Vec3Prod(jacobian, impulseStrength)))

	// tangentRadius := radius.DecVec3(c.Normal.Mul(math.Vec3DotProduct(radius, c.Normal)))
	// impulse := -normalVelocity * (1 + c.Body.RestitutionCoef*restitutionClamp) / ((1 / c.Body.Mass) + (tangentRadius.LengthSquared() / c.Body.MomentOfInertia))
	// c.Body.ApplyOffsetImpulse(c.Normal.Mul(impulse), radius)

	// 	firstTransformComp := r.First.Transform
	// 	secondTransformComp := r.Second.Transform

	// 	firstAnchorRelativePosition := firstTransformComp.Orientation.MulVec3(r.FirstAnchor)
	// 	secondAnchorRelativePosition := secondTransformComp.Orientation.MulVec3(r.SecondAnchor)

	// 	firstAnchorPosition := sprec.Vec3Sum(firstTransformComp.Position, firstAnchorRelativePosition)
	// 	secondAnchorPosition := sprec.Vec3Sum(secondTransformComp.Position, secondAnchorRelativePosition)

	// 	deltaPosition := sprec.Vec3Diff(secondAnchorPosition, firstAnchorPosition)
	// 	jacobian := sprec.UnitVec3(deltaPosition) // FIXME: Handle if deltaPosition == 0

	// 	firstMotionComp := r.First.Motion
	// 	secondMotionComp := r.Second.Motion

	// 	firstPointVelocity := sprec.Vec3Sum(firstMotionComp.Velocity, sprec.Vec3Cross(firstMotionComp.AngularVelocity, firstAnchorRelativePosition))
	// 	secondPointVelocity := sprec.Vec3Sum(secondMotionComp.Velocity, sprec.Vec3Cross(secondMotionComp.AngularVelocity, secondAnchorRelativePosition))
	// 	deltaVelocity := sprec.Vec3Diff(secondPointVelocity, firstPointVelocity)

	// 	firstEffectiveMass := 1.0 / ((1.0 / firstMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)))
	// 	secondEffectiveMass := 1.0 / ((1.0 / secondMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)))
	// 	totalMass := firstEffectiveMass * secondEffectiveMass / (firstEffectiveMass + secondEffectiveMass)

	// 	impulseStrength := totalMass * sprec.Vec3Dot(jacobian, deltaVelocity)
	// 	impulse := sprec.Vec3Prod(jacobian, impulseStrength)
	// 	firstMotionComp.ApplyOffsetImpulse(firstAnchorRelativePosition, impulse)
	// 	secondMotionComp.ApplyOffsetImpulse(secondAnchorRelativePosition, sprec.InverseVec3(impulse))
}

func (c GroundCollisionConstraint) ApplyCorrectionTranslations() {
	// 	firstTransformComp := r.First.Transform
	// 	secondTransformComp := r.Second.Transform

	// 	firstAnchorRelativePosition := firstTransformComp.Orientation.MulVec3(r.FirstAnchor)
	// 	secondAnchorRelativePosition := secondTransformComp.Orientation.MulVec3(r.SecondAnchor)

	// 	firstAnchorPosition := sprec.Vec3Sum(firstTransformComp.Position, firstAnchorRelativePosition)
	// 	secondAnchorPosition := sprec.Vec3Sum(secondTransformComp.Position, secondAnchorRelativePosition)

	// 	deltaPosition := sprec.Vec3Diff(secondAnchorPosition, firstAnchorPosition)
	// 	jacobian := sprec.UnitVec3(deltaPosition) // FIXME: Handle if deltaPosition == 0

	// 	firstMotionComp := r.First.Motion
	// 	secondMotionComp := r.Second.Motion

	// 	firstEffectiveMass := 1.0 / ((1.0 / firstMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)))
	// 	secondEffectiveMass := 1.0 / ((1.0 / secondMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)))
	// 	totalMass := firstEffectiveMass * secondEffectiveMass / (firstEffectiveMass + secondEffectiveMass)

	// 	nudgeStrength := totalMass * sprec.Vec3Dot(jacobian, sprec.ResizedVec3(deltaPosition, deltaPosition.Length()-r.Length))
	// 	nudge := sprec.Vec3Prod(jacobian, nudgeStrength)

	// 	firstTransformComp.Translate(sprec.Vec3Quot(nudge, firstMotionComp.Mass))
	// 	secondTransformComp.Translate(sprec.InverseVec3(sprec.Vec3Quot(nudge, secondMotionComp.Mass)))

	// 	firstTransformComp.Rotate(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, nudge)))
	// 	secondTransformComp.Rotate(sprec.InverseVec3(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, nudge))))
}

// type BodyToGroundCollision struct {
// 	Body             *Body
// 	OriginalPosition math.Vec3
// 	ContactPoint     math.Vec3
// 	Normal           math.Vec3
// 	Depth            float32
// }

// func (c BodyToGroundCollision) ApplyImpulse() {
// 	radius := c.ContactPoint.DecVec3(c.Body.Position)
// 	tangentRadius := radius.DecVec3(c.Normal.Mul(math.Vec3DotProduct(radius, c.Normal)))
// 	pointVelocity := c.Body.Velocity.IncVec3(math.Vec3CrossProduct(c.Body.AngularVelocity, radius))
// 	normalVelocity := math.Vec3DotProduct(c.Normal, pointVelocity)
// 	if c.Body.Name == "front-left-wheel" {
// 		// fmt.Printf("normal velocity: %f\n", normalVelocity)
// 	}
// 	if normalVelocity > 0 {
// 		return // moving away from ground
// 	}
// 	restitutionClamp := float32(1.0)
// 	if math.Abs32(normalVelocity) < 2 {
// 		restitutionClamp = 0.1
// 	}
// 	if math.Abs32(normalVelocity) < 1 {
// 		restitutionClamp = 0.05
// 	}
// 	if math.Abs32(normalVelocity) < 0.5 {
// 		restitutionClamp = 0.0
// 	}
// 	impulse := -normalVelocity * (1 + c.Body.RestitutionCoef*restitutionClamp) / ((1 / c.Body.Mass) + (tangentRadius.LengthSquared() / c.Body.MomentOfInertia))
// 	c.Body.ApplyOffsetImpulse(c.Normal.Mul(impulse), radius)
// }

// func (c BodyToGroundCollision) ApplySeparation() {
// 	deltaPosition := c.Body.Position.DecVec3(c.OriginalPosition)
// 	deltaDepth := math.Vec3DotProduct(deltaPosition, c.Normal)
// 	actualDepth := c.Depth - deltaDepth
// 	if actualDepth > 0 {
// 		c.Body.Position = c.Body.Position.IncVec3(c.Normal.Mul(actualDepth * 1.0))
// 	}
// }
