package game

import (
	"el_poblador/board"
	"testing"
)

func TestBankTrade(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	// Skip initial setup and go directly to idle phase
	game.phase = PhaseIdle(game)

	// Give player 4 brick and 2 wheat
	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceBrick)
	}
	for i := 0; i < 2; i++ {
		player.AddResource(board.ResourceWheat)
	}

	// Verify initial resources
	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Expected 4 brick, got %d", player.Resources[board.ResourceBrick])
	}
	if player.Resources[board.ResourceWheat] != 2 {
		t.Fatalf("Expected 2 wheat, got %d", player.Resources[board.ResourceWheat])
	}
	if player.Resources[board.ResourceOre] != 0 {
		t.Fatalf("Expected 0 ore, got %d", player.Resources[board.ResourceOre])
	}

	// Select "Trade" option (second option in idle phase)
	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	// Should be in trade offer phase now
	tradePhase, ok := game.phase.(*phaseTradeOffer)
	if !ok {
		t.Fatalf("Should be in trade offer phase after selecting Trade, got: %T", game.phase)
	}

	// Verify initial offer is all zeros
	for _, resourceType := range board.RESOURCE_TYPES {
		if tradePhase.offer[resourceType] != 0 {
			t.Fatalf("Initial offer should be 0 for %s, got %d", resourceType, tradePhase.offer[resourceType])
		}
	}

	// Navigate to brick (index 4 in RESOURCE_TYPES: Ore, Wood, Sheep, Wheat, Brick)
	// RESOURCE_TYPES order is: Piedra (Ore), Madera (Wood), Lana (Sheep), Trigo (Wheat), Ladrillo (Brick)
	brickIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceBrick {
			brickIndex = i
			break
		}
	}
	if brickIndex == -1 {
		t.Fatal("Could not find brick in RESOURCE_TYPES")
	}

	// Move cursor to brick
	for i := 0; i < brickIndex; i++ {
		game.phase.MoveCursor("down")
	}

	// Increment brick offer to 4 using right arrow
	for i := 0; i < 4; i++ {
		game.phase.MoveCursor("right")
	}

	// Verify offer has 4 brick
	if tradePhase.offer[board.ResourceBrick] != 4 {
		t.Fatalf("Expected offer of 4 brick, got %d", tradePhase.offer[board.ResourceBrick])
	}

	// Navigate to "Confirm Trade" button (down from last resource)
	numResources := len(board.RESOURCE_TYPES)
	currentPosition := brickIndex
	stepsToConfirm := numResources - currentPosition
	for i := 0; i < stepsToConfirm; i++ {
		game.phase.MoveCursor("down")
	}

	// Confirm the trade offer
	game.ConfirmAction(nil)

	// Should be in trade select receive phase now
	receivePhase, ok := game.phase.(*phaseTradeSelectReceive)
	if !ok {
		t.Fatalf("Should be in trade select receive phase, got: %T", game.phase)
	}

	// Find ore in the options
	oreIndex := -1
	for i, option := range receivePhase.options {
		if option == string(board.ResourceOre) {
			oreIndex = i
			break
		}
	}
	if oreIndex == -1 {
		t.Fatal("Ore should be in receive options")
	}

	// Move to ore option
	for i := 0; i < oreIndex; i++ {
		game.phase.MoveCursor("down")
	}

	// Confirm selection of ore
	game.ConfirmAction(nil)

	// Should be back in idle phase with notification
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be back in idle phase after completing trade, got: %T", game.phase)
	}

	// Verify resources were traded correctly
	if player.Resources[board.ResourceBrick] != 0 {
		t.Fatalf("Expected 0 brick after trade, got %d", player.Resources[board.ResourceBrick])
	}
	if player.Resources[board.ResourceWheat] != 2 {
		t.Fatalf("Expected 2 wheat (unchanged), got %d", player.Resources[board.ResourceWheat])
	}
	if player.Resources[board.ResourceOre] != 1 {
		t.Fatalf("Expected 1 ore after trade, got %d", player.Resources[board.ResourceOre])
	}
}

