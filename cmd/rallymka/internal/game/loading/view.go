package loading

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewView(registry *resource.Registry) *View {
	return &View{
		registry: registry,

		screenFramebuffer: &graphics.Framebuffer{},
		program:           stream.GetProgram(registry, "diffuse"),
		texture:           stream.GetTwoDTexture(registry, "loading"),
		mesh:              stream.GetMesh(registry, "quad"),

		indicatorModelMatrix: sprec.IdentityMat4(),
	}
}

type View struct {
	registry *resource.Registry

	screenFramebuffer *graphics.Framebuffer
	program           stream.ProgramHandle
	texture           stream.TwoDTextureHandle
	mesh              stream.MeshHandle

	indicatorModelMatrix sprec.Mat4
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

func (v *View) Update(ctx game.UpdateContext) {
	elapsedSeconds := float32(ctx.ElapsedTime.Seconds())
	v.indicatorModelMatrix = sprec.Mat4Prod(v.indicatorModelMatrix,
		sprec.RotationMat4(sprec.Degrees(elapsedSeconds*360.0), 0.0, 0.0, -1.0),
	)
}

func (v *View) Render(ctx game.RenderContext) {
	v.screenFramebuffer.Width = int32(ctx.WindowSize.Width)
	v.screenFramebuffer.Height = int32(ctx.WindowSize.Height)

	screenHalfWidth := float32(ctx.WindowSize.Width) / float32(ctx.WindowSize.Height)
	screenHalfHeight := float32(1.0)
	projectionMatrix := sprec.OrthoMat4(
		-screenHalfWidth, screenHalfWidth, screenHalfHeight, -screenHalfHeight, -1.0, 1.0,
	)

	sequence := ctx.GFXPipeline.BeginSequence()
	sequence.TargetFramebuffer = v.screenFramebuffer
	sequence.BackgroundColor = sprec.NewVec4(0.0, 0.0, 0.0, 1.0)
	sequence.ClearColor = true
	sequence.ClearDepth = true
	sequence.WriteDepth = false
	sequence.ProjectionMatrix = projectionMatrix
	sequence.ViewMatrix = sprec.ScaleMat4(0.1, 0.1, 1.0)

	indicatorItem := sequence.BeginItem()
	indicatorItem.Program = v.program.Get()
	indicatorItem.ModelMatrix = v.indicatorModelMatrix
	indicatorItem.DiffuseTexture = v.texture.Get()
	indicatorItem.VertexArray = v.mesh.Get().VertexArray
	indicatorItem.IndexCount = v.mesh.Get().SubMeshes[0].IndexCount
	sequence.EndItem(indicatorItem)

	ctx.GFXPipeline.EndSequence(sequence)
}
