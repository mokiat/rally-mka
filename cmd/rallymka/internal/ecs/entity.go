package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
)

type Entity struct {
	Debug     *DebugComponent
	Transform *TransformComponent
	Motion    *MotionComponent
	Collision *CollisionComponent

	RenderMesh   *RenderMesh
	RenderModel  *RenderModel
	RenderSkybox *RenderSkybox
	Input        *Input
	CameraStand  *CameraStand
	Car          *Car
}

type DebugComponent struct {
	Name string
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
	FLWheelRotation Constraint
	FLWheel         *Entity
	FRWheelRotation Constraint
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
