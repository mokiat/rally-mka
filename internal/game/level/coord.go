package level

import "fmt"

func C(x, y int) Coord {
	return Coord{X: x, Y: y}
}

type Coord struct {
	X int
	Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c Coord) Neighbor(direction byte) Coord {
	return c.Neighbors()[direction%6]
}

func (c Coord) Neighbors() [6]Coord {
	if c.Y%2 == 0 {
		return [6]Coord{
			C(c.X+1, c.Y),
			C(c.X, c.Y+1),
			C(c.X-1, c.Y+1),
			C(c.X-1, c.Y),
			C(c.X-1, c.Y-1),
			C(c.X, c.Y-1),
		}
	} else {
		return [6]Coord{
			C(c.X+1, c.Y),
			C(c.X+1, c.Y+1),
			C(c.X, c.Y+1),
			C(c.X-1, c.Y),
			C(c.X, c.Y-1),
			C(c.X+1, c.Y-1),
		}
	}
}

func (c Coord) HasNeighborAt(coord Coord) bool {
	for _, neighborCoord := range c.Neighbors() {
		if neighborCoord == coord {
			return true
		}
	}
	return false
}
