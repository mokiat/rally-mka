package physics

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/engine/physics/collision"
)

const (
	gravity     = 9.8
	windDensity = 1.2

	impulseIterations = 100
	nudgeIterations   = 100
)

func NewEngine(step time.Duration) *Engine {
	return &Engine{
		step:            step,
		accumulatedTime: 0,

		gravity:      sprec.NewVec3(0.0, -gravity, 0.0),
		windVelocity: sprec.NewVec3(0.0, 0.0, 0.0),
		windDensity:  windDensity,
	}
}

type Engine struct {
	step            time.Duration
	accumulatedTime time.Duration

	gravity      sprec.Vec3
	windVelocity sprec.Vec3
	windDensity  float32

	bodies               []*Body
	constraints          []Constraint
	collisionConstraints []Constraint
}

func (e *Engine) Bodies() []*Body {
	return e.bodies
}

func (e *Engine) Update(elapsedTime time.Duration) {
	e.accumulatedTime += elapsedTime
	for e.accumulatedTime > e.step {
		e.accumulatedTime -= e.step
		e.runSimulation(float32(e.step.Seconds()))
	}
}

func (e *Engine) Add(aspect interface{}) {
	if body, ok := aspect.(*Body); ok {
		e.AddBody(body)
	}
	if constraint, ok := aspect.(Constraint); ok {
		e.AddConstraint(constraint)
	}
}

func (e *Engine) AddBody(body *Body) {
	e.bodies = append(e.bodies, body)
}

func (e *Engine) AddConstraint(constraint Constraint) {
	e.constraints = append(e.constraints, constraint)
}

func (e *Engine) runSimulation(elapsedSeconds float32) {
	e.applyForces()
	e.integrate(elapsedSeconds)
	e.applyImpulses()
	e.applyBaumgartes()
	e.applyMotion(elapsedSeconds)
	// TODO: Should the following two be swapped
	e.applyNudges()
	e.detectCollisions()
}

func (e *Engine) applyForces() {
	for _, body := range e.bodies {
		if body.IsStatic {
			continue
		}
		body.ResetAcceleration()
		body.ResetAngularAcceleration()

		body.AddAcceleration(e.gravity)
		deltaWindVelocity := sprec.Vec3Diff(e.windVelocity, body.Velocity)
		body.ApplyForce(sprec.Vec3Prod(deltaWindVelocity, e.windDensity*body.DragFactor*deltaWindVelocity.Length()))
		body.ApplyTorque(sprec.Vec3Prod(body.AngularVelocity, -e.windDensity*body.AngularDragFactor*body.AngularVelocity.Length()))

		// TODO: Where to get the radius and length (maybe a magnus tensor)?
		// radius := float32(0.3)
		// length := float32(0.4)
		// body.ApplyForce(sprec.Vec3Prod(sprec.Vec3Cross(deltaWindVelocity, sprec.Vec3Prod(body.AngularVelocity, 2*sprec.Pi*radius*radius)), e.windDensity*length))
	}

	for _, constraint := range e.constraints {
		constraint.ApplyForce()
	}
	for _, constraint := range e.collisionConstraints {
		constraint.ApplyForce()
	}

	// TODO: Restrict max linear + angular accelerations
}

func (e *Engine) integrate(elapsedSeconds float32) {
	// we use semi-implicit euler as it is simple and
	// stable with harmonic motion (like springs)

	for _, body := range e.bodies {
		if body.IsStatic {
			continue
		}
		deltaVelocity := sprec.Vec3Prod(body.Acceleration, elapsedSeconds)
		body.AddVelocity(deltaVelocity)
		deltaAngularVelocity := sprec.Vec3Prod(body.AngularAcceleration, elapsedSeconds)
		body.AddAngularVelocity(deltaAngularVelocity)

		// TODO: Restrict max linear + angular velocities
	}
}

func (e *Engine) applyImpulses() {
	for i := 0; i < impulseIterations; i++ {
		for _, constraint := range e.constraints {
			constraint.ApplyImpulse()
		}
		for _, constraint := range e.collisionConstraints {
			constraint.ApplyImpulse()
		}
	}
}

func (e *Engine) applyBaumgartes() {
	for _, constraint := range e.constraints {
		constraint.ApplyBaumgarte()
	}
	for _, constraint := range e.collisionConstraints {
		constraint.ApplyBaumgarte()
	}
}

func (e *Engine) applyMotion(elapsedSeconds float32) {
	for _, body := range e.bodies {
		deltaPosition := sprec.Vec3Prod(body.Velocity, elapsedSeconds)
		body.Translate(deltaPosition)
		deltaRotation := sprec.Vec3Prod(body.AngularVelocity, elapsedSeconds)
		body.Rotate(deltaRotation)
	}
}

func (e *Engine) applyNudges() {
	for i := 0; i < nudgeIterations; i++ {
		for _, constraint := range e.constraints {
			constraint.ApplyNudge()
		}
		for _, constraint := range e.collisionConstraints {
			constraint.ApplyNudge()
		}
	}
}

