package ecs

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

const (
	gravity             = 0.01
	wheelFriction       = 99.0 / 100.0
	speedFriction       = 99.0 / 100.0
	speedFriction2      = 1.0 - (99.0 / 100.0)
	rotationFriction    = 60.0 / 100.0
	rotationFriction2   = 1.0 - 60.0/100.0
	maxSuspensionLength = 0.4
	skidWheelSpeed      = 0.03
)

type CarInput struct {
	Forward   bool
	Backward  bool
	TurnLeft  bool
	TurnRight bool
	Handbrake bool
}

type Stage interface {
	CheckCollision(line collision.Line) (collision.LineCollision, bool)
}

func NewVehicleSystem(ecsManager *Manager, stage Stage) *VehicleSystem {
	return &VehicleSystem{
		ecsManager: ecsManager,
		stage:      stage,
	}
}

type VehicleSystem struct {
	ecsManager *Manager
	stage      Stage
}

func (s *VehicleSystem) Update(elapsedTime time.Duration, input CarInput) {
	for _, entity := range s.ecsManager.Entities() {
		if vehicle := entity.Vehicle; vehicle != nil {
			if entity.Input != nil {
				s.updateVehicleInput(entity, elapsedTime, input)
			}
			s.updateVehicle(entity, elapsedTime)
		}
	}
}

