package entities

import (
	"github.com/mokiat/rally-mka/data/m3d"
	"github.com/mokiat/rally-mka/render"
)

type Sky struct {
	Object     *m3d.Object
	RenderMesh *render.Mesh
}
