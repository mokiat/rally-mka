package system

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/internal/engine/physics"
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
	// TODO: Move constants as part of car descriptor
	// const turnSpeed = 100       // FIXME ORIGINAL: 120
	const turnSpeed = 50 // FIXME ORIGINAL: 120
	// const returnSpeed = 50      // FIXME ORIGINAL: 60
	const returnSpeed = 150     // FIXME ORIGINAL: 60
	const maxWheelAngle = 30    // FIXME ORIGINAL: 30
	const maxAcceleration = 0.8 // FIXME ORIGINAL: 0.6
	const maxDeceleration = 0.6 // FIXME ORIGINAL: 0.3

	elapsedSeconds := float32(elapsedTime.Seconds())

	switch {
	case input.TurnLeft == input.TurnRight:
		if car.SteeringAngle > 0.001 {
			if car.SteeringAngle -= sprec.Degrees(elapsedSeconds * returnSpeed); car.SteeringAngle < 0.0 {
				car.SteeringAngle = 0.0
			}
		}
		if car.SteeringAngle < -0.001 {
			if car.SteeringAngle += sprec.Degrees(elapsedSeconds * returnSpeed); car.SteeringAngle > 0.0 {
				car.SteeringAngle = 0.0
			}
		}
	case input.TurnLeft:
		if car.SteeringAngle += sprec.Degrees(elapsedSeconds * turnSpeed); car.SteeringAngle > sprec.Degrees(maxWheelAngle) {
			car.SteeringAngle = sprec.Degrees(maxWheelAngle)
		}
	case input.TurnRight:
		if car.SteeringAngle -= sprec.Degrees(elapsedSeconds * turnSpeed); car.SteeringAngle < -sprec.Degrees(maxWheelAngle) {
			car.SteeringAngle = -sprec.Degrees(maxWheelAngle)
		}
	}
	car.Acceleration = 0.0
	if input.Forward {
		car.Acceleration = maxAcceleration * elapsedSeconds
	}
	if input.Backward {
		car.Acceleration = -maxDeceleration * elapsedSeconds
	}
	car.HandbrakePulled = input.Handbrake
}

func (s *CarSystem) updateCar(car *ecs.Car) {
	flRotation := car.FLWheelRotation.(*physics.MatchAxisConstraint)
	frRotation := car.FRWheelRotation.(*physics.MatchAxisConstraint)

	rotationQuat := sprec.RotationQuat(car.SteeringAngle, sprec.BasisYVec3())
	flRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())
	frRotation.FirstBodyAxis = sprec.QuatVec3Rotation(rotationQuat, sprec.BasisXVec3())

	// FIXME: Acceleration, however, it gets erased at the moment, hence velocity

	// TODO: Remove, just for debugging
	if car.HandbrakePulled {
		car.Body.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.Body.Physics.Body.AngularVelocity, sprec.NewVec3(0.0, 0.0, 0.1))
		car.Body.Physics.Body.Velocity = sprec.Vec3Sum(car.Body.Physics.Body.Velocity, sprec.NewVec3(0.0, 0.2, 0.0))
	}

	// FIXME: With rotation this is no-longer correct as the Z axis moves around, making the wheel wobble
	// car.FLWheel.Motion.Velocity = sprec.Vec3Sum(car.FLWheel.Motion.Velocity, sprec.Vec3Prod(car.FLWheel.Transform.Orientation.OrientationZ(), car.Acceleration*20))
	// car.FRWheel.Motion.Velocity = sprec.Vec3Sum(car.FRWheel.Motion.Velocity, sprec.Vec3Prod(car.FRWheel.Transform.Orientation.OrientationZ(), car.Acceleration*20))

	accelRatio := float32(250)
	// accelRatio := float32(250)
	if car.Acceleration > 0.0001 {
		if sprec.Vec3Dot(car.Body.Physics.Body.Velocity, car.Body.Render.Matrix.OrientationZ()) < -5.0 {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FLWheel.Physics.Body.AngularVelocity, 0.85)
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FRWheel.Physics.Body.AngularVelocity, 0.85)
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BLWheel.Physics.Body.AngularVelocity, 0.90)
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BRWheel.Physics.Body.AngularVelocity, 0.90)
		} else {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
		}
	}
	if car.Acceleration < -0.001 {
		if sprec.Vec3Dot(car.Body.Physics.Body.Velocity, car.Body.Render.Matrix.OrientationZ()) > 5.0 {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FLWheel.Physics.Body.AngularVelocity, 0.85)
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.FRWheel.Physics.Body.AngularVelocity, 0.85)
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BLWheel.Physics.Body.AngularVelocity, 0.90)
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Prod(car.BRWheel.Physics.Body.AngularVelocity, 0.90)
		} else {
			car.FLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.FRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.FRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.FRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.BLWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BLWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BLWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
			car.BRWheel.Physics.Body.AngularVelocity = sprec.Vec3Sum(car.BRWheel.Physics.Body.AngularVelocity, sprec.Vec3Prod(car.BRWheel.Physics.Body.Orientation.OrientationX(), car.Acceleration*accelRatio))
		}
	}
}
