package entities

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/render"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

type Ground struct {
	CollisionMesh *collision.Mesh
	RenderMesh    *render.Mesh
}
