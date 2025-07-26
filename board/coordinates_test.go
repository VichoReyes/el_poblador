package board

import (
	"slices"
	"testing"
	"testing/quick"
)

func TestPathCoord(t *testing.T) {
	from, _ := NewCrossCoord(2, 2)
	to, _ := NewCrossCoord(2, 3)
	path := NewPathCoord(from, to)

	if path.From != from || path.To != to {
		t.Errorf("Expected path from %v to %v, got from %v to %v", from, to, path.From, path.To)
	}

	expected := "(2,2)-(2,3)"
	if path.String() != expected {
		t.Errorf("Expected %s, got %s", expected, path.String())
	}
}

func TestPathCoordCanonicalOrdering(t *testing.T) {
	// Test that path coordinates are always in canonical order
	to, _ := NewCrossCoord(2, 3)
	from, _ := NewCrossCoord(2, 2)
	path := NewPathCoord(from, to)

	// Should be reordered to canonical form
	expectedFrom, _ := NewCrossCoord(2, 2)
	expectedTo, _ := NewCrossCoord(2, 3)

	if path.From != expectedFrom || path.To != expectedTo {
		t.Errorf("Expected canonical ordering from %v to %v, got from %v to %v",
			expectedFrom, expectedTo, path.From, path.To)
	}
}

func TestCrossCoordTraversal(t *testing.T) {
	// Create a set to track visited coordinates
	visited := make(map[CrossCoord]bool)

	// Create a queue for BFS traversal
	initial, _ := NewCrossCoord(1, 2)
	queue := []CrossCoord{initial}

	// BFS traversal
	for len(queue) > 0 {
		// Pop first element
		current := queue[0]
		queue = queue[1:]

		// Skip if already visited
		if visited[current] {
			continue
		}

		// Mark as visited
		visited[current] = true

		// Add unvisited neighbors to queue
		for _, neighbor := range current.Neighbors() {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check total number of valid coordinates
	expectedCount := 54
	if len(visited) != expectedCount {
		t.Errorf("Expected %d valid coordinates, got %d", expectedCount, len(visited))
	}
}

func TestCrossCoordNeighborsProperty(t *testing.T) {
	f := func(x, y int) bool {
		// Limit coordinates to a reasonable range to find valid ones
		x = x % 7
		y = y % 11

		coord, valid := NewCrossCoord(x, y)
		if !valid {
			return true // Skip out of bounds coordinates
		}

		neighbors := coord.Neighbors()

		// Property 1: Should have at least 2 neighbors
		if len(neighbors) < 2 {
			return false
		}

		// Property 2: Original coord should be neighbor of its neighbors
		for _, n := range neighbors {
			nNeighbors := n.Neighbors()
			found := false
			for _, nn := range nNeighbors {
				if nn == coord {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("Property test failed:", err)
	}
}

func TestNewDesertBoardTileCount(t *testing.T) {
	board := NewDesertBoard()

	// Count number of tiles
	tileCount := len(board.tiles)

	// Standard Catan board has 19 hexagonal tiles
	expectedTiles := 19
	if tileCount != expectedTiles {
		t.Errorf("Expected %d tiles in desert board, got %d", expectedTiles, tileCount)
	}
}

func TestTileRender(t *testing.T) {
	tile := Tile{Terrain: Arcilla, DiceNumber: 5}
	rendered1 := tile.RenderTile(false)
	rendered2 := tile.RenderTile(true)

	lengthsExpected := [5]int{8, 10, 10, 10, 8}
	for i := range rendered1 {
		actualLength1 := terminalLength(rendered1[i])
		if actualLength1 != lengthsExpected[i] {
			t.Errorf("Expected line %d to be %d characters long, got %d", i, lengthsExpected[i], actualLength1)
		}
		actualLength2 := terminalLength(rendered2[i])
		if actualLength2 != lengthsExpected[i] {
			t.Errorf("Expected line %d to be %d characters long, got %d", i, lengthsExpected[i], actualLength2)
		}
	}
}

func TestAdjacentTilesBasic(t *testing.T) {
	tests := []struct {
		x, y          int
		expectedCount int
		expectedTiles []TileCoord
	}{
		{0, 2, 1, []TileCoord{{X: 0, Y: 3}}},
		{0, 3, 1, []TileCoord{{X: 0, Y: 3}}},
		{2, 4, 3, []TileCoord{{X: 1, Y: 4}, {X: 2, Y: 3}, {X: 2, Y: 5}}},
		{2, 5, 3, []TileCoord{{X: 2, Y: 5}, {X: 1, Y: 4}, {X: 1, Y: 6}}},
	}

	for _, tt := range tests {
		coord, _ := NewCrossCoord(tt.x, tt.y)
		tiles := coord.AdjacentTiles()
		if len(tiles) != tt.expectedCount {
			t.Errorf("For coord (%d,%d): expected %d adjacent tiles, got %d",
				tt.x, tt.y, tt.expectedCount, len(tiles))
		}
		if tt.expectedTiles != nil {
			if !slices.Equal(tiles, tt.expectedTiles) {
				t.Errorf("For coord (%d,%d): expected tiles %v, got %v",
					tt.x, tt.y, tt.expectedTiles, tiles)
			}
		}
	}
}

func TestAdjacentTilesFull(t *testing.T) {
	// each tile should be seen 6 times
	timesSeen := make(map[TileCoord]int)
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			coord, ok := NewCrossCoord(x, y)
			if !ok {
				continue
			}
			tiles := coord.AdjacentTiles()
			for _, tile := range tiles {
				timesSeen[tile]++
			}
		}
	}
	for tile, count := range timesSeen {
		if count != 6 {
			t.Errorf("Expected tile %v to be seen 6 times, got %d", tile, count)
		}
	}
	if len(timesSeen) != 19 {
		t.Errorf("Expected 19 tiles, got %d", len(timesSeen))
	}
}
