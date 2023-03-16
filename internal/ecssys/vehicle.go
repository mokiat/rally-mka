package ecssys

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/preset"
	"github.com/mokiat/lacking/ui"
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

func (s *VehicleSystem) OnMouseEvent(event ui.MouseEvent, viewport graphics.Viewport, camera *graphics.Camera, graphicsScene *graphics.Scene) bool {
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
		s.mouseX = event.Position.X
		s.mouseY = event.Position.Y
	}
	return true
}

func (s *VehicleSystem) Update(elapsedSeconds float64) {
	result := s.ecsScene.Find(ecs.Having(ecscomp.VehicleComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()
		vehicle := ecscomp.GetVehicle(entity)
		var controlled *preset.ControlledComponent
		if ecs.FetchComponent(entity, &controlled) {
			if s.mouseActive {
				s.updateVehicleControlMouse(vehicle, elapsedSeconds)
			} else {
				s.updateVehicleControlKeyboard(vehicle, elapsedSeconds)
			}
		}
	}
}

func (s *VehicleSystem) updateVehicleControlMouse(vehicle *ecscomp.Vehicle, elapsedSeconds float64) {
	vehicle.Recover = s.isRecover

	orientation := vehicle.Chassis.Body.Orientation()
	orientationX := orientation.OrientationX()
	orientationZ := orientation.OrientationZ()

	// TODO: Use sphere shape instead
	// TODO: Move Position and orientation into placement
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

	surface := shape.NewPlacement(shape.IdentityTransform(), shape.NewStaticMesh([]shape.StaticTriangle{
		shape.NewStaticTriangle(shape.Point(a), shape.Point(b), shape.Point(c)),
		shape.NewStaticTriangle(shape.Point(a), shape.Point(c), shape.Point(d)),
	}))

	line := s.graphicsScene.Ray(s.viewport, s.camera, s.mouseX, s.mouseY)

	result := shape.NewIntersectionResultSet(1)
	shape.CheckIntersectionLineWithMesh(line, surface, result)

	if result.Found() {
		intersection := result.Intersections()[0]
		mouseTarget := intersection.FirstContact

		delta := dprec.Vec3Diff(mouseTarget, vehicle.Chassis.Body.Position())
		delta.Y = 0.0

		forward := vehicle.Chassis.Body.Orientation().OrientationZ()
		forward.Y = 0.0

		sin := dprec.Vec3Cross(
			dprec.UnitVec3(delta),
			dprec.UnitVec3(forward),
		)
		angle := dprec.Angle(dprec.Sign(sin.Y)) * dprec.Asin(sin.Length())

		vehicle.SteeringAngle = dprec.Clamp(-dprec.Angle(angle), -vehicle.MaxSteeringAngle, vehicle.MaxSteeringAngle)
	} else {
		vehicle.SteeringAngle = 0
	}

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