func (s *VehicleSystem) updateVehicleInput(entity *Entity, elapsedTime time.Duration, input CarInput) {
	// TODO: Move constants as part of car descriptor
	const turnSpeed = 100       // FIXME ORIGINAL: 120
	const returnSpeed = 50      // FIXME ORIGINAL: 60
	const maxWheelAngle = 30    // FIXME ORIGINAL: 30
	const maxAcceleration = 0.6 // FIXME ORIGINAL: 0.01
	const maxDeceleration = 0.3 // FIXME ORIGINAL: 0.005

	elapsedSeconds := float32(elapsedTime.Seconds())
	vehicle := entity.Vehicle

	switch {
	case input.TurnLeft == input.TurnRight:
		if vehicle.SteeringAngle > 0.001 {
			if vehicle.SteeringAngle -= sprec.Degrees(elapsedSeconds * returnSpeed); vehicle.SteeringAngle < 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
		if vehicle.SteeringAngle < -0.001 {
			if vehicle.SteeringAngle += sprec.Degrees(elapsedSeconds * returnSpeed); vehicle.SteeringAngle > 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
	case input.TurnLeft:
		if vehicle.SteeringAngle += sprec.Degrees(elapsedSeconds * turnSpeed); vehicle.SteeringAngle > sprec.Degrees(maxWheelAngle) {
			vehicle.SteeringAngle = sprec.Degrees(maxWheelAngle)
		}
	case input.TurnRight:
		if vehicle.SteeringAngle -= sprec.Degrees(elapsedSeconds * turnSpeed); vehicle.SteeringAngle < -sprec.Degrees(maxWheelAngle) {
			vehicle.SteeringAngle = -sprec.Degrees(maxWheelAngle)
		}
	}
	vehicle.Acceleration = 0.0
	if input.Forward {
		vehicle.Acceleration = maxAcceleration * elapsedSeconds
	}
	if input.Backward {
		vehicle.Acceleration = -maxDeceleration * elapsedSeconds
	}
	vehicle.HandbrakePulled = input.Handbrake
}

func (s *VehicleSystem) updateVehicle(entity *Entity, elapsedTime time.Duration) {
	car := entity.Vehicle
	flWheel := entity.Vehicle.FLWheel.Wheel
	frWheel := entity.Vehicle.FRWheel.Wheel
	blWheel := entity.Vehicle.BLWheel.Wheel
	brWheel := entity.Vehicle.BRWheel.Wheel

	flWheel.SteeringAngle = car.SteeringAngle
	frWheel.SteeringAngle = car.SteeringAngle

	s.calculateWheelRotation(car, flWheel)
	s.calculateWheelRotation(car, frWheel)
	s.calculateWheelRotation(car, blWheel)
	s.calculateWheelRotation(car, brWheel)
	s.accelerateCar(car, flWheel, frWheel, blWheel, brWheel)
	s.translateCar(car)

	s.checkWheelSideCollision(car, flWheel, sprec.InverseVec3(car.Orientation.VectorX), sprec.InverseVec3(car.Orientation.VectorZ))
	s.checkWheelSideCollision(car, frWheel, car.Orientation.VectorX, sprec.InverseVec3(car.Orientation.VectorZ))
	s.checkWheelSideCollision(car, blWheel, sprec.InverseVec3(car.Orientation.VectorX), car.Orientation.VectorZ)
	s.checkWheelSideCollision(car, brWheel, car.Orientation.VectorX, car.Orientation.VectorZ)
	s.checkWheelBottomCollision(car, flWheel)
	s.checkWheelBottomCollision(car, frWheel)
	s.checkWheelBottomCollision(car, blWheel)
	s.checkWheelBottomCollision(car, brWheel)

	s.updateCarModelMatrix(entity)
	s.updateWheelModelMatrix(entity, entity.Vehicle.FLWheel)
	s.updateWheelModelMatrix(entity, entity.Vehicle.FRWheel)
	s.updateWheelModelMatrix(entity, entity.Vehicle.BLWheel)
	s.updateWheelModelMatrix(entity, entity.Vehicle.BRWheel)
}

func (s *VehicleSystem) calculateWheelRotation(car *Vehicle, wheel *Wheel) {
	// Handbrake locks all wheels
	if wheel.IsGrounded && !car.HandbrakePulled {
		wheelOrientation := car.Orientation
		wheelOrientation.Rotate(sprec.ResizedVec3(wheelOrientation.VectorY, wheel.SteeringAngle.Radians()))
		wheelTravelDistance := sprec.Vec3Dot(car.Velocity, wheelOrientation.VectorZ)
		s.rotateWheel(wheel, wheelTravelDistance)
	}

	// Add some wheel slip if we are accelerating
	if sprec.Abs(car.Acceleration) > 0.0001 && wheel.IsDriven {
		s.rotateWheel(wheel, sprec.Sign(car.Acceleration)*skidWheelSpeed)
	}
}

func (s *VehicleSystem) accelerateCar(car *Vehicle, flWheel, frWheel, blWheel, brWheel *Wheel) {
	shouldAccelerate :=
		(flWheel.IsDriven && flWheel.IsGrounded) ||
			(frWheel.IsDriven && frWheel.IsGrounded) ||
			(blWheel.IsDriven && blWheel.IsGrounded) ||
			(brWheel.IsDriven && brWheel.IsGrounded)

	if shouldAccelerate {
		acceleration := sprec.ResizedVec3(car.Orientation.VectorZ, car.Acceleration)
		car.Velocity = sprec.Vec3Sum(car.Velocity, acceleration)
	}

	shouldDecelerate := car.HandbrakePulled &&
		(flWheel.IsGrounded || frWheel.IsGrounded || blWheel.IsGrounded || brWheel.IsGrounded)

	if shouldDecelerate {
		car.Velocity = sprec.Vec3Prod(car.Velocity, wheelFriction)
	}

	car.Velocity = sprec.Vec3Diff(car.Velocity, sprec.NewVec3(0.0, gravity, 0.0))
	car.Velocity = sprec.Vec3Diff(car.Velocity, sprec.Vec3Prod(car.Velocity, speedFriction2))

	// no idea how I originally go to this
	magicVelocity := sprec.Vec3Dot(car.Velocity, car.Orientation.VectorZ) + sprec.Vec3Dot(car.Velocity, car.Orientation.VectorX)*sprec.Sin(car.SteeringAngle)
	var turnAcceleration float32
	if !car.HandbrakePulled {
		turnAcceleration = car.SteeringAngle.Degrees() * magicVelocity / (100.0 + 2.0*magicVelocity*magicVelocity)
	}
	turnAcceleration += car.SteeringAngle.Degrees() * car.Acceleration * 0.6
	turnAcceleration /= 2.0

	if flWheel.IsGrounded || frWheel.IsGrounded {
		car.AngularVelocity = sprec.Vec3Sum(car.AngularVelocity, sprec.ResizedVec3(car.Orientation.VectorY, turnAcceleration))
	}
	car.AngularVelocity = sprec.Vec3Prod(car.AngularVelocity, rotationFriction)
	car.AngularVelocity = sprec.Vec3Diff(car.AngularVelocity, sprec.Vec3Prod(car.AngularVelocity, rotationFriction2))
}

func (s *VehicleSystem) translateCar(car *Vehicle) {
	car.Position = sprec.Vec3Sum(car.Position, car.Velocity)
	car.Orientation.Rotate(car.AngularVelocity)
}

func (s *VehicleSystem) checkWheelSideCollision(car *Vehicle, wheel *Wheel, dirX sprec.Vec3, dirZ sprec.Vec3) {
	pa := s.wheelAbsolutePosition(car, wheel)

	p2 := sprec.Vec3Diff(pa, sprec.ResizedVec3(dirX, wheel.Length))
	p1 := sprec.Vec3Sum(pa, sprec.ResizedVec3(dirX, 1.2))
	result, active := s.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		f := sprec.ResizedVec3(result.Normal(), sprec.Abs(result.BottomHeight()))
		car.Position = sprec.Vec3Sum(car.Position, f)
		collisionVelocity := sprec.Vec3Dot(result.Normal(), car.Velocity)
		car.Velocity = sprec.Vec3Diff(car.Velocity, sprec.Vec3Prod(result.Normal(), collisionVelocity))
	}

	p2 = sprec.Vec3Sum(pa, sprec.ResizedVec3(dirZ, sprec.Abs(wheel.Radius)))
	p1 = sprec.Vec3Diff(pa, sprec.ResizedVec3(dirZ, 1.2))
	result, active = s.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		f := sprec.ResizedVec3(result.Normal(), sprec.Abs(result.BottomHeight()))
		car.Position = sprec.Vec3Sum(car.Position, f)
		collisionVelocity := sprec.Vec3Dot(result.Normal(), car.Velocity)
		car.Velocity = sprec.Vec3Diff(car.Velocity, sprec.Vec3Prod(result.Normal(), collisionVelocity))
	}
}

