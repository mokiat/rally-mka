package ecssys

import (
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
)

func NewRenderer(ecsScene *ecs.Scene) *Renderer {
	return &Renderer{
		ecsScene: ecsScene,
	}
}

type Renderer struct {
	ecsScene *ecs.Scene
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
		mesh := renderComp.Mesh

		mesh.SetPosition(body.Position())
		mesh.SetRotation(body.Orientation())
	}
}
