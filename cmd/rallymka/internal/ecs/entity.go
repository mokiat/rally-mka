package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/physics"
	"github.com/mokiat/lacking/world"
)

type Entity struct {
	Physics       *PhysicsComponent
	Render        *RenderComponent
	Vehicle       *Vehicle
	CameraStand   *CameraStand
	PlayerControl *PlayerControl
}

type PhysicsComponent struct {
	Body *physics.Body
}

type RenderComponent struct {
	Renderable *world.Renderable
}

type PlayerControl struct {
}

type Vehicle struct {
	MaxSteeringAngle sprec.Angle
	SteeringAngle    sprec.Angle
	Acceleration     float32
	Deceleration     float32
	Recover          bool

	Chassis *Chassis
	Wheels  []*Wheel
}

type Chassis struct {
	Body *physics.Body
}

type Wheel struct {
	Body                 *physics.Body
	RotationConstraint   *physics.MatchAxisConstraint
	AccelerationVelocity float32
	DecelerationVelocity float32
}

type CameraStand struct {
	Target         *Entity
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *world.Camera
}
