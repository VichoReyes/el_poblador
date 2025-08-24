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

func TestRoadPurchase(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// Directly set up a settlement for the current player to connect roads to
	playerId := game.playerTurn
	settlementCoord, _ := board.NewCrossCoord(2, 3)
	game.board.SetSettlement(settlementCoord, playerId)

	game.phase = PhaseIdle(game)

	// Give the player enough resources to build a road
	player := &game.players[game.playerTurn]
	player.AddResource(board.ResourceWood)
	player.AddResource(board.ResourceBrick)

	// Select "Build" option (first option in idle phase)
	game.ConfirmAction(nil)
	
	// Should be in building phase now
	if _, ok := game.phase.(*phaseBuilding); !ok {
		t.Fatalf("Should be in building phase after selecting Build, got: %T", game.phase)
	}

	// Select "Road" option (first option in building phase)  
	game.ConfirmAction(nil)

	// Should be in road start phase
	roadStartPhase, ok := game.phase.(*phaseRoadStart)
	if !ok {
		t.Fatalf("Should be in road start phase after selecting Road, got: %T", game.phase)
	}

	// Set cursor to the settlement position
	roadStartPhase.cursorCross = settlementCoord

	// Confirm road start position
	game.ConfirmAction(nil)

	// Should be in road end phase
	roadEndPhase, ok := game.phase.(*phaseRoadEnd)
	if !ok {
		t.Fatalf("Should be in road end phase after confirming start, got: %T", game.phase)
	}

	// Find a valid neighbor for the road end
	neighbors := settlementCoord.Neighbors()
	found := false
	for _, neighbor := range neighbors {
		pathCoord := board.NewPathCoord(settlementCoord, neighbor)
		if game.board.CanPlaceRoad(pathCoord, playerId) {
			roadEndPhase.cursorCross = neighbor
			found = true
			break
		}
	}
	
	if !found {
		t.Fatal("No valid neighbor found to complete road placement")
	}
	
	// Confirm road end position to complete the road purchase
	game.ConfirmAction(nil)

	// Verify the road was actually built by checking:
	// 1. Resources were consumed (player should no longer be able to build another road)
	// 2. The path can no longer be placed (because a road is there)
	pathCoord := board.NewPathCoord(settlementCoord, roadEndPhase.cursorCross)
	if player.CanBuildRoad() {
		t.Fatal("Resources were not consumed - player can still build roads")
	}
	if game.board.CanPlaceRoad(pathCoord, playerId) {
		t.Fatal("Road was not built - path is still available")
	}

	// BUG: This should be idle phase but currently returns to building phase
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be in idle phase after building road, got: %T", game.phase)
	}
}
