package level

import (
	"fmt"

	"github.com/mokiat/gomath/dprec"
)

const (
	ShapeKindNone ShapeKind = iota
	ShapeKindTerrain
	ShapeKindRoadStraight
	ShapeKindRoadCornerSmooth
	ShapeKindRoadCornerSharp
	ShapeKindRoadSplit
)

type ShapeKind byte

const (
	GroundKindNone GroundKind = iota
	GroundKindGrass
)

type GroundKind byte

const (
	RoadKindNone RoadKind = iota
	RoadKindDirt
)

type RoadKind byte

type Tile struct {
	Shape     ShapeKind  `json:"shape"`
	Ground    GroundKind `json:"ground"`
	Road      RoadKind   `json:"road"`
	Variation byte       `json:"variation"`
	Rotation  byte       `json:"rotation"`
}

func (t Tile) NodeName() string {
	switch {
	case t.Shape == ShapeKindNone:
		return ""
	case t.Shape == ShapeKindTerrain && t.Ground == GroundKindGrass:
		const availableVariations = 1
		return fmt.Sprintf("Tile.Grass.v%d", (t.Variation%availableVariations)+1)
	case t.Shape == ShapeKindRoadStraight && t.Ground == GroundKindGrass && t.Road == RoadKindDirt:
		const availableVariations = 1
		return fmt.Sprintf("Tile.Grass.Dirt.Straight.v%d", (t.Variation%availableVariations)+1)
	case t.Shape == ShapeKindRoadCornerSmooth && t.Ground == GroundKindGrass && t.Road == RoadKindDirt:
		const availableVariations = 1
		return fmt.Sprintf("Tile.Grass.Dirt.Corner.Smooth.v%d", (t.Variation%availableVariations)+1)
	case t.Shape == ShapeKindRoadCornerSharp && t.Ground == GroundKindGrass && t.Road == RoadKindDirt:
		const availableVariations = 1
		return fmt.Sprintf("Tile.Grass.Dirt.Corner.Sharp.v%d", (t.Variation%availableVariations)+1)
	case t.Shape == ShapeKindRoadSplit && t.Ground == GroundKindGrass && t.Road == RoadKindDirt:
		const availableVariations = 1
		return fmt.Sprintf("Tile.Grass.Dirt.Split.v%d", (t.Variation%availableVariations)+1)
	default:
		panic("cannot resolve node name for tile")
	}
}

func (t Tile) RotationQuat() dprec.Quat {
	angle := dprec.Degrees(float64(t.Rotation) * 60.0)
	return dprec.RotationQuat(angle, dprec.BasisYVec3())
}

func (t Tile) HasRoad(direction byte) bool {
	localDirection := (direction + t.Rotation) % 6
	switch t.Shape {
	case ShapeKindRoadStraight:
		return localDirection == 0 || localDirection == 3
	case ShapeKindRoadCornerSmooth:
		return localDirection == 1 || localDirection == 3
	case ShapeKindRoadCornerSharp:
		return localDirection == 2 || localDirection == 3
	case ShapeKindRoadSplit:
		return localDirection == 1 || localDirection == 3 || localDirection == 5
	default:
		return false
	}
}
