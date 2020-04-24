package ecs

import (
	"fmt"
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

func NewPhysicsSystem(ecsManager *Manager, step time.Duration) *PhysicsSystem {
	return &PhysicsSystem{
		ecsManager: ecsManager,

		step:            step,
		accumulatedTime: 0,

		gravity: sprec.NewVec3(0.0, -9.8, 0.0),
		// gravity:      sprec.NewVec3(0.0, -1.8, 0.0),
		windVelocity: sprec.NewVec3(0.0, 0.0, 0.0),
		// windVelocity: sprec.NewVec3(0.0, 0.0, 1.0),
		windDensity: 1.2,
		// windDensity: 0.0,
		lines: make([]DebugLine, 0, 1024),
	}
}

type PhysicsSystem struct {
	ecsManager *Manager

	step            time.Duration
	accumulatedTime time.Duration

	constraints          []Constraint
	collisionConstraints []Constraint

	gravity      sprec.Vec3
	windVelocity sprec.Vec3
	windDensity  float32

	debugSkips int

	lines []DebugLine
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

type DebugLine struct {
	A     sprec.Vec3
	B     sprec.Vec3
	Color sprec.Vec4
}

func (s *PhysicsSystem) GetDebug() []DebugLine {
	return s.lines
}

func (s *PhysicsSystem) runSimulation(elapsedSeconds float32) {
	s.lines = s.lines[:0]
	s.applyForces()
	s.applyCorrectionForces()

	s.integrate(elapsedSeconds)
	s.applyCorrectionImpulses()
	s.applyCorrectionBaumgarte()
	s.applyMotion(elapsedSeconds)
	s.applyCorrectionTranslations()
	s.detectCollisions()

	s.renderDebug()
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

		// TODO: Where to get the radius and length (maybe a magnus tensor)?
		radius := float32(0.3)
		length := float32(0.4)
		motionComp.ApplyForce(sprec.Vec3Prod(sprec.Vec3Cross(deltaWindVelocity, sprec.Vec3Prod(motionComp.AngularVelocity, 2*sprec.Pi*radius*radius)), s.windDensity*length))
	}

	for _, constraint := range s.constraints {
		constraint.ApplyForces()
	}
	for _, constraint := range s.collisionConstraints {
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
		for _, constraint := range s.collisionConstraints {
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
	const accuracy = 100
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionImpulses()
		}
		for _, constraint := range s.collisionConstraints {
			constraint.ApplyCorrectionImpulses()
		}
	}
}

func (s *PhysicsSystem) applyCorrectionBaumgarte() {
	const accuracy = 1
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionBaumgarte()
		}
		for _, constraint := range s.collisionConstraints {
			constraint.ApplyCorrectionBaumgarte()
		}
	}
}

func (s *PhysicsSystem) applyCorrectionTranslations() {
	const accuracy = 100
	for i := 0; i < accuracy; i++ {
		for _, constraint := range s.constraints {
			constraint.ApplyCorrectionTranslations()
		}
		for _, constraint := range s.collisionConstraints {
			constraint.ApplyCorrectionTranslations()
		}
	}
}

func (s *PhysicsSystem) detectCollisions() {
	s.collisionConstraints = s.collisionConstraints[:0]
	for _, entity := range s.ecsManager.Entities() {
		collisionComp := entity.Collision
		if collisionComp == nil {
			continue
		}

		// FIXME: This does collision checks twice!
		for _, otherEntity := range s.ecsManager.Entities() {
			otherCollisionComp := otherEntity.Collision
			if otherCollisionComp == nil {
				continue
			}
			s.checkEntitiesCollision(entity, otherEntity)
		}
	}
}

