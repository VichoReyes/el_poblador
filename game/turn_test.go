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

func TestBuildingPhaseOptions(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// Complete initial settlements
	for i := 0; i < 2*len(game.players); i++ {
		game.MoveCursorToPlaceSettlement()
		game.ConfirmAction(nil) // place settlement
		game.ConfirmAction(nil) // place road
	}
	game.ConfirmAction(nil) // roll dice

	// Give player some resources to test building options
	player := game.players[game.playerTurn]
	player.AddResource(board.ResourceWood)
	player.AddResource(board.ResourceBrick)
	player.AddResource(board.ResourceWheat)
	player.AddResource(board.ResourceSheep)

	// Select "Build" option (it's already the first option, so no need to move cursor)
	game.ConfirmAction(nil)

	// Should be in building phase now
	buildingPhase, ok := game.phase.(*phaseBuilding)
	if !ok {
		t.Fatal("Should be in building phase after selecting Build")
	}

	// Check that we have the expected number of options (4 building options + Cancel)
	expectedOptions := 5
	if len(buildingPhase.options) != expectedOptions {
		t.Fatalf("Expected %d options, got %d", expectedOptions, len(buildingPhase.options))
	}

	// Check that the first option is Road
	if buildingPhase.options[0] != "Road (1 Wood, 1 Brick)" {
		t.Fatalf("Expected first option to be 'Road (1 Wood, 1 Brick)', got '%s'", buildingPhase.options[0])
	}

	// Check that the last option is Cancel
	if buildingPhase.options[len(buildingPhase.options)-1] != "Cancel" {
		t.Fatalf("Expected last option to be 'Cancel', got '%s'", buildingPhase.options[len(buildingPhase.options)-1])
	}
}
