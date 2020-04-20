package ecs

import "github.com/mokiat/gomath/sprec"

type MotionComponent struct {
	Mass            float32
	MomentOfInertia sprec.Mat3

	Acceleration        sprec.Vec3
	AngularAcceleration sprec.Vec3

	Velocity        sprec.Vec3
	AngularVelocity sprec.Vec3

	DragFactor        float32
	AngularDragFactor float32
}

func (c *MotionComponent) ResetAcceleration() {
	c.Acceleration = sprec.ZeroVec3()
}

func (c *MotionComponent) ResetAngularAcceleration() {
	c.AngularAcceleration = sprec.ZeroVec3()
}

func (c *MotionComponent) AddAcceleration(amount sprec.Vec3) {
	c.Acceleration = sprec.Vec3Sum(c.Acceleration, amount)
}

func (c *MotionComponent) AddAngularAcceleration(amount sprec.Vec3) {
	c.AngularAcceleration = sprec.Vec3Sum(c.AngularAcceleration, amount)
}

func (c *MotionComponent) AddAngularVelocity(amount sprec.Vec3) {
	c.AngularVelocity = sprec.Vec3Sum(c.AngularVelocity, amount)
}

func (c *MotionComponent) AddVelocity(amount sprec.Vec3) {
	c.Velocity = sprec.Vec3Sum(c.Velocity, amount)
}

func (c *MotionComponent) ApplyForce(force sprec.Vec3) {
	c.AddAcceleration(sprec.Vec3Quot(force, c.Mass))
}

func (c *MotionComponent) ApplyTorque(torque sprec.Vec3) {
	// FIXME: the moment of intertia is in local space, whereas the torque is in world space
	c.AddAngularAcceleration(sprec.Mat3Vec3Prod(sprec.InverseMat3(c.MomentOfInertia), torque))
}

func (c *MotionComponent) ApplyOffsetForce(offset, force sprec.Vec3) {
	c.ApplyForce(force)
	c.ApplyTorque(sprec.Vec3Cross(offset, force))
}

func (c *MotionComponent) ApplyImpulse(impulse sprec.Vec3) {
	c.Velocity = sprec.Vec3Sum(c.Velocity, sprec.Vec3Quot(impulse, c.Mass))
}

func (c *MotionComponent) ApplyAngularImpulse(impulse sprec.Vec3) {
	// FIXME: the moment of intertia is in local space, whereas the impulse is in world space
	c.AddAngularVelocity(sprec.Mat3Vec3Prod(sprec.InverseMat3(c.MomentOfInertia), impulse))
}

func (c *MotionComponent) ApplyOffsetImpulse(offset, impulse sprec.Vec3) {
	c.ApplyImpulse(impulse)
	c.ApplyAngularImpulse(sprec.Vec3Cross(offset, impulse))
}
