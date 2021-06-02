package ecs

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
)

func NewCameraStandSystem(ecsManager *Manager) *CameraStandSystem {
	return &CameraStandSystem{
		ecsManager: ecsManager,
	}
}

type CameraStandSystem struct {
	ecsManager *Manager
}

func (s *CameraStandSystem) Update(elapsedTime time.Duration, gamepad *app.GamepadState) {
	for _, entity := range s.ecsManager.Entities() {
		if cameraStand := entity.CameraStand; cameraStand != nil {
			s.updateCameraStand(cameraStand, elapsedTime, gamepad)
		}
	}
}

func (s *CameraStandSystem) updateCameraStand(cameraStand *CameraStand, elapsedTime time.Duration, gamepad *app.GamepadState) {
	var targetPosition sprec.Vec3
	switch {
	case cameraStand.Target.Physics != nil:
		targetPosition = cameraStand.Target.Physics.Body.Position()
	case cameraStand.Target.Render != nil:
		targetPosition = cameraStand.Target.Render.Renderable.Matrix.Translation()
	}
	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := sprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = sprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)

	cameraVectorZ := anchorVector
	cameraVectorX := sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	if gamepad != nil {
		elapsedSeconds := float32(elapsedTime.Seconds())
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

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ = anchorVector
	cameraVectorX = sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY = sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	cameraStand.Camera.SetMatrix(sprec.Mat4MultiProd(
		sprec.TranslationMat4(targetPosition.X, targetPosition.Y, targetPosition.Z),
		sprec.TransformationMat4(
			sprec.UnitVec3(cameraVectorX),
			sprec.UnitVec3(cameraVectorY),
			sprec.UnitVec3(cameraVectorZ),
			sprec.ZeroVec3(),
		),
		sprec.RotationMat4(sprec.Degrees(-25.0), 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, cameraStand.CameraDistance),
	))
}
