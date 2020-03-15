package graphics

import "github.com/mokiat/go-whiskey/math"

type RenderPrimitive int

const (
	RenderPrimitiveTriangles RenderPrimitive = iota
	RenderPrimitiveLines
)

func createItem() Item {
	return Item{}
}

type Item struct {
	Primitive      RenderPrimitive
	Program        *Program
	ModelMatrix    math.Mat4x4
	SkyboxTexture  *CubeTexture
	DiffuseTexture *TwoDTexture
	VertexArray    *VertexArray
	IndexCount     int32
}

func (i *Item) reset() {
	i.IndexCount = 0
}
