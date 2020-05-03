package system

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/internal/engine/physics"
)

const (
	// TODO: Move constants as part of car descriptor
	carSteeringAngle        = 40
	carSteeringSpeed        = 80
	carSteeringRestoreSpeed = 150

	carFrontAcceleration        = 160
	carRearAcceleration         = 80
	carReverseAccelerationRatio = 0.75

	carFrontBrakeRatio = 0.1
	carRearBrakeRatio  = 0.1
)

func NewCarSystem(ecsManager *ecs.Manager) *CarSystem {
	return &CarSystem{
		ecsManager: ecsManager,
	}
}

type CarSystem struct {
	ecsManager *ecs.Manager
}

func (s *CarSystem) Update(elapsedTime time.Duration, input ecs.CarInput) {
	for _, entity := range s.ecsManager.Entities() {
		if car := entity.Car; car != nil {
			s.updateCarInput(car, elapsedTime, input)
			s.updateCar(car)
		}
	}
}

func (s *CarSystem) updateCarInput(car *ecs.Car, elapsedTime time.Duration, input ecs.CarInput) {
	elapsedSeconds := float32(elapsedTime.Seconds())

	actualSteeringAngle := carSteeringAngle / (1.0 + 0.03*car.Body.Physics.Body.Velocity.Length())

	switch {
	case input.TurnLeft == input.TurnRight:
		if car.SteeringAngle > 0.001 {
			if car.SteeringAngle -= sprec.Degrees(elapsedSeconds * carSteeringRestoreSpeed); car.SteeringAngle < 0.0 {
				car.SteeringAngle = 0.0
			}
		}
		if car.SteeringAngle < -0.001 {
			if car.SteeringAngle += sprec.Degrees(elapsedSeconds * carSteeringRestoreSpeed); car.SteeringAngle > 0.0 {
				car.SteeringAngle = 0.0
			}
		}
	case input.TurnLeft:
		if car.SteeringAngle += sprec.Degrees(elapsedSeconds * carSteeringSpeed); car.SteeringAngle > sprec.Degrees(actualSteeringAngle) {
			car.SteeringAngle = sprec.Degrees(actualSteeringAngle)
		}
	case input.TurnRight:
		if car.SteeringAngle -= sprec.Degrees(elapsedSeconds * carSteeringSpeed); car.SteeringAngle < -sprec.Degrees(actualSteeringAngle) {
			car.SteeringAngle = -sprec.Degrees(actualSteeringAngle)
		}
	}
	car.Acceleration = 0.0
	if input.Forward {
		car.Acceleration = 1.0 * elapsedSeconds
	}
	if input.Backward {
		car.Acceleration = -1.0 * carReverseAccelerationRatio * elapsedSeconds
	}
	car.HandbrakePulled = input.Handbrake
}

func (s *CarSystem) updateCar(car *ecs.Car) {
	flRotation := car.FLWheelRotation.(*physics.MatchAxisConstraint)
	frRotation := car.FRWheelRotation.(*physics.MatchAxisConstraint)

	rotationQuat := sprec.RotationQuat(car.SteeringAngle, sprec.BasisYVec3())
	flRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
	frRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())

	// TODO: Remove, just for debugging
	if car.HandbrakePulled {
		car.Body.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.Body.Physics.Body.AngularVelocity, sprec.NewVec3(0.0, 0.0, 0.1))
		car.Body.Physics.Body.Velocity = sprec.Vec3Sum(car.Body.Physics.Body.Velocity, sprec.NewVec3(0.0, 0.2, 0.0))
	}

	if car.Acceleration > 0.0001 {
		if sprec.Vec3Dot(car.Body.Physics.Body.Velocity, car.Body.Render.Matrix.OrientationZ()) < -5.0 {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FLWheel.Physics.Body.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FRWheel.Physics.Body.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BLWheel.Physics.Body.AngularVelocity, 1.0-carRearBrakeRatio)
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BRWheel.Physics.Body.AngularVelocity, 1.0-carRearBrakeRatio)
		} else {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carFrontAcceleration))
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carFrontAcceleration))
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carRearAcceleration))
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carRearAcceleration))
		}
	}
	if car.Acceleration < -0.001 {
		if sprec.Vec3Dot(car.Body.Physics.Body.Velocity, car.Body.Render.Matrix.OrientationZ()) > 5.0 {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FLWheel.Physics.Body.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FRWheel.Physics.Body.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BLWheel.Physics.Body.AngularVelocity, 1.0-carRearBrakeRatio)
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BRWheel.Physics.Body.AngularVelocity, 1.0-carRearBrakeRatio)
		} else {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carFrontAcceleration))
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carFrontAcceleration))
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carRearAcceleration))
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*carRearAcceleration))
		}
	}
}
