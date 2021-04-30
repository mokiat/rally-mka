package loading

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/graphics"
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

	program   *resource.Program
	texture   *resource.TwoDTexture
	quadModel *resource.Model

	indicatorModelMatrix sprec.Mat4
}

func (v *View) Load(window app.Window, cb func()) {
	loadOutcome := async.NewCompositeOutcome(
		v.registry.LoadProgram("forward-albedo").OnSuccess(resource.InjectProgram(&v.program)),
		v.registry.LoadTwoDTexture("loading").OnSuccess(resource.InjectTwoDTexture(&v.texture)),
		v.registry.LoadModel("quad").OnSuccess(resource.InjectModel(&v.quadModel)),
	)
	go func() {
		loadOutcome.Wait()
		window.Schedule(func() error {
			cb()
			return nil
		})
	}()
}

func (v *View) Unload(window app.Window) {
	async.NewCompositeOutcome(
		v.registry.UnloadProgram(v.program),
		v.registry.UnloadTwoDTexture(v.texture),
		v.registry.UnloadModel(v.quadModel),
	)
}

func (v *View) Open(window app.Window) {}

func (v *View) Close(window app.Window) {}

func (v *View) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	return false
}

func (v *View) Update(window app.Window, elapsedTime time.Duration) {
	elapsedSeconds := float32(elapsedTime.Seconds())
	v.indicatorModelMatrix = sprec.Mat4Prod(v.indicatorModelMatrix,
		sprec.RotationMat4(sprec.Degrees(elapsedSeconds*360.0), 0.0, 0.0, -1.0),
	)
}

func (v *View) Render(window app.Window, width, height int, pipeline *graphics.Pipeline) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	projectionMatrix := sprec.OrthoMat4(
		-screenHalfWidth, screenHalfWidth, screenHalfHeight, -screenHalfHeight, -1.0, 1.0,
	)

	sequence := pipeline.BeginSequence()
	sequence.TargetFramebuffer = &graphics.Framebuffer{
		Width:  int32(width),
		Height: int32(height),
	}
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

	quadMesh := v.quadModel.Nodes[0].Mesh
	quadSubMesh := quadMesh.SubMeshes[0]
	indicatorItem.VertexArray = quadMesh.GFXVertexArray
	indicatorItem.IndexOffset = quadSubMesh.IndexOffset
	indicatorItem.IndexCount = quadSubMesh.IndexCount
	sequence.EndItem(indicatorItem)

	pipeline.EndSequence(sequence)
}
