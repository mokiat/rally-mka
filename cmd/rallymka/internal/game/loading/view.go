package loading

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry) *View {
	return &View{
		registry: registry,

		program: stream.GetProgram(registry, "diffuse"),
		texture: stream.GetTwoDTexture(registry, "loading"),
		mesh:    stream.GetMesh(registry, "quad"),

		projectionMatrix:     math.IdentityMat4x4(),
		indicatorModelMatrix: math.IdentityMat4x4(),
	}
}

type View struct {
	registry *resource.Registry

	program stream.ProgramHandle
	texture stream.TwoDTextureHandle
	mesh    stream.MeshHandle

	projectionMatrix     math.Mat4x4
	indicatorModelMatrix math.Mat4x4
}

func (v *View) Load() {
	v.registry.Request(v.program.Handle)
	v.registry.Request(v.texture.Handle)
	v.registry.Request(v.mesh.Handle)
}

func (v *View) Unload() {
	v.registry.Dismiss(v.program.Handle)
	v.registry.Dismiss(v.texture.Handle)
	v.registry.Dismiss(v.mesh.Handle)
}

func (v *View) IsAvailable() bool {
	return v.program.IsAvailable() && v.texture.IsAvailable() && v.mesh.IsAvailable()
}

func (v *View) Open() {}

func (v *View) Close() {}

func (v *View) Resize(width, height int) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	v.projectionMatrix = math.OrthoMat4x4(
		-screenHalfWidth, screenHalfWidth, screenHalfHeight, -screenHalfHeight, -1.0, 1.0,
	)
}

func (v *View) Update(elapsedSeconds float32, actions input.ActionSet) {
	v.indicatorModelMatrix = v.indicatorModelMatrix.MulMat4x4(
		math.RotationMat4x4(elapsedSeconds*360.0, 0.0, 0.0, -1.0),
	)
}

func (v *View) Render(pipeline *graphics.Pipeline) {
	sequence := pipeline.BeginSequence()
	sequence.BackgroundColor = math.MakeVec4(0.0, 0.0, 0.0, 1.0)
	sequence.ClearColor = true
	sequence.ClearDepth = true
	sequence.WriteDepth = false
	sequence.ProjectionMatrix = v.projectionMatrix
	sequence.ViewMatrix = math.ScaleMat4x4(0.1, 0.1, 1.0)

	indicatorItem := sequence.BeginItem()
	indicatorItem.Program = v.program.Get()
	indicatorItem.ModelMatrix = v.indicatorModelMatrix
	indicatorItem.DiffuseTexture = v.texture.Get()
	indicatorItem.VertexArray = v.mesh.Get().VertexArray
	indicatorItem.IndexCount = v.mesh.Get().SubMeshes[0].IndexCount
	sequence.EndItem(indicatorItem)

	pipeline.EndSequence(sequence)
}
