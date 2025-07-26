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
		if _, ok := game.phase.(*phaseInitialSettlements); !ok {
			t.Fatalf("Phase is not initial settlements. Iteration %d", i)
		}
		// find a valid location and move the cursor there
		coord, ok := findValidSettlementLocation(game)
		if !ok {
			t.Fatalf("Could not find a valid settlement location. Iteration %d", i)
		}
		game.phase.(*phaseInitialSettlements).cursorCross = coord

		game.ConfirmAction() // place settlement
		if _, ok := game.phase.(*phaseInitialRoad); !ok {
			t.Fatalf("Phase is not initial road. Iteration %d", i)
		}
		game.ConfirmAction() // place road
	}

	// there should be at least one resource
	// perhaps this can fail with very low probability
	resources := 0
	for _, player := range game.players {
		resources += player.TotalResources()
	}
	if resources == 0 {
		t.Fatal("No resources found")
	}

	if _, ok := game.phase.(*phaseDiceRoll); !ok {
		t.Fatal("Phase should be dice roll after initial flow is complete")
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
