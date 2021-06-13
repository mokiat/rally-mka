package ecssys

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
)

const (
	steeringSpeed        = 80
	steeringRestoreSpeed = steeringSpeed * 2
)

func NewVehicleSystem(ecsScene *ecs.Scene) *VehicleSystem {
	return &VehicleSystem{
		ecsScene: ecsScene,
	}
}

type VehicleSystem struct {
	ecsScene *ecs.Scene

	isSteerLeft  bool
	isSteerRight bool
	isAccelerate bool
	isDecelerate bool
	isRecover    bool
}

func (s *VehicleSystem) OnKeyboardEvent(event app.KeyboardEvent) bool {
	active := event.Type != app.KeyboardEventTypeKeyUp
	switch event.Code {
	case app.KeyCodeArrowLeft:
		s.isSteerLeft = active
		return true
	case app.KeyCodeArrowRight:
		s.isSteerRight = active
		return true
	case app.KeyCodeArrowUp:
		s.isAccelerate = active
		return true
	case app.KeyCodeArrowDown:
		s.isDecelerate = active
		return true
	case app.KeyCodeEnter:
		s.isRecover = active
		return true
	}
	return false
}

func (s *VehicleSystem) Update(elapsedSeconds float32, gamepad *app.GamepadState) {
	result := s.ecsScene.Find(ecs.Having(ecscomp.VehicleComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()
		vehicle := ecscomp.GetVehicle(entity)
		if ecscomp.GetPlayerControl(entity) != nil {
			if gamepad != nil {
				s.updateVehicleControlGamepad(vehicle, elapsedSeconds, gamepad)
			} else {
				s.updateVehicleControlKeyboard(vehicle, elapsedSeconds)
			}
		}
		s.updateVehiclePhysics(vehicle, elapsedSeconds)
	}
}

func (s *VehicleSystem) updateVehicleControlGamepad(vehicle *ecscomp.Vehicle, elapsedSeconds float32, gamepad *app.GamepadState) {
	steeringAmount := gamepad.LeftStickX * sprec.Abs(gamepad.LeftStickX)
	vehicle.SteeringAngle = -sprec.Degrees(steeringAmount * vehicle.MaxSteeringAngle.Degrees())
	vehicle.Acceleration = gamepad.RightTrigger
	vehicle.Deceleration = gamepad.LeftTrigger
	vehicle.Recover = gamepad.LeftBumper
}

func (s *VehicleSystem) updateVehicleControlKeyboard(vehicle *ecscomp.Vehicle, elapsedSeconds float32) {
	vehicle.Recover = s.isRecover

	autoMaxSteeringAngle := sprec.Degrees(vehicle.MaxSteeringAngle.Degrees() / (1.0 + 0.05*vehicle.Chassis.Body.Velocity().Length()))
	switch {
	case s.isSteerLeft == s.isSteerRight:
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
	case s.isSteerLeft:
		vehicle.SteeringAngle += sprec.Degrees(elapsedSeconds * steeringSpeed)
		if vehicle.SteeringAngle > autoMaxSteeringAngle {
			vehicle.SteeringAngle = autoMaxSteeringAngle
		}
	case s.isSteerRight:
		vehicle.SteeringAngle -= sprec.Degrees(elapsedSeconds * steeringSpeed)
		if vehicle.SteeringAngle < -autoMaxSteeringAngle {
			vehicle.SteeringAngle = -autoMaxSteeringAngle
		}
	}

	if s.isAccelerate {
		vehicle.Acceleration = 0.8
	} else {
		vehicle.Acceleration = 0.0
	}
	if s.isDecelerate {
		vehicle.Deceleration = 0.8
	} else {
		vehicle.Deceleration = 0.0
	}
}

func (s *VehicleSystem) updateVehiclePhysics(vehicle *ecscomp.Vehicle, elapsedSeconds float32) {
	if vehicle.Recover {
		vehicle.Chassis.Body.SetAngularVelocity(sprec.Vec3Sum(vehicle.Chassis.Body.AngularVelocity(), sprec.NewVec3(0.0, 0.0, 0.1)))
		vehicle.Chassis.Body.SetVelocity(sprec.Vec3Sum(vehicle.Chassis.Body.Velocity(), sprec.NewVec3(0.0, 0.2, 0.0)))
	}

	steeringQuat := sprec.RotationQuat(vehicle.SteeringAngle, sprec.BasisYVec3())
	isMovingForward := sprec.Vec3Dot(vehicle.Chassis.Body.Velocity(), vehicle.Chassis.Body.Orientation().OrientationZ()) > 5.0
	isMovingBackward := sprec.Vec3Dot(vehicle.Chassis.Body.Velocity(), vehicle.Chassis.Body.Orientation().OrientationZ()) < -5.0

	for _, wheel := range vehicle.Wheels {
		if wheel.RotationConstraint != nil {
			wheel.RotationConstraint.SetPrimaryAxis(sprec.QuatVec3Rotation(steeringQuat, sprec.BasisXVec3()))
		}

		if vehicle.Acceleration > 0.0 {
			if isMovingBackward {
				if wheelVelocity := sprec.Vec3Dot(wheel.Body.AngularVelocity(), wheel.Body.Orientation().OrientationX()); wheelVelocity < 0.0 {
					correction := sprec.Max(-vehicle.Acceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.SetAngularVelocity(sprec.Vec3Prod(wheel.Body.AngularVelocity(), 1.0-correction/wheelVelocity))
				}
			} else {
				wheel.Body.SetAngularVelocity(sprec.Vec3Sum(wheel.Body.AngularVelocity(),
					sprec.Vec3Prod(wheel.Body.Orientation().OrientationX(), vehicle.Acceleration*wheel.AccelerationVelocity*elapsedSeconds),
				))
			}
		}

		if vehicle.Deceleration > 0.0 {
			if isMovingForward {
				if wheelVelocity := sprec.Vec3Dot(wheel.Body.AngularVelocity(), wheel.Body.Orientation().OrientationX()); wheelVelocity > 0.0 {
					correction := sprec.Min(vehicle.Deceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.SetAngularVelocity(sprec.Vec3Prod(wheel.Body.AngularVelocity(), 1.0-correction/wheelVelocity))
				}
			} else {
				wheel.Body.SetAngularVelocity(sprec.Vec3Sum(wheel.Body.AngularVelocity(),
					sprec.Vec3Prod(wheel.Body.Orientation().OrientationX(), -vehicle.Deceleration*wheel.AccelerationVelocity*elapsedSeconds),
				))
			}
		}
	}
}
