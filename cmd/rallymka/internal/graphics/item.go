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
	Program        Program
	ModelMatrix    math.Mat4x4
	SkyboxTexture  uint32
	DiffuseTexture uint32
	VertexArrayID  uint32
	IndexCount     int
}

func (i *Item) reset() {
	i.IndexCount = 0
}
