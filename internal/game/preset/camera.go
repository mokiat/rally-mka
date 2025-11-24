package preset

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/ui"
)

type GamepadProvider interface {
	Gamepads() [4]app.Gamepad
}

func NewFollowCameraSystem(ecsScene *ECSScene, gamepadProvider GamepadProvider) *FollowCameraSystem {
	return &FollowCameraSystem{
		ecsScene:        ecsScene,
		gamepadProvider: gamepadProvider,
	}
}

type FollowCameraSystem struct {
	ecsScene        *ECSScene
	gamepadProvider GamepadProvider

	rotationSpeed    dprec.Angle
	rotationStrength dprec.Angle
	zoomSpeed        float64

	hasKeyboardConsumer bool

	keyRotateLeft  ui.KeyCode
	keyRotateRight ui.KeyCode
	keyRotateUp    ui.KeyCode
	keyRotateDown  ui.KeyCode
	keyZoomIn      ui.KeyCode
	keyZoomOut     ui.KeyCode

	isRotateLeft  bool
	isRotateRight bool
	isRotateUp    bool
	isRotateDown  bool
	isZoomIn      bool
	isZoomOut     bool
}

func (s *FollowCameraSystem) UseDefaults() {
	s.rotationSpeed = dprec.Degrees(100)
	s.rotationStrength = dprec.Degrees(300.0)
	s.zoomSpeed = 1.0

	s.keyRotateUp = ui.KeyCodeW
	s.keyRotateDown = ui.KeyCodeS
	s.keyRotateLeft = ui.KeyCodeA
	s.keyRotateRight = ui.KeyCodeD
	s.keyZoomIn = ui.KeyCodeE
	s.keyZoomOut = ui.KeyCodeQ
}

func (s *FollowCameraSystem) OnKeyboardEvent(event ui.KeyboardEvent) bool {
	if !s.hasKeyboardConsumer {
		return false
	}
	active := event.Action != ui.KeyboardActionUp
	switch event.Code {
	case s.keyRotateUp:
		s.isRotateUp = active
		return true
	case s.keyRotateDown:
		s.isRotateDown = active
		return true
	case s.keyRotateLeft:
		s.isRotateLeft = active
		return true
	case s.keyRotateRight:
		s.isRotateRight = active
		return true
	case s.keyZoomIn:
		s.isZoomIn = active
		return true
	case s.keyZoomOut:
		s.isZoomOut = active
		return true
	default:
		return false
	}
}

func (s *FollowCameraSystem) OnFixedUpdate(elapsedSeconds float64) {
	result := s.ecsScene.Query(
		s.ecsScene.HasNodeComponent(),
		s.ecsScene.HasFollowCameraComponent(),
	)
	defer result.Release()

	var hasKeyboardConsumer bool
	result.Each(func(entity Entity) {
		controlled := entity.RefControlledComponent()
		if controlled != nil {
			if controlled.Inputs.Is(ControlInputKeyboard) {
				hasKeyboardConsumer = true
				s.updateKeyboard(elapsedSeconds, entity)
			}
			if controlled.Inputs.Is(ControlInputGamepad0) {
				if gamepad := s.gamepadProvider.Gamepads()[0]; gamepad.Connected() && gamepad.Supported() {
					s.updateGamepad(elapsedSeconds, gamepad, entity)
				}
			}
		}
		s.updateCamera(elapsedSeconds, entity)
	})
	s.hasKeyboardConsumer = hasKeyboardConsumer
}

func (s *FollowCameraSystem) updateKeyboard(elapsedSeconds float64, entity Entity) {
	cameraComp := entity.RefFollowCameraComponent()

	if s.isRotateUp {
		cameraComp.PitchAngle -= s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}
	if s.isRotateDown {
		cameraComp.PitchAngle += s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}

	if s.isRotateLeft {
		cameraComp.YawAngle -= s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}
	if s.isRotateRight {
		cameraComp.YawAngle += s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}

	if s.isZoomIn {
		cameraComp.Zoom -= cameraComp.Zoom * elapsedSeconds
	}
	if s.isZoomOut {
		cameraComp.Zoom += cameraComp.Zoom * elapsedSeconds
	}
}

