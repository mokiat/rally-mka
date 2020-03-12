package resource

import (
	"encoding/json"
	"fmt"
	"io"
)

type Model struct {
	Meshes []Mesh `json:"meshes"`
	Nodes  []Node `json:"nodes"`
}

type Node struct {
	ParentIndex int         `json:"parent_index"`
	Name        string      `json:"name"`
	Matrix      [16]float32 `json:"matrix"`
	MeshIndex   int         `json:"mesh_index"`
}

func NewModelDecoder() *ModelDecoder {
	return &ModelDecoder{}
}

type ModelDecoder struct{}

func (d *ModelDecoder) Decode(in io.Reader) (*Model, error) {
	var model Model
	if err := json.NewDecoder(in).Decode(&model); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return &model, nil
}

func NewModelEncoder() *ModelEncoder {
	return &ModelEncoder{}
}

type ModelEncoder struct{}

func (e *ModelEncoder) Encode(out io.Writer, model *Model) error {
	if err := json.NewEncoder(out).Encode(model); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}