func (e *Engine) detectCollisions() {
	e.collisionConstraints = e.collisionConstraints[:0]
	for i := 0; i < len(e.bodies); i++ {
		for j := i + 1; j < len(e.bodies); j++ {
			first := e.bodies[i]
			second := e.bodies[j]
			e.checkCollisionTwoBodies(first, second)
		}
	}
}

func (e *Engine) checkCollisionTwoBodies(first, second *Body) {
	// FIXME: Temp hack section
	if !first.IsStatic && !second.IsStatic {
		return
	}
	if first.IsStatic {
		first, second = second, first
	}

	secondCollisionShape := second.CollisionShape.(MeshShape)

	checkLineCollision := func(a, b sprec.Vec3) {
		if result, ok := secondCollisionShape.Mesh.LineCollision(collision.MakeLine(a, b)); ok {
			e.collisionConstraints = append(e.collisionConstraints, GroundCollisionConstraint{
				Body:             first,
				OriginalPosition: first.Position,
				Normal:           result.Normal(),
				ContactPoint:     result.Intersection(),
				Depth:            sprec.Abs(result.BottomHeight()), // FIXME: Shouldn't it just be negative or positive
			})
		}
		if result, ok := secondCollisionShape.Mesh.LineCollision(collision.MakeLine(b, a)); ok {
			e.collisionConstraints = append(e.collisionConstraints, GroundCollisionConstraint{
				Body:             first,
				OriginalPosition: first.Position,
				Normal:           result.Normal(),
				ContactPoint:     result.Intersection(),
				Depth:            sprec.Abs(result.BottomHeight()), // FIXME: Shouldn't it just be negative or positive
			})
		}
	}

	switch firstCollisionShape := first.CollisionShape.(type) {
	case BoxShape:
		minX := sprec.Vec3Prod(first.Orientation.OrientationX(), firstCollisionShape.MinX)
		maxX := sprec.Vec3Prod(first.Orientation.OrientationX(), firstCollisionShape.MaxX)
		minY := sprec.Vec3Prod(first.Orientation.OrientationY(), firstCollisionShape.MinY)
		maxY := sprec.Vec3Prod(first.Orientation.OrientationY(), firstCollisionShape.MaxY)
		minZ := sprec.Vec3Prod(first.Orientation.OrientationZ(), firstCollisionShape.MinZ)
		maxZ := sprec.Vec3Prod(first.Orientation.OrientationZ(), firstCollisionShape.MaxZ)

		p1 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, minX), minZ), maxY)
		p2 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, minX), maxZ), maxY)
		p3 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, maxX), maxZ), maxY)
		p4 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, maxX), minZ), maxY)
		p5 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, minX), minZ), minY)
		p6 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, minX), maxZ), minY)
		p7 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, maxX), maxZ), minY)
		p8 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(first.Position, maxX), minZ), minY)

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

	case CylinderShape:
		halfWidth := sprec.Vec3Prod(first.Orientation.OrientationX(), firstCollisionShape.Length)
		halfHeight := sprec.Vec3Prod(first.Orientation.OrientationY(), firstCollisionShape.Radius)
		halfLength := sprec.Vec3Prod(first.Orientation.OrientationZ(), firstCollisionShape.Radius)

		checkLineCollision(first.Position, sprec.Vec3Sum(first.Position, halfWidth))
		checkLineCollision(first.Position, sprec.Vec3Diff(first.Position, halfWidth))

		const precision = 48
		for i := 0; i < precision; i++ {
			cos := sprec.Cos(sprec.Degrees(360.0 * float32(i) / 12.0))
			sin := sprec.Sin(sprec.Degrees(360.0 * float32(i) / 12.0))
			direction := sprec.Vec3Sum(sprec.Vec3Prod(halfLength, cos), sprec.Vec3Prod(halfHeight, sin))
			checkLineCollision(first.Position, sprec.Vec3Sum(first.Position, direction))
		}

	case SphereShape:
		var bestCollision collision.LineCollision
		var found bool
		for _, triangle := range secondCollisionShape.Mesh.Triangles() {
			deltaPosition := sprec.Vec3Diff(first.Position, triangle.Center())
			if deltaPosition.Length() > firstCollisionShape.Radius+triangle.BoudingSphereRadius() {
				continue
			}

			distance := sprec.Vec3Dot(triangle.Normal(), deltaPosition)
			if distance > firstCollisionShape.Radius {
				continue
			}
			projectedPoint := sprec.Vec3Diff(first.Position, sprec.Vec3Prod(triangle.Normal(), distance))
			if triangle.Contains(projectedPoint) {
				bestCollision = collision.NewLineCollision(
					projectedPoint,
					triangle.Normal(),
					distance,
					firstCollisionShape.Radius-distance,
				)
				found = true
			}
		}
		if found {
			e.collisionConstraints = append(e.collisionConstraints, GroundCollisionConstraint{
				Body:             first,
				OriginalPosition: first.Position,
				Normal:           bestCollision.Normal(),
				ContactPoint:     bestCollision.Intersection(),
				Depth:            sprec.Abs(bestCollision.BottomHeight()), // FIXME: Shouldn't it just be negative or positive
			})
		}
	}
}
