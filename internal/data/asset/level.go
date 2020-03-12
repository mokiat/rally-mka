package asset

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Level struct {
	SkyboxTexture      string
	Waypoints          []Point
	StartCollisionMesh LevelCollisionMesh
	StaticEntities     []LevelEntity
	StaticMeshes       []Mesh
	CollisionMeshes    []LevelCollisionMesh
}

type LevelEntity struct {
	Model  string
	Matrix [16]float32
}

type LevelCollisionMesh struct {
	Triangles []Triangle
}

type Triangle [3]Point

type Point [3]float32

func NewLevelDecoder() *LevelDecoder {
	return &LevelDecoder{}
}

type LevelDecoder struct{}

func (d *LevelDecoder) Decode(in io.Reader) (*Level, error) {
	var level Level
	if err := gob.NewDecoder(in).Decode(&level); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &level, nil
}

func NewLevelEncoder() *LevelEncoder {
	return &LevelEncoder{}
}

type LevelEncoder struct{}

func (e *LevelEncoder) Encode(out io.Writer, level *Level) error {
	if err := gob.NewEncoder(out).Encode(level); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}
