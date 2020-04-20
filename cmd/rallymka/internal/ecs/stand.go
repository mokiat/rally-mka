package ecs

import "github.com/mokiat/gomath/sprec"

func NewCameraStandSystem(ecsManager *Manager) *CameraStandSystem {
	return &CameraStandSystem{
		ecsManager: ecsManager,
	}
}

type CameraStandSystem struct {
	ecsManager *Manager
}

func (s *CameraStandSystem) Update() {
	for _, entity := range s.ecsManager.Entities() {
		if cameraStand := entity.CameraStand; cameraStand != nil {
			s.updateCameraStand(cameraStand)
		}
	}
}

// var angleSpeed sprec.Angle = sprec.Degrees(0.5)
// var angle sprec.Angle = sprec.Degrees(200)

var angleSpeed sprec.Angle = sprec.Degrees(0.0)
var angle sprec.Angle = sprec.Degrees(0)

// var angle sprec.Angle = sprec.Degrees(300)

func (s *CameraStandSystem) updateCameraStand(cameraStand *CameraStand) {
	angle += angleSpeed

	var targetPosition sprec.Vec3
	if cameraStand.Target.Transform != nil {
		targetPosition = cameraStand.Target.Transform.Position
	}
	if cameraStand.Target.Vehicle != nil {
		targetPosition = cameraStand.Target.Vehicle.Position
	}
	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := sprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = sprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)
	cameraStand.AnchorPosition = sprec.Vec3Sum(targetPosition, anchorVector)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ := anchorVector
	cameraVectorX := sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	cameraStand.Camera.SetViewMatrix(sprec.Mat4MultiProd(
		sprec.TranslationMat4(targetPosition.X, targetPosition.Y, targetPosition.Z),
		sprec.TransformationMat4(
			sprec.UnitVec3(cameraVectorX),
			sprec.UnitVec3(cameraVectorY),
			sprec.UnitVec3(cameraVectorZ),
			sprec.ZeroVec3(),
		),
		sprec.RotationMat4(angle, 0.0, 1.0, 0.0),
		sprec.RotationMat4(sprec.Degrees(-25.0), 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, cameraStand.CameraDistance),
	))
}
