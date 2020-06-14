package loading

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/resource"
)

func NewView(registry *resource.Registry) *View {
	return &View{
		registry:             registry,
		indicatorModelMatrix: sprec.IdentityMat4(),
	}
}

type View struct {
	registry *resource.Registry

	loadOutcome async.Outcome
	program     *resource.Program
	texture     *resource.TwoDTexture
	mesh        *resource.Mesh

	indicatorModelMatrix sprec.Mat4
}

func (v *View) Load() {
	v.loadOutcome = async.NewCompositeOutcome(
		v.registry.LoadProgram("forward-albedo").OnSuccess(resource.InjectProgram(&v.program)),
		v.registry.LoadTwoDTexture("loading").OnSuccess(resource.InjectTwoDTexture(&v.texture)),
		v.registry.LoadMesh("quad").OnSuccess(resource.InjectMesh(&v.mesh)),
	)
}

func (v *View) Unload() {
	async.NewCompositeOutcome(
		v.registry.UnloadProgram(v.program),
		v.registry.UnloadTwoDTexture(v.texture),
		v.registry.UnloadMesh(v.mesh),
	)
}

func (v *View) IsAvailable() bool {
	if v.loadOutcome.IsAvailable() {
		if err := v.loadOutcome.Wait().Err; err != nil {
			panic(err)
		}
		return true
	}
	return false
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
	screenHalfWidth := float32(ctx.WindowSize.Width) / float32(ctx.WindowSize.Height)
	screenHalfHeight := float32(1.0)
	projectionMatrix := sprec.OrthoMat4(
		-screenHalfWidth, screenHalfWidth, screenHalfHeight, -screenHalfHeight, -1.0, 1.0,
	)

	sequence := ctx.GFXPipeline.BeginSequence()
	sequence.BackgroundColor = sprec.NewVec4(0.0, 0.0, 0.0, 1.0)
	sequence.ClearColor = true
	sequence.ClearDepth = true
	sequence.WriteDepth = false
	sequence.TestDepth = false
	sequence.ProjectionMatrix = projectionMatrix
	sequence.ViewMatrix = sprec.ScaleMat4(0.1, 0.1, 1.0)

	indicatorItem := sequence.BeginItem()
	indicatorItem.Program = v.program.GFXProgram
	indicatorItem.ModelMatrix = v.indicatorModelMatrix
	indicatorItem.AlbedoTwoDTexture = v.texture.GFXTexture
	indicatorItem.VertexArray = v.mesh.GFXVertexArray
	indicatorItem.IndexOffset = v.mesh.SubMeshes[0].IndexOffset
	indicatorItem.IndexCount = v.mesh.SubMeshes[0].IndexCount
	sequence.EndItem(indicatorItem)

	ctx.GFXPipeline.EndSequence(sequence)
}
