package graphics

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
	ID             uint32
	VertexBufferID uint32
	IndexBufferID  uint32
}

type VertexArrayData struct {
	VertexData     []byte
	VertexStride   int32
	CoordOffset    int
	NormalOffset   int
	TexCoordOffset int
	IndexData      []byte
}

func (a *VertexArray) Allocate(data VertexArrayData) error {
	gl.GenVertexArrays(1, &a.ID)
	if a.ID == 0 {
		return fmt.Errorf("failed to allocate vertex array")
	}
	gl.BindVertexArray(a.ID)

	gl.GenBuffers(1, &a.VertexBufferID)
	if a.VertexBufferID == 0 {
		return fmt.Errorf("failed to allocate vertex buffer")
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, a.VertexBufferID)
	gl.BufferData(gl.ARRAY_BUFFER, len(data.VertexData), gl.Ptr(data.VertexData), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, data.VertexStride, gl.PtrOffset(data.CoordOffset))
	if data.NormalOffset != 0 {
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, data.VertexStride, gl.PtrOffset(data.NormalOffset))
	}
	if data.TexCoordOffset != 0 {
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, data.VertexStride, gl.PtrOffset(data.TexCoordOffset))
	}

	gl.GenBuffers(1, &a.IndexBufferID)
	if a.IndexBufferID == 0 {
		return fmt.Errorf("failed to allocate index buffer")
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, a.IndexBufferID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data.IndexData), gl.Ptr(data.IndexData), gl.STATIC_DRAW)
	return nil
}

func (a *VertexArray) Release() error {
	gl.DeleteBuffers(1, &a.IndexBufferID)
	gl.DeleteBuffers(1, &a.VertexBufferID)
	gl.DeleteVertexArrays(1, &a.ID)
	a.ID = 0
	a.VertexBufferID = 0
	a.IndexBufferID = 0
	return nil
}