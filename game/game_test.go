package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestInitialSettlementsRender(t *testing.T) {
	game := &Game{}
	game.Start([]string{"john", "jane", "jim"})

	expectedTurns := []int{0, 1, 2, 2, 1, 0}

	// 2 rounds of placing settlements and roads
	for i := 0; i < 2*len(game.players); i++ {
		help := game.helpText(100)
		expectedTurn := game.players[expectedTurns[i]].Name
		if !strings.Contains(help, expectedTurn) {
			t.Fatalf("Help text '%s' does not contain expected turn '%s'. Iteration %d", help, expectedTurn, i)
		}

		if _, ok := game.phase.(*phaseInitialSettlements); !ok {
			t.Fatalf("Phase is not initial settlements. Iteration %d", i)
		}
		// find a valid location and move the cursor there
		coord, ok := findValidSettlementLocation(game)
		if !ok {
			t.Fatalf("Could not find a valid settlement location. Iteration %d", i)
		}
		game.phase.(*phaseInitialSettlements).cursorCross = coord

		game.ConfirmAction(nil) // place settlement
		if _, ok := game.phase.(*phaseInitialRoad); !ok {
			t.Fatalf("Phase is not initial road. Iteration %d", i)
		}
		game.ConfirmAction(nil) // place road
	}
	if _, ok := game.phase.(*phaseDiceRoll); !ok {
		t.Fatal("Phase should be dice roll after initial flow is complete")
	}
	game.ConfirmAction(nil) // roll dice
	if game.lastDice[0] == 0 {
		t.Fatal("Dice should be rolled")
	}
	diceText := fmt.Sprintf("Dice: %d (%d + %d)", game.lastDice[0]+game.lastDice[1], game.lastDice[0], game.lastDice[1])
	fullRender := game.Print(60, 50, nil)
	if !strings.Contains(fullRender, diceText) {
		t.Fatalf("Game render does not contain expected dice text '%s'", diceText)
	}
}
