package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/input"
)

const (
	steeringSpeed        = 80
	steeringRestoreSpeed = steeringSpeed * 2
)

func NewVehicleSystem(ecsManager *Manager) *VehicleSystem {
	return &VehicleSystem{
		ecsManager: ecsManager,
	}
}

type VehicleSystem struct {
	ecsManager *Manager
}

func (s *VehicleSystem) Update(ctx game.UpdateContext) {
	for _, entity := range s.ecsManager.Entities() {
		if vehicle := entity.Vehicle; vehicle != nil {
			if entity.PlayerControl != nil {
				switch {
				case ctx.Gamepad.Available:
					s.updateVehicleControlGamepad(vehicle, ctx)
				default:
					s.updateVehicleControlKeyboard(vehicle, ctx)
				}
			}
			s.updateVehiclePhysics(vehicle, ctx)
		}
	}
}

func (s *VehicleSystem) updateVehicleControlGamepad(vehicle *Vehicle, ctx game.UpdateContext) {
	gamepad := ctx.Gamepad

	steeringAmount := gamepad.LeftStickX * sprec.Abs(gamepad.LeftStickX)
	vehicle.SteeringAngle = -sprec.Degrees(steeringAmount * vehicle.MaxSteeringAngle.Degrees())
	vehicle.Acceleration = gamepad.RightTrigger
	vehicle.Deceleration = gamepad.LeftTrigger
	vehicle.Recover = gamepad.LeftBumper
}

func (s *VehicleSystem) updateVehicleControlKeyboard(vehicle *Vehicle, ctx game.UpdateContext) {
	elapsedSeconds := float32(ctx.ElapsedTime.Seconds())
	keyboard := ctx.Keyboard
	isSteerLeft := keyboard.IsPressed(input.KeyLeft)
	isSteerRight := keyboard.IsPressed(input.KeyRight)
	isAccelerate := keyboard.IsPressed(input.KeyUp)
	isDecelerate := keyboard.IsPressed(input.KeyDown)
	vehicle.Recover = keyboard.IsPressed(input.KeyEnter)

	autoMaxSteeringAngle := sprec.Degrees(vehicle.MaxSteeringAngle.Degrees() / (1.0 + 0.05*vehicle.Chassis.Body.Velocity.Length()))
	switch {
	case isSteerLeft == isSteerRight:
		if vehicle.SteeringAngle > 0.001 {
			vehicle.SteeringAngle -= sprec.Degrees(elapsedSeconds * steeringRestoreSpeed)
			if vehicle.SteeringAngle < 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
		if vehicle.SteeringAngle < -0.001 {
			vehicle.SteeringAngle += sprec.Degrees(elapsedSeconds * steeringRestoreSpeed)
			if vehicle.SteeringAngle > 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
	case isSteerLeft:
		vehicle.SteeringAngle += sprec.Degrees(elapsedSeconds * steeringSpeed)
		if vehicle.SteeringAngle > autoMaxSteeringAngle {
			vehicle.SteeringAngle = autoMaxSteeringAngle
		}
	case isSteerRight:
		vehicle.SteeringAngle -= sprec.Degrees(elapsedSeconds * steeringSpeed)
		if vehicle.SteeringAngle < -autoMaxSteeringAngle {
			vehicle.SteeringAngle = -autoMaxSteeringAngle
		}
	}

	if isAccelerate {
		vehicle.Acceleration = 0.8
	} else {
		vehicle.Acceleration = 0.0
	}
	if isDecelerate {
		vehicle.Deceleration = 0.8
	} else {
		vehicle.Deceleration = 0.0
	}
}

func (s *VehicleSystem) updateVehiclePhysics(vehicle *Vehicle, ctx game.UpdateContext) {
	elapsedSeconds := float32(ctx.ElapsedTime.Seconds())

	if vehicle.Recover {
		vehicle.Chassis.Body.AngularVelocity = sprec.Vec3Sum(vehicle.Chassis.Body.AngularVelocity, sprec.NewVec3(0.0, 0.0, 0.1))
		vehicle.Chassis.Body.Velocity = sprec.Vec3Sum(vehicle.Chassis.Body.Velocity, sprec.NewVec3(0.0, 0.2, 0.0))
	}

	steeringQuat := sprec.RotationQuat(vehicle.SteeringAngle, sprec.BasisYVec3())
	isMovingForward := sprec.Vec3Dot(vehicle.Chassis.Body.Velocity, vehicle.Chassis.Body.Orientation.OrientationZ()) > 5.0
	isMovingBackward := sprec.Vec3Dot(vehicle.Chassis.Body.Velocity, vehicle.Chassis.Body.Orientation.OrientationZ()) < -5.0

	for _, wheel := range vehicle.Wheels {
		if wheel.RotationConstraint != nil {
			wheel.RotationConstraint.FirstBodyAxis = sprec.QuatVec3Rotation(steeringQuat, sprec.BasisXVec3())
		}

		if vehicle.Acceleration > 0.0 {
			if isMovingBackward {
				if wheelVelocity := sprec.Vec3Dot(wheel.Body.AngularVelocity, wheel.Body.Orientation.OrientationX()); wheelVelocity < 0.0 {
					correction := sprec.Max(-vehicle.Acceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.AngularVelocity = sprec.Vec3Prod(wheel.Body.AngularVelocity, 1.0-correction/wheelVelocity)
				}
			} else {
				wheel.Body.AngularVelocity = sprec.Vec3Sum(wheel.Body.AngularVelocity,
					sprec.Vec3Prod(wheel.Body.Orientation.OrientationX(), vehicle.Acceleration*wheel.AccelerationVelocity*elapsedSeconds),
				)
			}
		}

		if vehicle.Deceleration > 0.0 {
			if isMovingForward {
				if wheelVelocity := sprec.Vec3Dot(wheel.Body.AngularVelocity, wheel.Body.Orientation.OrientationX()); wheelVelocity > 0.0 {
					correction := sprec.Min(vehicle.Deceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.AngularVelocity = sprec.Vec3Prod(wheel.Body.AngularVelocity, 1.0-correction/wheelVelocity)
				}
			} else {
				wheel.Body.AngularVelocity = sprec.Vec3Sum(wheel.Body.AngularVelocity,
					sprec.Vec3Prod(wheel.Body.Orientation.OrientationX(), -vehicle.Deceleration*wheel.AccelerationVelocity*elapsedSeconds),
				)
			}
		}
	}
}