func (s *VehicleSystem) checkWheelBottomCollision(car *Vehicle, wheel *Wheel) {
	wheel.IsGrounded = false

	pa := sprec.Vec3Sum(car.Position, car.Orientation.MulVec3(wheel.AnchorPosition))
	p1 := sprec.Vec3Sum(pa, sprec.ResizedVec3(car.Orientation.VectorY, 3.6))
	p2 := sprec.Vec3Diff(pa, sprec.ResizedVec3(car.Orientation.VectorY, wheel.Radius+maxSuspensionLength))

	result, active := s.stage.CheckCollision(collision.MakeLine(p1, p2))
	if !active {
		wheel.SuspensionLength = maxSuspensionLength
		return
	}

	dis := sprec.Vec3Diff(result.Intersection(), p1).Length()
	if dis > 3.6+wheel.Radius {
		wheel.SuspensionLength = dis - (3.6 + wheel.Radius)
	} else {
		wheel.SuspensionLength = 0.0
	}
	wheel.IsGrounded = true

	if dis < 3.6+wheel.Radius {
		dis2 := sprec.Vec3Dot(result.Normal(), car.Velocity)
		f := sprec.ResizedVec3(result.Normal(), dis2)
		car.Velocity = sprec.Vec3Diff(car.Velocity, f)
		f = sprec.ResizedVec3(car.Orientation.VectorY, 3.6+wheel.Radius-dis)
		car.Position = sprec.Vec3Sum(car.Position, f)
	}

	wheelAbsolutePosition := s.wheelAbsolutePosition(car, wheel)
	relativePosition := sprec.Vec3Diff(wheelAbsolutePosition, car.Position)
	cross := sprec.Vec3Cross(result.Normal(), sprec.InverseVec3(relativePosition))
	cross = sprec.UnitVec3(cross)

	koef := sprec.ZeroVec2()
	koef.Y = sprec.Vec3Dot(sprec.InverseVec3(result.Intersection()), result.Intersection())
	tmp := sprec.InverseVec3(result.Intersection()).Length()
	koef.X = sprec.Sqrt(tmp*tmp - koef.Y*koef.Y)

	if koef.Length() > 0.0000001 {
		koef = sprec.UnitVec2(koef)
	} else {
		koef.X = 1.0
		koef.Y = 0.0
	}

	if sprec.Abs(koef.X) > 0.0000001 {
		cross = sprec.ResizedVec3(cross, koef.X*koef.X*(1.0-(dis-(3.6+wheel.Radius))/maxSuspensionLength)/10)
		car.AngularVelocity = sprec.Vec3Sum(car.AngularVelocity, cross)
	}
}

func (s *VehicleSystem) wheelAbsolutePosition(car *Vehicle, wheel *Wheel) sprec.Vec3 {
	worldPosition := car.Position
	worldPosition = sprec.Vec3Sum(worldPosition, car.Orientation.MulVec3(wheel.AnchorPosition))
	worldPosition = sprec.Vec3Sum(worldPosition, car.Orientation.MulVec3(sprec.NewVec3(0.0, -wheel.SuspensionLength, 0.0)))
	return worldPosition
}

func (s *VehicleSystem) rotateWheel(wheel *Wheel, speed float32) {
	wheel.RotationAngle += sprec.Radians(speed / wheel.Radius)
}

func (s *VehicleSystem) updateCarModelMatrix(car *Entity) {
	vehicleComp := car.Vehicle
	if car.Transform != nil {
		car.Transform.Position = vehicleComp.Position
		// car.Transform.Orientation = vehicleComp.Orientation
	}
}

func (s *VehicleSystem) updateWheelModelMatrix(car *Entity, wheel *Entity) {
	vehicleComp := car.Vehicle
	wheelComp := wheel.Wheel

	modelMatrix := sprec.Mat4MultiProd(
		sprec.TransformationMat4(
			vehicleComp.Orientation.VectorX,
			vehicleComp.Orientation.VectorY,
			vehicleComp.Orientation.VectorZ,
			vehicleComp.Position,
		),
		sprec.TranslationMat4(wheelComp.AnchorPosition.X, wheelComp.AnchorPosition.Y-wheelComp.SuspensionLength, wheelComp.AnchorPosition.Z),
		sprec.RotationMat4(wheelComp.SteeringAngle, 0.0, 1.0, 0.0),
		sprec.RotationMat4(wheelComp.RotationAngle, 1.0, 0.0, 0.0),
	)

	if wheel.Transform != nil {
		wheel.Transform.Position = modelMatrix.Translation()
		// wheel.Transform.Orientation = Orientation{
		// 	VectorX: modelMatrix.OrientationX(),
		// 	VectorY: modelMatrix.OrientationY(),
		// 	VectorZ: modelMatrix.OrientationZ(),
		// }
	}
}
