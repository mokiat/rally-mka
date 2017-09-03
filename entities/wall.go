package entities

import (
	"github.com/mokiat/rally-mka/collision"
	"github.com/mokiat/rally-mka/render"
)

type Wall struct {
	CollisionMesh *collision.Mesh
	RenderMesh    *render.Mesh
}
