package ecssys

import (
	"github.com/mokiat/gomath/sprec"
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
	zoom     float32

	isRotateLeft  bool
	isRotateRight bool
	isRotateUp    bool
	isRotateDown  bool
	isZoomIn      bool
	isZoomOut     bool

	rotationX sprec.Angle
	rotationY sprec.Angle
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

func (s *CameraStandSystem) Update(elapsedSeconds float32, gamepad *app.GamepadState) {
	result := s.ecsScene.Find(ecs.Having(ecscomp.CameraStandComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()
		cameraStand := ecscomp.GetCameraStand(entity)
		s.updateCameraStand(cameraStand, elapsedSeconds, gamepad)
	}
}

func (s *CameraStandSystem) updateCameraStand(cameraStand *ecscomp.CameraStand, elapsedSeconds float32, gamepad *app.GamepadState) {
	var (
		target = cameraStand.Target
	)
	var targetPosition sprec.Vec3
	switch {
	case target.Body() != nil:
		targetPosition = target.Body().Position()
	default:
		targetPosition = target.AbsoluteMatrix().Translation()
	}

	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := sprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = sprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)

	cameraVectorZ := anchorVector
	cameraVectorX := sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	if s.isRotateLeft {
		s.rotationY -= sprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateRight {
		s.rotationY += sprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateUp {
		s.rotationX -= sprec.Degrees(elapsedSeconds * 100)
	}
	if s.isRotateDown {
		s.rotationX += sprec.Degrees(elapsedSeconds * 100)
	}
	if s.isZoomIn {
		s.zoom -= 0.3 * elapsedSeconds * s.zoom
	}
	if s.isZoomOut {
		s.zoom += 0.3 * elapsedSeconds * s.zoom
	}

	if gamepad != nil {
		if gamepad.RightBumper {
			s.zoom = s.zoom - 0.3*elapsedSeconds*s.zoom
		}
		if gamepad.LeftBumper {
			s.zoom = s.zoom + 0.3*elapsedSeconds*s.zoom
		}

		rotationAmount := 200 * elapsedSeconds
		if sprec.Abs(gamepad.RightStickY) > 0.2 {
			rotation := sprec.RotationQuat(sprec.Degrees(gamepad.RightStickY*rotationAmount), cameraVectorX)
			anchorVector = sprec.QuatVec3Rotation(rotation, anchorVector)
		}
		if sprec.Abs(gamepad.RightStickX) > 0.2 {
			rotation := sprec.RotationQuat(sprec.Degrees(-gamepad.RightStickX*rotationAmount), cameraVectorY)
			anchorVector = sprec.QuatVec3Rotation(rotation, anchorVector)
		}
	}

	cameraStand.AnchorPosition = sprec.Vec3Sum(targetPosition, anchorVector)
	// cameraStand.AnchorPosition = sprec.NewVec3(10.0, 60.0, 40.0)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ = anchorVector
	cameraVectorX = sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY = sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	matrix := sprec.Mat4MultiProd(
		sprec.TranslationMat4(targetPosition.X, targetPosition.Y, targetPosition.Z),
		sprec.TransformationMat4(
			sprec.UnitVec3(cameraVectorX),
			sprec.UnitVec3(cameraVectorY),
			sprec.UnitVec3(cameraVectorZ),
			sprec.ZeroVec3(),
		),
		sprec.RotationMat4(s.rotationY, 0.0, 1.0, 0.0),
		sprec.RotationMat4(sprec.Degrees(-25.0), 1.0, 0.0, 0.0),
		sprec.RotationMat4(s.rotationX, 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, cameraStand.CameraDistance*s.zoom),
	)
	cameraStand.Camera.SetMatrix(matrix)
}
