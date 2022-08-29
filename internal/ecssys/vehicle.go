package ecssys

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/util/shape"
	"github.com/mokiat/rally-mka/internal/ecscomp"
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

	mouseLeft   bool
	mouseRight  bool
	mouseActive bool
	mouseX      int
	mouseY      int

	viewport      graphics.Viewport
	camera        *graphics.Camera
	graphicsScene *graphics.Scene
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

func (s *VehicleSystem) OnMouseEvent(event app.MouseEvent, viewport graphics.Viewport, camera *graphics.Camera, graphicsScene *graphics.Scene) bool {
	s.viewport = viewport
	s.camera = camera
	s.graphicsScene = graphicsScene

	switch event.Type {
	case app.MouseEventTypeDown:
		switch event.Button {
		case app.MouseButtonLeft:
			s.mouseLeft = true
		case app.MouseButtonMiddle:
			s.mouseActive = !s.mouseActive
		case app.MouseButtonRight:
			s.mouseRight = true
		}
		// if s.mouseLeft && s.mouseRight {
		// 	s.mouseLeft = false
		// 	s.mouseRight = false
		// 	s.mouseActive = !s.mouseActive
		// }
	case app.MouseEventTypeUp:
		switch event.Button {
		case app.MouseButtonLeft:
			s.mouseLeft = false
		case app.MouseButtonRight:
			s.mouseRight = false
		}
	case app.MouseEventTypeMove:
		if !s.mouseActive {
			return false
		}
		s.mouseX = event.X
		s.mouseY = event.Y
	}
	return true
}

