package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
)

type Entity struct {
	Transform    *TransformComponent
	RenderMesh   *RenderMesh
	RenderModel  *RenderModel
	RenderSkybox *RenderSkybox
	Vehicle      *Vehicle
	Wheel        *Wheel
	Input        *Input
	CameraStand  *CameraStand
}

type RenderModel struct {
	Model       *stream.Model
	GeomProgram *graphics.Program
}

type RenderMesh struct {
	Mesh        *stream.Mesh
	GeomProgram *graphics.Program
}

type RenderSkybox struct {
	Program *graphics.Program
	Texture *graphics.CubeTexture
	Mesh    *stream.Mesh
}

type Vehicle struct {
	SteeringAngle   sprec.Angle
	Acceleration    float32
	HandbrakePulled bool

	Position        sprec.Vec3
	Orientation     Orientation
	Velocity        sprec.Vec3
	AngularVelocity sprec.Vec3

	FLWheel *Entity
	FRWheel *Entity
	BLWheel *Entity
	BRWheel *Entity
}

type Wheel struct {
	SteeringAngle sprec.Angle
	RotationAngle sprec.Angle

	AnchorPosition sprec.Vec3

	IsDriven bool

	Length float32
	Radius float32

	SuspensionLength float32
	IsGrounded       bool
}

type Input struct{}

type CameraStand struct {
	Target         *Entity
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *Camera
}
