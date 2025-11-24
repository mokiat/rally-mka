package preset

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game/hierarchy"
	"github.com/mokiat/lacking/ui"
)

type NodeComponent struct {
	Node hierarchy.Node
}

type ControlledComponent struct {
	Inputs ControlInput
}

type FollowCameraComponent struct {
	Target         hierarchy.Node
	AnchorPosition dprec.Vec3
	AnchorDistance float64
	CameraDistance float64
	PitchAngle     dprec.Angle
	YawAngle       dprec.Angle
	Zoom           float64
}

type CarComponent struct {
	Car            *Car
	Gear           CarGear
	SteeringAmount float64
	Acceleration   float64
	Deceleration   float64
	Recover        bool
	LightsOn       bool
}

type CarKeyboardComponent struct {
	AccelerateKey ui.KeyCode
	DecelerateKey ui.KeyCode
	TurnLeftKey   ui.KeyCode
	TurnRightKey  ui.KeyCode
	ShiftUpKey    ui.KeyCode
	ShiftDownKey  ui.KeyCode
	RecoverKey    ui.KeyCode

	AccelerationChangeSpeed float64
	DecelerationChangeSpeed float64
	SteeringAmount          float64
	SteeringChangeSpeed     float64
	SteeringRestoreSpeed    float64
}

type CarMouseComponent struct {
	AccelerationChangeSpeed float64
	DecelerationChangeSpeed float64
	Destination             dprec.Vec3
}

type CarGamepadComponent struct {
	Gamepad app.Gamepad
}

const (
	ControlInputKeyboard ControlInput = 1 << iota
	ControlInputMouse
	ControlInputGamepad0
)

type ControlInput int

func (i ControlInput) Is(query ControlInput) bool {
	return i&query != 0
}

const (
	CarGearNeutral CarGear = iota
	CarGearForward
	CarGearReverse
)

type CarGear int
