package asset

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Mesh struct {
	VertexData     []byte
	VertexStride   uint8
	CoordOffset    uint8
	NormalOffset   uint8
	TexCoordOffset uint8
	IndexData      []byte
	SubMeshes      []SubMesh
}

type SubMesh struct {
	Name           string
	IndexOffset    uint32
	IndexCount     uint32
	DiffuseTexture string
}

func NewMeshDecoder() *MeshDecoder {
	return &MeshDecoder{}
}

type MeshDecoder struct{}

func (d *MeshDecoder) Decode(in io.Reader) (*Mesh, error) {
	var mesh Mesh
	if err := gob.NewDecoder(in).Decode(&mesh); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &mesh, nil
}

func NewMeshEncoder() *MeshEncoder {
	return &MeshEncoder{}
}

type MeshEncoder struct{}

func (e *MeshEncoder) Encode(out io.Writer, mesh *Mesh) error {
	if err := gob.NewEncoder(out).Encode(mesh); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}
