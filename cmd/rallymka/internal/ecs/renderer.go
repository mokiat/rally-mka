package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/render"
)

func NewRenderer(ecsManager *Manager, scene *render.Scene) *Renderer {
	return &Renderer{
		ecsManager: ecsManager,
		scene:      scene,
	}
}

type Renderer struct {
	ecsManager *Manager
	scene      *render.Scene
}

func (r *Renderer) Update(ctx game.UpdateContext) {
	for _, entity := range r.ecsManager.Entities() {
		renderComp := entity.Render
		physicsComp := entity.Physics
		if renderComp == nil || physicsComp == nil {
			continue
		}
		body := physicsComp.Body
		renderComp.Renderable.Matrix = sprec.TransformationMat4(
			body.Orientation.OrientationX(),
			body.Orientation.OrientationY(),
			body.Orientation.OrientationZ(),
			body.Position,
		)
		r.scene.Layout().InvalidateRenderable(renderComp.Renderable)
	}
}