func (s *VehicleSystem) Update(elapsedSeconds float64, gamepad *app.GamepadState) {
	result := s.ecsScene.Find(ecs.Having(ecscomp.VehicleComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()
		vehicle := ecscomp.GetVehicle(entity)
		if ecscomp.GetPlayerControl(entity) != nil {
			if gamepad != nil {
				s.updateVehicleControlGamepad(vehicle, elapsedSeconds, gamepad)
			} else {
				if s.mouseActive {
					s.updateVehicleControlMouse(vehicle, elapsedSeconds)
				} else {
					s.updateVehicleControlKeyboard(vehicle, elapsedSeconds)
				}
			}
		}
		s.updateVehiclePhysics(vehicle, elapsedSeconds)
	}
}

func (s *VehicleSystem) updateVehicleControlGamepad(vehicle *ecscomp.Vehicle, elapsedSeconds float64, gamepad *app.GamepadState) {
	steeringAmount := float64(gamepad.LeftStickX) * dprec.Abs(float64(gamepad.LeftStickX))
	vehicle.SteeringAngle = -dprec.Degrees(steeringAmount * vehicle.MaxSteeringAngle.Degrees())
	vehicle.Acceleration = float64(gamepad.RightTrigger)
	vehicle.Deceleration = float64(gamepad.LeftTrigger)
	vehicle.Recover = gamepad.CrossButton
}

func (s *VehicleSystem) updateVehicleControlMouse(vehicle *ecscomp.Vehicle, elapsedSeconds float64) {
	vehicle.Recover = s.isRecover

	orientation := vehicle.Chassis.Body.Orientation()
	orientationX := orientation.OrientationX()
	// orientationY := orientation.OrientationY()
	orientationZ := orientation.OrientationZ()

	position := vehicle.Chassis.Body.Position()
	a := dprec.Vec3Sum(dprec.Vec3Sum(
		dprec.Vec3Prod(orientationX, -1000),
		dprec.Vec3Prod(orientationZ, -1000),
	), position)
	b := dprec.Vec3Sum(dprec.Vec3Sum(
		dprec.Vec3Prod(orientationX, -1000),
		dprec.Vec3Prod(orientationZ, 1000),
	), position)
	c := dprec.Vec3Sum(dprec.Vec3Sum(
		dprec.Vec3Prod(orientationX, 1000),
		dprec.Vec3Prod(orientationZ, 1000),
	), position)
	d := dprec.Vec3Sum(dprec.Vec3Sum(
		dprec.Vec3Prod(orientationX, 1000),
		dprec.Vec3Prod(orientationZ, -1000),
	), position)

	// a := dprec.NewVec3(
	// 	-1000+position.X,
	// 	position.Y,
	// 	-1000+position.Z,
	// )
	// b := dprec.NewVec3(
	// 	-1000+position.X,
	// 	position.Y,
	// 	1000+position.Z,
	// )
	// c := dprec.NewVec3(
	// 	1000+position.X,
	// 	position.Y,
	// 	1000+position.Z,
	// )
	// d := dprec.NewVec3(
	// 	1000+position.X,
	// 	position.Y,
	// 	-1000+position.Z,
	// )

	surface := shape.NewStaticMesh([]shape.StaticTriangle{
		shape.NewStaticTriangle(a, b, c),
		shape.NewStaticTriangle(a, c, d),
	})

	line := s.graphicsScene.Ray(s.viewport, s.camera, s.mouseX, s.mouseY)

	result := shape.NewIntersectionResultSet(1)
	shape.CheckLineIntersection(line, surface, result)

	if result.Found() {
		intersection := result.Intersections()[0]
		mouseTarget := intersection.FirstContact
		// log.Info("CONTACT: %#v", intersection.FirstContact)

		delta := dprec.Vec3Diff(mouseTarget, vehicle.Chassis.Body.Position())
		delta.Y = 0.0

		forward := vehicle.Chassis.Body.Orientation().OrientationZ()
		forward.Y = 0.0

		sin := dprec.Vec3Cross(
			dprec.UnitVec3(delta),
			dprec.UnitVec3(forward),
		)
		angle := dprec.Angle(dprec.Sign(sin.Y)) * dprec.Asin(sin.Length())
		// log.Info("Angle: %.4f", angle.Degrees())

		vehicle.SteeringAngle = dprec.Clamp(-dprec.Angle(angle), -vehicle.MaxSteeringAngle, vehicle.MaxSteeringAngle)
	} else {
		vehicle.SteeringAngle = 0
	}
	// log.Info("Line: %#v, %#v", line.A(), line.B())

	// halfWidth := float64(width) / 2.0
	// halfHeight := float64(height) / 2.0
	// s.mouseTurn = 2 * (float64(event.X) - halfWidth) / halfWidth
	// s.mouseAcceleration = (halfHeight - float64(event.Y)) / halfHeight
	// return true

	if s.mouseLeft {
		vehicle.Acceleration = 0.8
	} else {
		vehicle.Acceleration = 0.0
	}
	if s.mouseRight {
		vehicle.Deceleration = 0.8
	} else {
		vehicle.Deceleration = 0.0
	}
}

func (s *VehicleSystem) updateVehicleControlKeyboard(vehicle *ecscomp.Vehicle, elapsedSeconds float64) {
	vehicle.Recover = s.isRecover

	autoMaxSteeringAngle := dprec.Degrees(vehicle.MaxSteeringAngle.Degrees() / (1.0 + 0.05*vehicle.Chassis.Body.Velocity().Length()))
	switch {
	case s.isSteerLeft == s.isSteerRight:
		if vehicle.SteeringAngle > 0.0 {
			vehicle.SteeringAngle -= dprec.Degrees(elapsedSeconds * steeringRestoreSpeed)
			if vehicle.SteeringAngle < 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
		if vehicle.SteeringAngle < 0.0 {
			vehicle.SteeringAngle += dprec.Degrees(elapsedSeconds * steeringRestoreSpeed)
			if vehicle.SteeringAngle > 0.0 {
				vehicle.SteeringAngle = 0.0
			}
		}
	case s.isSteerLeft:
		vehicle.SteeringAngle += dprec.Degrees(elapsedSeconds * steeringSpeed)
		if vehicle.SteeringAngle > autoMaxSteeringAngle {
			vehicle.SteeringAngle = autoMaxSteeringAngle
		}
	case s.isSteerRight:
		vehicle.SteeringAngle -= dprec.Degrees(elapsedSeconds * steeringSpeed)
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

func (s *VehicleSystem) updateVehiclePhysics(vehicle *ecscomp.Vehicle, elapsedSeconds float64) {
	if vehicle.Recover {
		vehicle.Chassis.Body.SetAngularVelocity(dprec.Vec3Sum(vehicle.Chassis.Body.AngularVelocity(), dprec.NewVec3(0.0, 0.0, 0.1)))
		vehicle.Chassis.Body.SetVelocity(dprec.Vec3Sum(vehicle.Chassis.Body.Velocity(), dprec.NewVec3(0.0, 0.2, 0.0)))
	}

	steeringQuat := dprec.RotationQuat(vehicle.SteeringAngle, dprec.BasisYVec3())
	isMovingForward := dprec.Vec3Dot(vehicle.Chassis.Body.Velocity(), vehicle.Chassis.Body.Orientation().OrientationZ()) > 5.0
	isMovingBackward := dprec.Vec3Dot(vehicle.Chassis.Body.Velocity(), vehicle.Chassis.Body.Orientation().OrientationZ()) < -5.0

	for _, wheel := range vehicle.Wheels {
		if wheel.RotationConstraint != nil {
			wheel.RotationConstraint.SetPrimaryAxis(dprec.QuatVec3Rotation(steeringQuat, dprec.BasisXVec3()))
		}

		if vehicle.Acceleration > 0.0 {
			if isMovingBackward {
				if wheelVelocity := dprec.Vec3Dot(wheel.Body.AngularVelocity(), wheel.Body.Orientation().OrientationX()); wheelVelocity < 0.0 {
					correction := dprec.Max(-vehicle.Acceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.SetAngularVelocity(dprec.Vec3Prod(wheel.Body.AngularVelocity(), 1.0-correction/wheelVelocity))
				}
			} else {
				wheel.Body.SetAngularVelocity(dprec.Vec3Sum(wheel.Body.AngularVelocity(),
					dprec.Vec3Prod(wheel.Body.Orientation().OrientationX(), vehicle.Acceleration*wheel.AccelerationVelocity*elapsedSeconds),
				))
			}
		}

		if vehicle.Deceleration > 0.0 {
			if isMovingForward {
				if wheelVelocity := dprec.Vec3Dot(wheel.Body.AngularVelocity(), wheel.Body.Orientation().OrientationX()); wheelVelocity > 0.0 {
					correction := dprec.Min(vehicle.Deceleration*wheel.DecelerationVelocity*elapsedSeconds, wheelVelocity)
					wheel.Body.SetAngularVelocity(dprec.Vec3Prod(wheel.Body.AngularVelocity(), 1.0-correction/wheelVelocity))
				}
			} else {
				wheel.Body.SetAngularVelocity(dprec.Vec3Sum(wheel.Body.AngularVelocity(),
					dprec.Vec3Prod(wheel.Body.Orientation().OrientationX(), -vehicle.Deceleration*wheel.AccelerationVelocity*elapsedSeconds),
				))
			}
		}
	}
}
