package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/input"
)

const (
	// TODO: Move constants as part of car descriptor
	carMaxSteeringAngle     = 30
	carSteeringSpeed        = 80
	carSteeringRestoreSpeed = 160

	carFrontAcceleration        = 100
	carRearAcceleration         = 200
	carReverseAccelerationRatio = 0.75

	carFrontBrakeRatio = 0.1
	carRearBrakeRatio  = 0.1
)

func NewCarSystem(ecsManager *Manager) *CarSystem {
	return &CarSystem{
		ecsManager: ecsManager,
	}
}

type CarSystem struct {
	ecsManager *Manager
}

func (s *CarSystem) Update(ctx game.UpdateContext) {
	for _, entity := range s.ecsManager.Entities() {
		if car := entity.Car; car != nil && entity.HumanInput {
			s.updateCarSteering(car, ctx)
			s.updateCarAcceleration(car, ctx)
		}
	}
}

func (s *CarSystem) updateCarSteering(car *Car, ctx game.UpdateContext) {
	if ctx.Gamepad.Available {
		gamepad := ctx.Gamepad
		car.SteeringAngle = -sprec.Degrees(gamepad.LeftStickX * sprec.Abs(gamepad.LeftStickX) * carMaxSteeringAngle)
	} else {
		elapsedSeconds := float32(ctx.ElapsedTime.Seconds())
		keyboard := ctx.Keyboard
		turnLeft := keyboard.IsPressed(input.KeyLeft)
		turnRight := keyboard.IsPressed(input.KeyRight)
		actualSteeringAngle := carMaxSteeringAngle / (1.0 + 0.2*car.Chassis.Velocity.Length())
		switch {
		case turnLeft == turnRight:
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
		case turnLeft:
			if car.SteeringAngle += sprec.Degrees(elapsedSeconds * carSteeringSpeed); car.SteeringAngle > sprec.Degrees(actualSteeringAngle) {
				car.SteeringAngle = sprec.Degrees(actualSteeringAngle)
			}
		case turnRight:
			if car.SteeringAngle -= sprec.Degrees(elapsedSeconds * carSteeringSpeed); car.SteeringAngle < -sprec.Degrees(actualSteeringAngle) {
				car.SteeringAngle = -sprec.Degrees(actualSteeringAngle)
			}
		}
	}

	rotationQuat := sprec.RotationQuat(car.SteeringAngle, sprec.BasisYVec3())
	car.FLWheelRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
	car.FRWheelRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
}

func (s *CarSystem) updateCarAcceleration(car *Car, ctx game.UpdateContext) {
	elapsedSeconds := float32(ctx.ElapsedTime.Seconds())

	// TODO: Remove, just for debugging
	if ctx.Keyboard.IsPressed(input.KeyEnter) {
		car.Chassis.AngularVelocity = sprec.Vec3Sum(car.Chassis.AngularVelocity, sprec.NewVec3(0.0, 0.0, 0.1))
		car.Chassis.Velocity = sprec.Vec3Sum(car.Chassis.Velocity, sprec.NewVec3(0.0, 0.2, 0.0))
	}

	if ctx.Gamepad.Available {
		gamepad := ctx.Gamepad
		if gamepad.LeftTrigger > 0.05 {
			car.FLWheel.AngularVelocity = sprec.Vec3Prod(car.FLWheel.AngularVelocity, 1.0-carFrontBrakeRatio*gamepad.LeftTrigger)
			car.FRWheel.AngularVelocity = sprec.Vec3Prod(car.FRWheel.AngularVelocity, 1.0-carFrontBrakeRatio*gamepad.LeftTrigger)
			car.BLWheel.AngularVelocity = sprec.Vec3Prod(car.BLWheel.AngularVelocity, 1.0-carRearBrakeRatio*gamepad.LeftTrigger)
			car.BRWheel.AngularVelocity = sprec.Vec3Prod(car.BRWheel.AngularVelocity, 1.0-carRearBrakeRatio*gamepad.LeftTrigger)
		}
		if gamepad.RightTrigger > 0.05 {
			car.FLWheel.AngularVelocity = sprec.Vec3Sum(car.FLWheel.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds*gamepad.RightTrigger))
			car.FRWheel.AngularVelocity = sprec.Vec3Sum(car.FRWheel.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds*gamepad.RightTrigger))
			car.BLWheel.AngularVelocity = sprec.Vec3Sum(car.BLWheel.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds*gamepad.RightTrigger))
			car.BRWheel.AngularVelocity = sprec.Vec3Sum(car.BRWheel.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds*gamepad.RightTrigger))
		}
	} else {
		if ctx.Keyboard.IsPressed(input.KeyUp) {
			if sprec.Vec3Dot(car.Chassis.Velocity, car.Chassis.Orientation.OrientationZ()) < -5.0 {
				car.FLWheel.AngularVelocity = sprec.Vec3Prod(car.FLWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
				car.FRWheel.AngularVelocity = sprec.Vec3Prod(car.FRWheel.AngularVelocity, 1.0-carFrontBrakeRatio)
				car.BLWheel.AngularVelocity = sprec.Vec3Prod(car.BLWheel.AngularVelocity, 1.0-carRearBrakeRatio)
				car.BRWheel.AngularVelocity = sprec.Vec3Prod(car.BRWheel.AngularVelocity, 1.0-carRearBrakeRatio)
			} else {
				if car.Chassis.Velocity.Length()*3.6 < 90 {
					car.FLWheel.AngularVelocity = sprec.Vec3Sum(car.FLWheel.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds))
					car.FRWheel.AngularVelocity = sprec.Vec3Sum(car.FRWheel.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Orientation.OrientationX(), carFrontAcceleration*elapsedSeconds))
					car.BLWheel.AngularVelocity = sprec.Vec3Sum(car.BLWheel.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds))
					car.BRWheel.AngularVelocity = sprec.Vec3Sum(car.BRWheel.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Orientation.OrientationX(), carRearAcceleration*elapsedSeconds))
				}
			}
		}
		if ctx.Keyboard.IsPressed(input.KeyDown) {
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
}
