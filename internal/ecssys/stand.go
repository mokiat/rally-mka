package ecssys

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/rally-mka/internal/ecscomp"
)

func NewCameraStandSystem(ecsScene *ecs.Scene) *CameraStandSystem {
	return &CameraStandSystem{
		ecsScene: ecsScene,
		zoom:     1.0,
	}
}

type CameraStandSystem struct {
	ecsScene *ecs.Scene
	zoom     float64

	isRotateLeft  bool
	isRotateRight bool
	isRotateUp    bool
	isRotateDown  bool
	isZoomIn      bool
	isZoomOut     bool

	rotationX dprec.Angle
	rotationY dprec.Angle
}

func (s *CameraStandSystem) OnKeyboardEvent(event app.KeyboardEvent) bool {
	active := event.Type != app.KeyboardEventTypeKeyUp
	switch event.Code {
	case app.KeyCodeA:
		s.isRotateLeft = active
		return true
	case app.KeyCodeD:
		s.isRotateRight = active
		return true
	case app.KeyCodeW:
		s.isRotateUp = active
		return true
	case app.KeyCodeS:
		s.isRotateDown = active
		return true
	case app.KeyCodeE:
		s.isZoomIn = active
		return true
	case app.KeyCodeQ:
		s.isZoomOut = active
		return true
	}
	return false
}

func (s *CameraStandSystem) Update(elapsedSeconds float64, gamepad *app.GamepadState) {
	result := s.ecsScene.Find(ecs.Having(ecscomp.CameraStandComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()
		cameraStand := ecscomp.GetCameraStand(entity)
		s.updateCameraStand(cameraStand, elapsedSeconds, gamepad)
	}
}

func (s *CameraStandSystem) updateCameraStand(cameraStand *ecscomp.CameraStand, elapsedSeconds float64, gamepad *app.GamepadState) {
	var (
		target = cameraStand.Target
	)
	var targetPosition dprec.Vec3
	switch {
	case target.Body() != nil:
		targetPosition = target.Body().Position()
	default:
		targetPosition = target.Position()
	}

	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := dprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = dprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)

	cameraVectorZ := anchorVector
	cameraVectorX := dprec.Vec3Cross(dprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := dprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	if s.isRotateLeft || (gamepad != nil && gamepad.DpadLeftButton) {
		s.rotationY -= dprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateRight || (gamepad != nil && gamepad.DpadRightButton) {
		s.rotationY += dprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateUp || (gamepad != nil && gamepad.DpadUpButton) {
		s.rotationX -= dprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateDown || (gamepad != nil && gamepad.DpadDownButton) {
		s.rotationX += dprec.Degrees(elapsedSeconds * 100)
	}
	if s.isZoomIn {
		s.zoom -= elapsedSeconds * s.zoom
	}
	if s.isZoomOut {
		s.zoom += elapsedSeconds * s.zoom
	}

	if gamepad != nil {
		if gamepad.RightBumper {
			s.zoom = s.zoom - 0.3*elapsedSeconds*s.zoom
		}
		if gamepad.LeftBumper {
			s.zoom = s.zoom + 0.3*elapsedSeconds*s.zoom
		}

		rotationAmount := 200 * elapsedSeconds
		if dprec.Abs(gamepad.RightStickY) > 0.2 {
			rotation := dprec.RotationQuat(dprec.Degrees(gamepad.RightStickY*rotationAmount), cameraVectorX)
			anchorVector = dprec.QuatVec3Rotation(rotation, anchorVector)
		}
		if dprec.Abs(gamepad.RightStickX) > 0.2 {
			rotation := dprec.RotationQuat(dprec.Degrees(-gamepad.RightStickX*rotationAmount), cameraVectorY)
			anchorVector = dprec.QuatVec3Rotation(rotation, anchorVector)
		}
	}

	cameraStand.AnchorPosition = dprec.Vec3Sum(targetPosition, anchorVector)
	// cameraStand.AnchorPosition = dprec.NewVec3(10.0, 60.0, 40.0)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ = anchorVector
	cameraVectorX = dprec.Vec3Cross(dprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY = dprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	matrix := dprec.Mat4MultiProd(
		dprec.TranslationMat4(targetPosition.X, targetPosition.Y, targetPosition.Z),
		dprec.TransformationMat4(
			dprec.UnitVec3(cameraVectorX),
			dprec.UnitVec3(cameraVectorY),
			dprec.UnitVec3(cameraVectorZ),
			dprec.ZeroVec3(),
		),
		dprec.RotationMat4(s.rotationY, 0.0, 1.0, 0.0),
		dprec.RotationMat4(dprec.Degrees(-25.0), 1.0, 0.0, 0.0),
		dprec.RotationMat4(s.rotationX, 1.0, 0.0, 0.0),
		dprec.TranslationMat4(0.0, 0.0, cameraStand.CameraDistance*s.zoom),
	)
	cameraStand.Camera.SetMatrix(matrix)
}
