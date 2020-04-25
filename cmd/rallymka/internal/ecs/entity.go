package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/physics"
)

type Entity struct {
	Physics      *PhysicsComponent
	Render       *RenderComponent
	RenderSkybox *RenderSkybox

	Input       *Input
	CameraStand *CameraStand
	Car         *Car
}

type PhysicsComponent struct {
	Body *physics.Body
}

type RenderComponent struct {
	GeomProgram *graphics.Program
	Model       *stream.Model
	Mesh        *stream.Mesh
	Matrix      sprec.Mat4
}

type RenderSkybox struct {
	Program *graphics.Program
	Texture *graphics.CubeTexture
	Mesh    *stream.Mesh
}

type CarInput struct {
	Forward   bool
	Backward  bool
	TurnLeft  bool
	TurnRight bool
	Handbrake bool
}

type Car struct {
	SteeringAngle   sprec.Angle
	Acceleration    float32
	HandbrakePulled bool

	Body            *Entity
	FLWheelRotation physics.Constraint
	FLWheel         *Entity
	FRWheelRotation physics.Constraint
	FRWheel         *Entity
	BLWheel         *Entity
	BRWheel         *Entity
}

type Input struct{}

type CameraStand struct {
	Target         *Entity
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *Camera
}
