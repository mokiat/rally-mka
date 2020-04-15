package ecs

import (
	"fmt"
	"time"

	"github.com/mokiat/gomath/sprec"
)

func NewPhysicsSystem(ecsManager *Manager, step time.Duration) *PhysicsSystem {
	return &PhysicsSystem{
		ecsManager: ecsManager,

		step:            step,
		accumulatedTime: 0,

		// gravity: sprec.NewVec3(0.0, -9.8, 0.0),
		gravity:      sprec.NewVec3(0.0, -0.2, 0.0),
		windVelocity: sprec.NewVec3(0.0, 0.0, 0.0),
		// windVelocity: sprec.NewVec3(0.0, 0.0, 1.0),
		// windDensity: 1.2,
		windDensity: 0.0,
	}
}

type PhysicsSystem struct {
	ecsManager *Manager

	step            time.Duration
	accumulatedTime time.Duration

	constraints []Constraint

	gravity      sprec.Vec3
	windVelocity sprec.Vec3
	windDensity  float32
}

func (s *PhysicsSystem) AddConstraint(constraint Constraint) {
	s.constraints = append(s.constraints, constraint)
}

func (s *PhysicsSystem) Update(elapsedTime time.Duration) {
	s.accumulatedTime += elapsedTime
	for s.accumulatedTime > s.step {
		s.accumulatedTime -= s.step
		s.runSimulation(float32(s.step.Seconds()))
	}
}

func (s *PhysicsSystem) runSimulation(elapsedSeconds float32) {
	s.applyForces()
	s.applyCorrectionForces()

	s.integrate(elapsedSeconds)
	s.applyMotion(elapsedSeconds)
	// TODO: Collision detection

	s.applyCorrectionImpulses()
	s.applyCorrectionTranslations()
	s.printDebug()
}

func (s *PhysicsSystem) applyForces() {
	for _, entity := range s.ecsManager.Entities() {
		transformComp := entity.Transform
		motionComp := entity.Motion
		if transformComp == nil || motionComp == nil {
			continue
		}

		motionComp.ResetAcceleration()
		motionComp.ResetAngularAcceleration()

		motionComp.AddAcceleration(s.gravity)
		deltaWindVelocity := sprec.Vec3Diff(s.windVelocity, motionComp.Velocity)
		motionComp.ApplyForce(sprec.Vec3Prod(deltaWindVelocity, s.windDensity*motionComp.DragFactor*deltaWindVelocity.Length()))
		motionComp.ApplyTorque(sprec.Vec3Prod(motionComp.AngularVelocity, -s.windDensity*motionComp.AngularDragFactor*motionComp.AngularVelocity.Length()))

		radius := float32(0.3)
		length := float32(0.4)
		motionComp.ApplyForce(sprec.Vec3Prod(sprec.Vec3Cross(deltaWindVelocity, sprec.Vec3Prod(motionComp.AngularVelocity, 2*sprec.Pi*radius*radius)), s.windDensity*length)) // TODO: Where to get the radius and length (maybe a magnus tensor)?
		// TODO: Add magnus force ?
	}

	for _, constraint := range s.constraints {
		constraint.ApplyForces()
	}

	// TODO: Restrict max linear + angular accelerations
}

func (s *PhysicsSystem) applyCorrectionForces() {
	const accuracy = 1
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionForces()
		}
	}
}

func (s *PhysicsSystem) integrate(elapsedSeconds float32) {
	// we use semi-implicit euler as it is simple and
	// stable with harmonic motion (like springs)

	for _, entity := range s.ecsManager.Entities() {
		transformComp := entity.Transform
		motionComp := entity.Motion
		if transformComp == nil || motionComp == nil {
			continue
		}

		deltaVelocity := sprec.Vec3Prod(motionComp.Acceleration, elapsedSeconds)
		motionComp.AddVelocity(deltaVelocity)
		deltaAngularVelocity := sprec.Vec3Prod(motionComp.AngularAcceleration, elapsedSeconds)
		motionComp.AddAngularVelocity(deltaAngularVelocity)

		// TODO: Restrict max linear + angular velocities
	}
}

func (s *PhysicsSystem) applyMotion(elapsedSeconds float32) {
	for _, entity := range s.ecsManager.Entities() {
		transformComp := entity.Transform
		motionComp := entity.Motion
		if transformComp == nil || motionComp == nil {
			continue
		}

		deltaPosition := sprec.Vec3Prod(motionComp.Velocity, elapsedSeconds)
		transformComp.Translate(deltaPosition)
		deltaRotation := sprec.Vec3Prod(motionComp.AngularVelocity, elapsedSeconds)
		transformComp.Rotate(deltaRotation)
	}
}

func (s *PhysicsSystem) applyCorrectionImpulses() {
	const accuracy = 10
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionImpulses()
		}
	}
}

func (s *PhysicsSystem) applyCorrectionTranslations() {
	const accuracy = 10
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionTranslations()
		}
	}
}

var debugSkip int

func (s *PhysicsSystem) printDebug() {
	if debugSkip++; debugSkip%10 != 0 {
		return
	}
	for _, entity := range s.ecsManager.Entities() {
		transformComp := entity.Transform
		motionComp := entity.Motion
		debugComp := entity.Debug
		if transformComp == nil || motionComp == nil || debugComp == nil {
			continue
		}
		fmt.Printf("Entity [%s]:\n", debugComp.Name)
		fmt.Printf("- position: %#v\n", transformComp.Position)
		fmt.Printf("- velocity: %#v\n", motionComp.Velocity)
	}
}
