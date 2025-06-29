package board

import (
	"testing"
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
