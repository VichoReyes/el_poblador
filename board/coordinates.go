package board

import (
	"fmt"
	"slices"
)

// CrossCoord represents the coordinates of an intersection point
// where three hexagons meet.
type CrossCoord struct {
	X, Y int
}

// TileCoord represents the coordinates of a hexagonal tile.
type TileCoord struct {
	X, Y int
}

// PathCoord represents a path/edge between two intersections.
type PathCoord struct {
	From, To CrossCoord
}

// NewCrossCoord creates a intersection coordinate and returns whether it is valid
func NewCrossCoord(x, y int) (CrossCoord, bool) {
	coord := CrossCoord{X: x, Y: y}
	if coord.IsInBounds() {
		return coord, true
	}
	return CrossCoord{}, false
}

// IsInBounds checks if the coordinate is within the bounds of the board
// currently only for 3-4 player games
func (c CrossCoord) IsInBounds() bool {
	x := c.X
	y := c.Y
	// left edge
	if x < 0 {
		return false
	}
	// top left edge
	if x+y < 2 {
		return false
	}
	// top right edge
	if x-y > 3 {
		return false
	}
	// right edge
	if x > 5 {
		return false
	}
	// bottom right edge
	if x+y > 13 {
		return false
	}
	// bottom left edge
	if y-x > 8 {
		return false
	}
	return true
}

func (c CrossCoord) Up() (CrossCoord, bool) {
	return NewCrossCoord(c.X, c.Y-1)
}

func (c CrossCoord) Down() (CrossCoord, bool) {
	return NewCrossCoord(c.X, c.Y+1)
}

func (c CrossCoord) Left() (CrossCoord, bool) {
	if (c.X+c.Y)%2 == 0 {
		return CrossCoord{}, false
	}
	return NewCrossCoord(c.X-1, c.Y)
}

func (c CrossCoord) Right() (CrossCoord, bool) {
	if (c.X+c.Y)%2 == 1 {
		return CrossCoord{}, false
	}
	return NewCrossCoord(c.X+1, c.Y)
}

func (c TileCoord) Up() (TileCoord, bool) {
	return NewTileCoord(c.X, c.Y-2)
}

func (c TileCoord) Down() (TileCoord, bool) {
	return NewTileCoord(c.X, c.Y+2)
}

func (c TileCoord) Left() (TileCoord, bool) {
	return NewTileCoord(c.X-1, c.Y)
}

func (c TileCoord) Right() (TileCoord, bool) {
	return NewTileCoord(c.X+1, c.Y)
}

func (c CrossCoord) Neighbors() []CrossCoord {
	// TODO: use Up, Down, Left, Right methods
	var potential []CrossCoord
	if (c.X+c.Y)%2 == 0 {
		potential = []CrossCoord{
			{X: c.X + 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	} else {
		potential = []CrossCoord{
			{X: c.X - 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	}
	neighbors := []CrossCoord{}
	for _, p := range potential {
		if p.IsInBounds() {
			neighbors = append(neighbors, p)
		}
	}
	return neighbors
}

func (c CrossCoord) adjacentTileCoords() []TileCoord {
	var potential []TileCoord
	if (c.X+c.Y)%2 == 0 {
		potential = []TileCoord{
			{X: c.X - 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	} else {
		potential = []TileCoord{
			{X: c.X, Y: c.Y},
			{X: c.X - 1, Y: c.Y - 1},
			{X: c.X - 1, Y: c.Y + 1},
		}
	}
	neighbors := []TileCoord{}
	for _, p := range potential {
		if _, ok := NewTileCoord(p.X, p.Y); ok {
			neighbors = append(neighbors, p)
		}
	}
	return neighbors
}

// NewTileCoord creates a new tile coordinate and returns whether it is valid
func NewTileCoord(x, y int) (TileCoord, bool) {
	if (x+y)%2 == 0 {
		return TileCoord{}, false
	}
	_, leftOk := NewCrossCoord(x, y)
	_, acrossOk := NewCrossCoord(x+1, y)
	if leftOk && acrossOk {
		return TileCoord{X: x, Y: y}, true
	}
	return TileCoord{}, false
}

// NewPathCoord creates a new path coordinate between two intersections
// panics if from and to are not neighbors
func NewPathCoord(from, to CrossCoord) PathCoord {
	fromNeighbors := from.Neighbors()
	if !slices.Contains(fromNeighbors, to) {
		panic("from and to are not neighbors")
	}
	// Ensure canonical ordering (ascending)
	if (from.X < to.X) || (from.X == to.X && from.Y < to.Y) {
		return PathCoord{From: from, To: to}
	}
	return PathCoord{From: to, To: from}
}

// String returns the string representation of an intersection coordinate
func (ic CrossCoord) String() string {
	return fmt.Sprintf("(%d,%d)", ic.X, ic.Y)
}

// String returns the string representation of a tile coordinate
func (tc TileCoord) String() string {
	return fmt.Sprintf("[%d,%d]", tc.X, tc.Y)
}

// String returns the string representation of a path coordinate
func (pc PathCoord) String() string {
	return fmt.Sprintf("%s-%s", pc.From.String(), pc.To.String())
}