func (s *PhysicsSystem) checkEntitiesCollision(firstEntity, secondEntity *Entity) {
	firstTransformComp := firstEntity.Transform
	firstMotionComp := firstEntity.Motion
	firstCollisionComp := firstEntity.Collision

	if firstMotionComp == nil {
		return
	}

	secondCollisionComp := secondEntity.Collision
	if _, ok := secondCollisionComp.CollisionShape.(MeshShape); !ok {
		return
	}
	secondCollisionShape := secondEntity.Collision.CollisionShape.(MeshShape)

	checkLineCollision := func(a, b sprec.Vec3) {
		if result, ok := secondCollisionShape.Mesh.LineCollision(collision.MakeLine(a, b)); ok {
			s.collisionConstraints = append(s.collisionConstraints, GroundCollisionConstraint{
				Entity:           firstEntity,
				OriginalPosition: firstTransformComp.Position,
				Normal:           result.Normal(),
				ContactPoint:     result.Intersection(),
				Depth:            sprec.Abs(result.BottomHeight()), // FIXME: Shouldn't it just be negative or positive
			})
		}
		if result, ok := secondCollisionShape.Mesh.LineCollision(collision.MakeLine(b, a)); ok {
			s.collisionConstraints = append(s.collisionConstraints, GroundCollisionConstraint{
				Entity:           firstEntity,
				OriginalPosition: firstTransformComp.Position,
				Normal:           result.Normal(),
				ContactPoint:     result.Intersection(),
				Depth:            sprec.Abs(result.BottomHeight()), // FIXME: Shouldn't it just be negative or positive
			})
		}
	}

	switch firstCollisionShape := firstCollisionComp.CollisionShape.(type) {
	case BoxShape:
		minX := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationX(), firstCollisionShape.MinX)
		maxX := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationX(), firstCollisionShape.MaxX)
		minY := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationY(), firstCollisionShape.MinY)
		maxY := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationY(), firstCollisionShape.MaxY)
		minZ := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationZ(), firstCollisionShape.MinZ)
		maxZ := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationZ(), firstCollisionShape.MaxZ)

		p1 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, minX), minZ), maxY)
		p2 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, minX), maxZ), maxY)
		p3 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, maxX), maxZ), maxY)
		p4 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, maxX), minZ), maxY)
		p5 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, minX), minZ), minY)
		p6 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, minX), maxZ), minY)
		p7 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, maxX), maxZ), minY)
		p8 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(firstTransformComp.Position, maxX), minZ), minY)

		checkLineCollision(p1, p2)
		checkLineCollision(p2, p3)
		checkLineCollision(p3, p4)
		checkLineCollision(p4, p1)

		checkLineCollision(p5, p6)
		checkLineCollision(p6, p7)
		checkLineCollision(p7, p8)
		checkLineCollision(p8, p5)

		checkLineCollision(p1, p5)
		checkLineCollision(p2, p6)
		checkLineCollision(p3, p7)
		checkLineCollision(p4, p8)

		// checkLineCollision(p2, p1)
		// checkLineCollision(p3, p2)
		// checkLineCollision(p4, p3)
		// checkLineCollision(p1, p4)
		// checkLineCollision(p6, p5)
		// checkLineCollision(p7, p6)
		// checkLineCollision(p8, p7)
		// checkLineCollision(p5, p8)

		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfHeight))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfWidth))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfWidth))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfLength))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfLength))

		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfHeight))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfWidth))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfWidth))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfLength))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfLength))

		// s.addCollisionLine(p1, p2)
		// s.addCollisionLine(p2, p3)
		// s.addCollisionLine(p3, p4)
		// s.addCollisionLine(p4, p1)

		// s.addCollisionLine(p5, p6)
		// s.addCollisionLine(p6, p7)
		// s.addCollisionLine(p7, p8)
		// s.addCollisionLine(p8, p5)

		// s.addCollisionLine(p1, p5)
		// s.addCollisionLine(p2, p6)
		// s.addCollisionLine(p3, p7)
		// s.addCollisionLine(p4, p8)

	case CylinderShape:
		halfWidth := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationX(), firstCollisionShape.Length)
		halfHeight := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationY(), firstCollisionShape.Radius)
		halfLength := sprec.Vec3Prod(firstTransformComp.Orientation.OrientationZ(), firstCollisionShape.Radius)

		checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfWidth))
		checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfWidth))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfHeight))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfHeight))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, halfLength))
		// checkLineCollision(firstTransformComp.Position, sprec.Vec3Diff(firstTransformComp.Position, halfLength))

		const precision = 48
		for i := 0; i < precision; i++ {
			cos := sprec.Cos(sprec.Degrees(360.0 * float32(i) / 12.0))
			sin := sprec.Sin(sprec.Degrees(360.0 * float32(i) / 12.0))
			direction := sprec.Vec3Sum(sprec.Vec3Prod(halfLength, cos), sprec.Vec3Prod(halfHeight, sin))
			checkLineCollision(firstTransformComp.Position, sprec.Vec3Sum(firstTransformComp.Position, direction))
		}
	}
}

func (s *PhysicsSystem) addCollisionLine(a, b sprec.Vec3) {
	s.lines = append(s.lines, DebugLine{
		A:     a,
		B:     b,
		Color: sprec.NewVec4(1.0, 1.0, 1.0, 1.0),
	})
}

func (s *PhysicsSystem) renderDebug() {
	for _, constraint := range s.constraints {
		if renderable, ok := constraint.(RenderableConstraint); ok {
			s.lines = append(s.lines, renderable.Lines()...)
		}
	}
}

func (s *PhysicsSystem) printDebug() {
	if s.debugSkips++; s.debugSkips != 10 {
		return
	}
	s.debugSkips = 0

	for _, entity := range s.ecsManager.Entities() {
		transformComp := entity.Transform
		motionComp := entity.Motion
		debugComp := entity.Debug
		if transformComp == nil || motionComp == nil || debugComp == nil {
			continue
		}
		fmt.Printf("Entity [%s]:\n", debugComp.Name)
		fmt.Printf("- position: %#v\n", transformComp.Position)
		fmt.Printf("- orientation: %#v\n", transformComp.Orientation)
		fmt.Printf("- velocity: %#v\n", motionComp.Velocity)
		fmt.Printf("- angular velocity: %#v\n", motionComp.AngularVelocity)
	}
	for _, constraint := range s.constraints {
		if debuggable, ok := constraint.(DebuggableConstraint); ok {
			fmt.Printf("constraint error: %f\n", debuggable.Error())
		}
	}
}