func TestBankTradeInvalidOffer(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	// Give player only 3 brick (not enough for bank trade)
	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 3; i++ {
		player.AddResource(board.ResourceBrick)
	}

	// Go to trade phase
	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseTradeOffer); !ok {
		t.Fatalf("Should be in trade offer phase, got: %T", game.phase)
	}

	// Try to offer 3 brick
	brickIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceBrick {
			brickIndex = i
			break
		}
	}

	for i := 0; i < brickIndex; i++ {
		game.phase.MoveCursor("down")
	}

	for i := 0; i < 3; i++ {
		game.phase.MoveCursor("right")
	}

	// Navigate to confirm button
	numResources := len(board.RESOURCE_TYPES)
	stepsToConfirm := numResources - brickIndex
	for i := 0; i < stepsToConfirm; i++ {
		game.phase.MoveCursor("down")
	}

	// Try to confirm (should fail validation and stay in same phase)
	game.ConfirmAction(nil)

	// Should still be in trade offer phase (validation failed)
	if _, ok := game.phase.(*phaseTradeOffer); !ok {
		t.Fatalf("Should still be in trade offer phase after invalid trade, got: %T", game.phase)
	}

	// Resources should be unchanged
	if player.Resources[board.ResourceBrick] != 3 {
		t.Fatalf("Resources should be unchanged, expected 3 brick, got %d", player.Resources[board.ResourceBrick])
	}
}

func TestBankTradeCancel(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceBrick)
	}

	// Go to trade phase
	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseTradeOffer); !ok {
		t.Fatalf("Should be in trade offer phase, got: %T", game.phase)
	}

	// Navigate to cancel button (resources + confirm + cancel)
	numResources := len(board.RESOURCE_TYPES)
	for i := 0; i <= numResources; i++ {
		game.phase.MoveCursor("down")
	}

	// Confirm cancel
	game.ConfirmAction(nil)

	// Should be back in idle phase
	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be back in idle phase after cancel, got: %T", game.phase)
	}

	// Resources should be unchanged
	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Resources should be unchanged, expected 4 brick, got %d", player.Resources[board.ResourceBrick])
	}
}

func TestBankTradeCancelFromReceivePhase(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceBrick)
	}

	// Go to trade phase
	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	// Set up valid offer
	brickIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceBrick {
			brickIndex = i
			break
		}
	}

	for i := 0; i < brickIndex; i++ {
		game.phase.MoveCursor("down")
	}

	for i := 0; i < 4; i++ {
		game.phase.MoveCursor("right")
	}

	// Navigate to confirm
	numResources := len(board.RESOURCE_TYPES)
	stepsToConfirm := numResources - brickIndex
	for i := 0; i < stepsToConfirm; i++ {
		game.phase.MoveCursor("down")
	}

	game.ConfirmAction(nil)

	// Should be in receive phase
	receivePhase, ok := game.phase.(*phaseTradeSelectReceive)
	if !ok {
		t.Fatalf("Should be in trade select receive phase, got: %T", game.phase)
	}

	// Navigate to cancel (last option)
	for i := 0; i < len(receivePhase.options)-1; i++ {
		game.phase.MoveCursor("down")
	}

	// Confirm cancel
	game.ConfirmAction(nil)

	// Should be back in trade offer phase with offer preserved
	newTradePhase, ok := game.phase.(*phaseTradeOffer)
	if !ok {
		t.Fatalf("Should be back in trade offer phase after cancel from receive, got: %T", game.phase)
	}

	// Verify offer was preserved
	if newTradePhase.offer[board.ResourceBrick] != 4 {
		t.Fatalf("Offer should be preserved, expected 4 brick, got %d", newTradePhase.offer[board.ResourceBrick])
	}

	// Resources should be unchanged
	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Resources should be unchanged, expected 4 brick, got %d", player.Resources[board.ResourceBrick])
	}
}
