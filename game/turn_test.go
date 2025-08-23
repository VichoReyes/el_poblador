package game

import (
	"el_poblador/board"
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

	// Handle both possible outcomes: idle phase (normal roll) or robber placement (rolled 7)
	if robberPhase, ok := game.phase.(*phasePlaceRobber); ok {
		// If a 7 was rolled, complete the robber placement to get back to idle
		// Move robber to a different tile
		robberPhase.tileCoord, _ = board.NewTileCoord(1, 2)
		game.ConfirmAction(nil) // place robber
		// This might go to steal phase if there are players adjacent to robber
		// Keep confirming until we get back to idle
		for i := 0; i < 10; i++ { // safety limit
			if _, ok := game.phase.(*phaseIdle); ok {
				break
			}
			game.ConfirmAction(nil)
		}
	}
	
	// Should be in idle phase now
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be in idle phase after dice roll and robber handling, got: %T", game.phase)
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
