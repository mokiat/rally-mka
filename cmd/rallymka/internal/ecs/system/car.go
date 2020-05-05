package system

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

const (
	// TODO: Move constants as part of car descriptor
	carMaxSteeringAngle     = 40
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
	elapsedSeconds := float32(elapsedTime.Seconds())

	for _, entity := range s.ecsManager.Entities() {
		if car := entity.Car; car != nil && entity.HumanInput {
			s.updateCarSteering(car, elapsedSeconds, input)
			s.updateCarAcceleration(car, elapsedSeconds, input)
		}
	}
}

func (s *CarSystem) updateCarSteering(car *ecs.Car, elapsedSeconds float32, input ecs.CarInput) {
	actualSteeringAngle := carMaxSteeringAngle / (1.0 + 0.03*car.Chassis.Velocity.Length())
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

	rotationQuat := sprec.RotationQuat(car.SteeringAngle, sprec.BasisYVec3())
	car.FLWheelRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
	car.FRWheelRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
}

func (s *CarSystem) updateCarAcceleration(car *ecs.Car, elapsedSeconds float32, input ecs.CarInput) {
	// TODO: Remove, just for debugging
	if input.Handbrake {
		car.Chassis.AngularVelocity = sprec.Vec3Sum(car.Chassis.AngularVelocity, sprec.NewVec3(0.0, 0.0, 0.1))
		car.Chassis.Velocity = sprec.Vec3Sum(car.Chassis.Velocity, sprec.NewVec3(0.0, 0.2, 0.0))
	}

	if input.Forward {
		if sprec.Vec3Dot(car.Chassis.Velocity, car.Chassis.Orientation.OrientationZ()) < -5.0 {
			car.FLWheel.AngularVelocity = sprec.Vec3Prod(car.FLWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.FRWheel.AngularVelocity = sprec.Vec3Prod(car.FRWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.BLWheel.AngularVelocity = sprec.Vec3Prod(car.BLWheel.AngularVelocity, 1.0-carRearBrakeRatio)
			car.BRWheel.AngularVelocity = sprec.Vec3Prod(car.BRWheel.AngularVelocity, 1.0-carRearBrakeRatio)
		} else {
			car.FLWheel.AngularVelocity = sprec.Vec3Sum(car.FLWheel.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds))
			car.FRWheel.AngularVelocity = sprec.Vec3Sum(car.FRWheel.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds))
			car.BLWheel.AngularVelocity = sprec.Vec3Sum(car.BLWheel.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds))
			car.BRWheel.AngularVelocity = sprec.Vec3Sum(car.BRWheel.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds))
		}
	}
	if input.Backward {
		if sprec.Vec3Dot(car.Chassis.Velocity, car.Chassis.Orientation.OrientationZ()) > 5.0 {
			car.FLWheel.AngularVelocity = sprec.Vec3Prod(car.FLWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.FRWheel.AngularVelocity = sprec.Vec3Prod(car.FRWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
			car.BLWheel.AngularVelocity = sprec.Vec3Prod(car.BLWheel.AngularVelocity, 1.0-carRearBrakeRatio)
			car.BRWheel.AngularVelocity = sprec.Vec3Prod(car.BRWheel.AngularVelocity, 1.0-carRearBrakeRatio)
		} else {
			car.FLWheel.AngularVelocity = sprec.Vec3Sum(car.FLWheel.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Orientation.OrientationX(), -carFrontAcceleration*carReverseAccelerationRatio*elapsedSeconds))
			car.FRWheel.AngularVelocity = sprec.Vec3Sum(car.FRWheel.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Orientation.OrientationX(), -carFrontAcceleration*carReverseAccelerationRatio*elapsedSeconds))
			car.BLWheel.AngularVelocity = sprec.Vec3Sum(car.BLWheel.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Orientation.OrientationX(), -carRearAcceleration*carReverseAccelerationRatio*elapsedSeconds))
			car.BRWheel.AngularVelocity = sprec.Vec3Sum(car.BRWheel.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Orientation.OrientationX(), -carRearAcceleration*carReverseAccelerationRatio*elapsedSeconds))
		}
	}
}
