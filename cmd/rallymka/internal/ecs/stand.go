package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
)

func NewCameraStandSystem(ecsManager *Manager) *CameraStandSystem {
	return &CameraStandSystem{
		ecsManager: ecsManager,
	}
}

type CameraStandSystem struct {
	ecsManager *Manager
}

func (s *CameraStandSystem) Update(ctx game.UpdateContext) {
	for _, entity := range s.ecsManager.Entities() {
		if cameraStand := entity.CameraStand; cameraStand != nil {
			s.updateCameraStand(cameraStand, ctx)
		}
	}
}

func (s *CameraStandSystem) updateCameraStand(cameraStand *CameraStand, ctx game.UpdateContext) {
	var targetPosition sprec.Vec3
	switch {
	case cameraStand.Target.Physics != nil:
		targetPosition = cameraStand.Target.Physics.Body.Position
	case cameraStand.Target.Render != nil:
		targetPosition = cameraStand.Target.Render.Matrix.Translation()
	}
	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := sprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = sprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)

	cameraVectorZ := anchorVector
	cameraVectorX := sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	if ctx.Gamepad.Available {
		if sprec.Abs(ctx.Gamepad.RightStickY) > 0.1 {
			rotation := sprec.RotationQuat(sprec.Degrees(ctx.Gamepad.RightStickY*2), cameraVectorX)
			anchorVector = sprec.QuatVec3Rotation(rotation, anchorVector)
		}
		if sprec.Abs(ctx.Gamepad.RightStickX) > 0.1 {
			rotation := sprec.RotationQuat(sprec.Degrees(-ctx.Gamepad.RightStickX*2), cameraVectorY)
			anchorVector = sprec.QuatVec3Rotation(rotation, anchorVector)
		}
	}

	cameraStand.AnchorPosition = sprec.Vec3Sum(targetPosition, anchorVector)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ = anchorVector
	cameraVectorX = sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY = sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	cameraStand.Camera.SetViewMatrix(sprec.Mat4MultiProd(
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
