package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
	"github.com/mokiat/lacking/render"
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
	Renderable *render.Renderable
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
	RotationConstraint   *solver.MatchAxis
	AccelerationVelocity float32
	DecelerationVelocity float32
}

type CameraStand struct {
	Target         *Entity
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *render.Camera
}
