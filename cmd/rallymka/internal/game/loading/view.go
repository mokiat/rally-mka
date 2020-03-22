package loading

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry) *View {
	program := stream.GetProgram(registry, "diffuse")
	registry.Request(program.Handle)
	texture := stream.GetTwoDTexture(registry, "loading")
	registry.Request(texture.Handle)
	mesh := stream.GetMesh(registry, "quad")
	registry.Request(mesh.Handle)

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

	program stream.ProgramHandle
	texture stream.TwoDTextureHandle
	mesh    stream.MeshHandle
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

	if v.program.IsAvailable() && v.texture.IsAvailable() && v.mesh.IsAvailable() {
		indicatorItem := sequence.BeginItem()
		indicatorItem.Program = v.program.Get()
		indicatorItem.ModelMatrix = v.indicatorModelMatrix
		indicatorItem.DiffuseTexture = v.texture.Get()
		indicatorItem.VertexArray = v.mesh.Get().VertexArray
		indicatorItem.IndexCount = v.mesh.Get().SubMeshes[0].IndexCount
		sequence.EndItem(indicatorItem)
	}
	pipeline.EndSequence(sequence)
}
