package level

import (
	"encoding/json"
)

func NewBoard(size int) *Board {
	return &Board{
		size:  size,
		tiles: make([]Tile, size*size),
	}
}

type Board struct {
	size  int
	tiles []Tile
}

func (b *Board) Size() int {
	return b.size
}

func (b *Board) Center() Coord {
	return C(b.size/2, b.size/2)
}

func (b *Board) ContainsCoord(coord Coord) bool {
	return coord.X >= 0 && coord.X < b.size && coord.Y >= 0 && coord.Y < b.size
}

func (b *Board) Tile(coord Coord) Tile {
	if !b.ContainsCoord(coord) {
		panic("coord out of bounds")
	}
	return b.tiles[coord.X+coord.Y*b.size]
}

func (b *Board) SetTile(coord Coord, tile Tile) {
	if !b.ContainsCoord(coord) {
		panic("coord out of bounds")
	}
	b.tiles[coord.X+coord.Y*b.size] = tile
}

func SerializeBoard(board *Board) ([]byte, error) {
	return json.Marshal(struct {
		Size  int    `json:"size"`
		Tiles []Tile `json:"tiles"`
	}{
		Size:  board.size,
		Tiles: board.tiles,
	})
}

func ParseBoard(data []byte) (*Board, error) {
	var parsed struct {
		Size  int    `json:"size"`
		Tiles []Tile `json:"tiles"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return &Board{
		size:  parsed.Size,
		tiles: parsed.Tiles,
	}, nil
}