func (s *FollowCameraSystem) updateGamepad(elapsedSeconds float64, gamepad app.Gamepad, entity Entity) {
	cameraComp := entity.RefFollowCameraComponent()

	if gamepad.DpadUpButton() {
		cameraComp.PitchAngle -= s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}
	if gamepad.DpadDownButton() {
		cameraComp.PitchAngle += s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}

	if gamepad.DpadLeftButton() {
		cameraComp.YawAngle -= s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}
	if gamepad.DpadRightButton() {
		cameraComp.YawAngle += s.rotationSpeed * dprec.Angle(elapsedSeconds)
	}

	if gamepad.RightBumper() {
		cameraComp.Zoom -= cameraComp.Zoom * elapsedSeconds
	}
	if gamepad.LeftBumper() {
		cameraComp.Zoom += cameraComp.Zoom * elapsedSeconds
	}

	if dprec.Abs(gamepad.RightStickX()) > 0.0 || dprec.Abs(gamepad.RightStickY()) > 0.0 {
		target := cameraComp.Target
		targetPosition := target.AbsoluteMatrix().Translation()
		anchorVector := dprec.Vec3Diff(cameraComp.AnchorPosition, targetPosition)

		cameraVectorZ := anchorVector
		cameraVectorX := dprec.Vec3Cross(dprec.BasisYVec3(), cameraVectorZ)
		cameraVectorY := dprec.Vec3Cross(cameraVectorZ, cameraVectorX)

		angle := -s.rotationStrength * dprec.Angle(gamepad.RightStickY()*elapsedSeconds)
		anchorVector = dprec.QuatVec3Rotation(dprec.RotationQuat(angle, cameraVectorX), anchorVector)

		angle = -s.rotationStrength * dprec.Angle(gamepad.RightStickX()*elapsedSeconds)
		anchorVector = dprec.QuatVec3Rotation(dprec.RotationQuat(angle, cameraVectorY), anchorVector)

		cameraComp.AnchorPosition = dprec.Vec3Sum(targetPosition, anchorVector)
	}
}

func (s *FollowCameraSystem) updateCamera(_ float64, entity Entity) {
	nodeComp := entity.RefNodeComponent()
	cameraComp := entity.RefFollowCameraComponent()

	target := cameraComp.Target
	targetPosition := target.AbsoluteMatrix().Translation()

	// We use a camera anchor to achieve the smooth effect of a
	// camera following the target.
	anchorVector := dprec.Vec3Diff(cameraComp.AnchorPosition, targetPosition)
	anchorVector = dprec.ResizedVec3(anchorVector, cameraComp.AnchorDistance)
	cameraComp.AnchorPosition = dprec.Vec3Sum(targetPosition, anchorVector)

	cameraVectorZ := anchorVector
	cameraVectorX := dprec.Vec3Cross(dprec.BasisYVec3(), cameraVectorZ)
	cameraVectorY := dprec.Vec3Cross(cameraVectorZ, cameraVectorX)

	matrix := dprec.Mat4MultiProd(
		dprec.TransformationMat4(
			dprec.UnitVec3(cameraVectorX),
			dprec.UnitVec3(cameraVectorY),
			dprec.UnitVec3(cameraVectorZ),
			targetPosition,
		),
		dprec.RotationMat4(cameraComp.YawAngle, 0.0, 1.0, 0.0),
		dprec.RotationMat4(cameraComp.PitchAngle, 1.0, 0.0, 0.0),
		dprec.TranslationMat4(0.0, 0.0, cameraComp.CameraDistance*cameraComp.Zoom),
	)
	nodeComp.Node.SetAbsoluteMatrix(matrix)
}
