package game

import (
	"testing"
)

func TestBuildingPhase(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// Complete initial settlements
	for i := 0; i < 2*len(game.players); i++ {
		game.MoveCursorToPlaceSettlement()
		game.ConfirmAction(nil) // place settlement
		game.ConfirmAction(nil) // place road
	}
	game.ConfirmAction(nil) // roll dice

	// Should be in idle phase now
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatal("Should be in idle phase after dice roll")
	}

	// Select "Build" option (it's already the first option, so no need to move cursor)
	game.ConfirmAction(nil)

	// Should be in building phase now
	if _, ok := game.phase.(*phaseBuilding); !ok {
		t.Fatal("Should be in building phase after selecting Build")
	}

	// Test cancel functionality
	// Move cursor to "Cancel" (last option, index 4)
	for i := 0; i < 4; i++ {
		game.phase.MoveCursor("down")
	}
	game.ConfirmAction(nil)

	// Should be back in idle phase
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatal("Should be back in idle phase after canceling")
	}
}
