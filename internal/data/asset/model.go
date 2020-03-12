package asset

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Model struct {
	Meshes []Mesh
	Nodes  []Node
}

type Node struct {
	ParentIndex int16
	Name        string
	Matrix      [16]float32
	MeshIndex   uint16
}

func NewModelDecoder() *ModelDecoder {
	return &ModelDecoder{}
}

type ModelDecoder struct{}

func (d *ModelDecoder) Decode(in io.Reader) (*Model, error) {
	var model Model
	if err := gob.NewDecoder(in).Decode(&model); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &model, nil
}

func NewModelEncoder() *ModelEncoder {
	return &ModelEncoder{}
}

type ModelEncoder struct{}

func (e *ModelEncoder) Encode(out io.Writer, model *Model) error {
	if err := gob.NewEncoder(out).Encode(model); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}
