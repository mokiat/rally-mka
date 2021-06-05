package ecssys

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
)

func NewRenderer(ecsScene *ecs.Scene, scene *render.Scene) *Renderer {
	return &Renderer{
		ecsScene: ecsScene,
		scene:    scene,
	}
}

type Renderer struct {
	ecsScene *ecs.Scene
	scene    *render.Scene
}

func (r *Renderer) Update() {
	result := r.ecsScene.Find(ecs.
		Having(ecscomp.PhysicsComponentID).
		And(ecscomp.RenderComponentID))
	defer result.Close()

	for result.HasNext() {
		entity := result.Next()

		physicsComp := ecscomp.GetPhysics(entity)
		body := physicsComp.Body

		renderComp := ecscomp.GetRender(entity)
		renderable := renderComp.Renderable

		renderable.Matrix = sprec.TransformationMat4(
			body.Orientation().OrientationX(),
			body.Orientation().OrientationY(),
			body.Orientation().OrientationZ(),
			body.Position(),
		)
		r.scene.Layout().InvalidateRenderable(renderable)
	}
}
