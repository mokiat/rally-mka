package level

import (
	"math/rand/v2"
)

type GeneratorConfig struct {
	Random *rand.Rand
	Size   int
}

func NewGenerator(config GeneratorConfig) *Generator {
	shapeSequences := make([][]ShapeKind, config.Size*config.Size)
	for i := range shapeSequences {
		shapeSequences[i] = []ShapeKind{
			ShapeKindTerrain,
			ShapeKindRoadStraight,
			ShapeKindRoadCornerSmooth,
			ShapeKindRoadCornerSharp,
			ShapeKindRoadSplit,
		}
	}

	rotationSequences := make([][]byte, config.Size*config.Size)
	for i := range rotationSequences {
		rotationSequences[i] = []byte{0, 1, 2, 3, 4, 5}
	}

	return &Generator{
		random:            config.Random,
		shapeSequences:    shapeSequences,
		rotationSequences: rotationSequences,
		boardSize:         config.Size,
	}
}

type Generator struct {
	random            *rand.Rand
	shapeSequences    [][]ShapeKind
	rotationSequences [][]byte
	boardSize         int
}

func (g *Generator) Generate() *Board {
	for _, shapeSequence := range g.shapeSequences {
		shuffleSlice(g.random, shapeSequence)
	}
	for _, rotationSequence := range g.rotationSequences {
		shuffleSlice(g.random, rotationSequence)
	}

	board := NewBoard(g.boardSize)
	if !g.generateTile(board, 0) {
		panic("failed to generate a level")
	}
	return board
}

func (g *Generator) generateTile(board *Board, index int) bool {
	if index >= g.boardSize*g.boardSize {
		return true
	}

	coord := C(index%g.boardSize, index/g.boardSize)
	shapeSequence := g.shapeSequences[index]
	rotationSequence := g.rotationSequences[index]

	for _, shape := range shapeSequence {
		for _, rotation := range rotationSequence {
			tile := Tile{
				Shape:     shape,
				Ground:    GroundKindGrass,
				Road:      RoadKindDirt,
				Variation: 0,
				Rotation:  rotation,
			}
			if !g.canPlaceTile(board, tile, coord) {
				continue
			}
			board.SetTile(coord, tile)
			if g.generateTile(board, index+1) {
				return true
			}
		}
	}
	return false
}

func (g *Generator) canPlaceTile(board *Board, tile Tile, coord Coord) bool {
	// Check starting tile.
	if coord == board.Center() {
		isValid := tile.Shape == ShapeKindRoadStraight && (tile.Rotation == 0 || tile.Rotation == 3)
		if !isValid {
			return false
		}
	}
	// Check connection to the left.
	if !g.isValidConnection(board, tile, coord, 3) {
		return false
	}
	// Check connection to the top-left.
	if !g.isValidConnection(board, tile, coord, 4) {
		return false
	}
	// Check connection to the top-right.
	if !g.isValidConnection(board, tile, coord, 5) {
		return false
	}
	// Check boundaries.
	for direction := range byte(6) {
		if !board.ContainsCoord(coord.Neighbor(direction)) && tile.HasRoad(direction) {
			return false
		}
	}
	return true
}

func (g *Generator) isValidConnection(board *Board, tile Tile, coord Coord, direction byte) bool {
	neighborCoord := coord.Neighbor(direction)
	if !board.ContainsCoord(neighborCoord) {
		return !tile.HasRoad(direction)
	}
	neighborTile := board.Tile(neighborCoord)
	return tile.HasRoad(direction) == neighborTile.HasRoad((direction+3)%6)
}
