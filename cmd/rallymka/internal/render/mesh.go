package render

import (
	"github.com/mokiat/go-whiskey-gl/buffer"
	"github.com/mokiat/go-whiskey-gl/texture"
)

const DisabledOffset = -1

type MeshData struct {
	VertexData     buffer.Float32DataPlayground
	CoordOffset    int
	NormalOffset   int
	ColorOffset    int
	TexCoordOffset int
	Stride         int
	IndexData      buffer.UInt16DataPlayground
	TextureData    texture.FlatDataPlayground
}

func NewMesh(data *MeshData) *Mesh {
	return &Mesh{
		data: data,
	}
}

type Mesh struct {
	data           *MeshData
	vertexBuffer   *buffer.VertexBuffer
	vertexStride   int32
	coordOffset    int
	normalOffset   int
	colorOffset    int
	texCoordOffset int
	indexBuffer    *buffer.IndexBuffer
	indexCount     int32
	texture        *texture.FlatTexture
}

func (m *Mesh) Generate() {
	m.vertexBuffer = buffer.NewVertexBuffer()
	if err := m.vertexBuffer.Allocate(); err != nil {
		panic(err)
	}
	m.vertexBuffer.Bind()
	m.vertexBuffer.CreateData(m.data.VertexData)

	m.vertexStride = int32(m.data.Stride)
	m.coordOffset = m.data.CoordOffset
	m.normalOffset = m.data.NormalOffset
	m.colorOffset = m.data.ColorOffset
	m.texCoordOffset = m.data.TexCoordOffset

	m.indexBuffer = buffer.NewIndexBuffer()
	if err := m.indexBuffer.Allocate(); err != nil {
		panic(err)
	}
	m.indexBuffer.Bind()
	m.indexBuffer.CreateData(m.data.IndexData)

	m.indexCount = int32(m.data.IndexData.Count())

	if m.data.TextureData != nil {
		m.texture = texture.NewFlatTexture()
		if err := m.texture.Allocate(); err != nil {
			panic(err)
		}
		m.texture.Bind()
		m.texture.CreateData(m.data.TextureData)
	}

	// release local resources
	m.data = nil
}
