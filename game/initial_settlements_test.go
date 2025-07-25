package game

import (
	"el_poblador/board"
	"testing"
)

func TestInitialFlow(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// 2 rounds of placing settlements and roads
	for i := 0; i < 2*len(game.players); i++ {
		if game.phase == nil {
			t.Fatalf("Phase is nil before initial flow is complete. Iteration %d", i)
		}
		// find a valid location and move the cursor there
		coord, ok := findValidSettlementLocation(game)
		if !ok {
			t.Fatalf("Could not find a valid settlement location. Iteration %d", i)
		}
		game.phase.(*phaseInitialSettlements).cursorCross = coord

		game.ConfirmAction() // place settlement
		if game.phase == nil {
			t.Fatalf("Phase is nil before initial flow is complete. Iteration %d", i)
		}
		game.ConfirmAction() // place road
	}

	if game.phase != nil {
		t.Fatal("Phase should be nil after initial flow is complete")
	}
}

func TestPlaceSettlementOnExisting(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// Place the first settlement
	game.ConfirmAction()
	// get placed settlement location
	settlementLocation := game.phase.(*phaseInitialRoad).sourceCross
	// place the road
	game.ConfirmAction()

	// Try to place another settlement in the same location
	// move cursor to the same spot
	game.phase.(*phaseInitialSettlements).cursorCross = settlementLocation
	// now in the same spot, try to place another settlement
	game.ConfirmAction()

	// after the second confirm, the phase should still be settlement placement
	if _, ok := game.phase.(*phaseInitialSettlements); !ok {
		t.Fatal("Should not be able to place a settlement on an existing one")
	}
}

func findValidSettlementLocation(game *Game) (board.CrossCoord, bool) {
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			coord, valid := board.NewCrossCoord(x, y)
			if valid && game.board.CanPlaceSettlement(coord) {
				return coord, true
			}
		}
	}
	return board.CrossCoord{}, false
}
