package loading

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
)

func NewView(registry *resource.Registry) *View {
	program := stream.GetProgram(registry, "diffuse")
	program.Request()
	texture := stream.GetTwoDTexture(registry, "loading")
	texture.Request()
	mesh := stream.GetMesh(registry, "quad")
	mesh.Request()

	return &View{
		projectionMatrix:     math.IdentityMat4x4(),
		indicatorModelMatrix: math.IdentityMat4x4(),

		program: program,
		texture: texture,
		mesh:    mesh,
	}
}

type View struct {
	projectionMatrix     math.Mat4x4
	indicatorModelMatrix math.Mat4x4

	program *stream.Program
	texture *stream.TwoDTexture
	mesh    *stream.Mesh
}

func (v *View) Resize(width, height int) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	// TODO: Switch to Ortho projection
	v.projectionMatrix = math.PerspectiveMat4x4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 300.0,
	)
}

func (v *View) Update(elapsedSeconds float32) {
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
	sequence.ViewMatrix = math.TranslationMat4x4(0.0, 0.0, -15.0)

	if v.program.Available() && v.texture.Available() && v.mesh.Available() {
		indicatorItem := sequence.BeginItem()
		indicatorItem.Program = v.program.Gfx()
		indicatorItem.ModelMatrix = v.indicatorModelMatrix
		indicatorItem.DiffuseTexture = v.texture.Gfx()
		indicatorItem.VertexArray = v.mesh.Gfx()
		indicatorItem.IndexCount = int32(v.mesh.SubMeshes()[0].IndexCount)
		sequence.EndItem(indicatorItem)
	}
	pipeline.EndSequence(sequence)
}
