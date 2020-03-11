package resource

import (
	"encoding/json"
	"fmt"
	"io"
)

type Mesh struct {
	Coords    []float32 `json:"coords"`
	Normals   []float32 `json:"normals"`
	TexCoords []float32 `json:"tex_coords"`
	Indices   []int     `json:"indices"`
	SubMeshes []SubMesh `json:"sub_meshes"`
}

type SubMesh struct {
	Name           string `json:"name"`
	IndexOffset    int    `json:"index_offset"`
	IndexCount     int    `json:"index_count"`
	DiffuseTexture string `json:"diffuse_texture"`
}

func NewMeshDecoder() *MeshDecoder {
	return &MeshDecoder{}
}

type MeshDecoder struct{}

func (d *MeshDecoder) Decode(in io.Reader) (*Mesh, error) {
	var mesh Mesh
	if err := json.NewDecoder(in).Decode(&mesh); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return &mesh, nil
}

func NewMeshEncoder() *MeshEncoder {
	return &MeshEncoder{}
}

type MeshEncoder struct{}

func (e *MeshEncoder) Encode(out io.Writer, mesh *Mesh) error {
	if err := json.NewEncoder(out).Encode(mesh); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}
