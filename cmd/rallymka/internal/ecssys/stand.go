package ecssys

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
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
		targetPhysicsComp = ecscomp.GetPhysics(cameraStand.Target)
		targetRenderComp  = ecscomp.GetRender(cameraStand.Target)
	)

	var targetPosition sprec.Vec3
	switch {
	case targetPhysicsComp != nil:
		targetPosition = targetPhysicsComp.Body.Position()
	case targetRenderComp != nil:
		targetPosition = targetRenderComp.Mesh.Position()
	}
	// we use a camera anchor to achieve the smooth effect of a
	// camera following the target
	anchorVector := sprec.Vec3Diff(cameraStand.AnchorPosition, targetPosition)
	anchorVector = sprec.ResizedVec3(anchorVector, cameraStand.AnchorDistance)

	cameraVectorZ := anchorVector
	cameraVectorX := sprec.Vec3Cross(sprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := sprec.Vec3Cross(cameraVectorZ, cameraVectorX)

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
		sprec.RotationMat4(sprec.Degrees(-25.0), 1.0, 0.0, 0.0),
		sprec.TranslationMat4(0.0, 0.0, cameraStand.CameraDistance*s.zoom),
	)

	cameraStand.Camera.SetPosition(matrix.Translation())
	cameraStand.Camera.SetRotation(matrixToQuat(matrix))
}

// TODO: Move to gomath library.
// This is calculated by inversing the formulas for
// quat.OrientationX, quat.OrientationY and quat.OrientationZ.
func matrixToQuat(matrix sprec.Mat4) sprec.Quat {
	sqrX := (1.0 + matrix.M11 - matrix.M22 - matrix.M33) / 4.0
	sqrY := (1.0 - matrix.M11 + matrix.M22 - matrix.M33) / 4.0
	sqrZ := (1.0 - matrix.M11 - matrix.M22 + matrix.M33) / 4.0

	var x, y, z, w float32
	if sqrZ > sqrX && sqrZ > sqrY {
		// Z is largest
		z = sprec.Sqrt(sqrZ)
		x = (matrix.M31 + matrix.M13) / (4 * z)
		y = (matrix.M32 + matrix.M23) / (4 * z)
		w = (matrix.M21 - matrix.M12) / (4 * z)
	} else if sqrY > sqrX {
		// Y is largest
		y = sprec.Sqrt(sqrY)
		x = (matrix.M21 + matrix.M12) / (4 * y)
		z = (matrix.M32 + matrix.M23) / (4 * y)
		w = (matrix.M13 - matrix.M31) / (4 * y)
	} else {
		// X is largest
		x = sprec.Sqrt(sqrX)
		y = (matrix.M21 + matrix.M12) / (4 * x)
		z = (matrix.M31 + matrix.M13) / (4 * x)
		w = (matrix.M32 - matrix.M23) / (4 * x)
	}
	return sprec.UnitQuat(sprec.NewQuat(w, x, y, z))
}
