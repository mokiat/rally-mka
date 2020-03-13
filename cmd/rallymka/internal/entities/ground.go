package entities

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/collision"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/render"
)

type Ground struct {
	CollisionMesh *collision.Mesh
	RenderMesh    *render.Mesh
}