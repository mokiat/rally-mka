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
	Car          *Car
	CameraStand  *CameraStand
	HumanInput   bool
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
	Chassis         *physics.Body
	FLWheelRotation *physics.MatchAxisConstraint
	FRWheelRotation *physics.MatchAxisConstraint
	FLWheel         *physics.Body
	FRWheel         *physics.Body
	BLWheel         *physics.Body
	BRWheel         *physics.Body
}

type CameraStand struct {
	Target         *Entity
	AnchorPosition sprec.Vec3
	AnchorDistance float32
	CameraDistance float32
	Camera         *Camera
}
