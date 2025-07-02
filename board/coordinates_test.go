package board

import (
	"testing"
	"testing/quick"
)

func TestPathCoord(t *testing.T) {
	from := NewCrossCoord(2, 1)
	to := NewCrossCoord(3, 1)
	path := NewPathCoord(from, to)

	if path.From != from || path.To != to {
		t.Errorf("Expected path from %v to %v, got from %v to %v", from, to, path.From, path.To)
	}

	expected := "(2,1)-(3,1)"
	if path.String() != expected {
		t.Errorf("Expected %s, got %s", expected, path.String())
	}
}

func TestPathCoordCanonicalOrdering(t *testing.T) {
	// Test that path coordinates are always in canonical order
	from := NewCrossCoord(3, 1)
	to := NewCrossCoord(2, 1)
	path := NewPathCoord(from, to)

	// Should be reordered to canonical form
	expectedFrom := NewCrossCoord(2, 1)
	expectedTo := NewCrossCoord(3, 1)

	if path.From != expectedFrom || path.To != expectedTo {
		t.Errorf("Expected canonical ordering from %v to %v, got from %v to %v",
			expectedFrom, expectedTo, path.From, path.To)
	}
}

func TestCoordinateTypesAreDistinct(t *testing.T) {
	// This test verifies that the types are different at runtime
	intersection := NewCrossCoord(1, 1)
	tile := NewTileCoord(1, 1)

	// They should have different string representations
	if intersection.String() == tile.String() {
		t.Error("Intersection and tile coordinates should have different string representations")
	}

	// They should be different types (this is checked at compile time, but we can verify the behavior)
	if intersection.String() == "(1,1)" && tile.String() == "[1,1]" {
		// This is the expected behavior
	} else {
		t.Error("Unexpected string representations")
	}
}
func TestCrossCoordTraversal(t *testing.T) {
	// Create a set to track visited coordinates
	visited := make(map[CrossCoord]bool)

	// Create a queue for BFS traversal
	queue := []CrossCoord{NewCrossCoord(0, 0)}

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
		x = x % 10
		y = y % 10

		coord := NewCrossCoord(x, y)
		if !coord.IsInBounds() {
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
