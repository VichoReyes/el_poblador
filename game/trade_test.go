package game

import (
	"el_poblador/board"
	"testing"
)

func TestBankTrade(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceBrick)
	}
	for i := 0; i < 2; i++ {
		player.AddResource(board.ResourceWheat)
	}

	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Expected 4 brick, got %d", player.Resources[board.ResourceBrick])
	}
	if player.Resources[board.ResourceOre] != 0 {
		t.Fatalf("Expected 0 ore, got %d", player.Resources[board.ResourceOre])
	}

	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	tradePhase, ok := game.phase.(*phaseTradeOffer)
	if !ok {
		t.Fatalf("Should be in trade offer phase, got: %T", game.phase)
	}

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

	if tradePhase.offer[board.ResourceBrick] != 4 {
		t.Fatalf("Expected offer of 4 brick, got %d", tradePhase.offer[board.ResourceBrick])
	}

	game.ConfirmAction(nil)

	receivePhase, ok := game.phase.(*phaseTradeSelectReceive)
	if !ok {
		t.Fatalf("Should be in receive phase, got: %T", game.phase)
	}

	oreIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceOre {
			oreIndex = i
			break
		}
	}

	for i := 0; i < oreIndex; i++ {
		game.phase.MoveCursor("down")
	}

	game.phase.MoveCursor("right")

	if receivePhase.request[board.ResourceOre] != 1 {
		t.Fatalf("Expected request of 1 ore, got %d", receivePhase.request[board.ResourceOre])
	}

	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be in idle phase after trade, got: %T", game.phase)
	}

	if player.Resources[board.ResourceBrick] != 0 {
		t.Fatalf("Expected 0 brick after trade, got %d", player.Resources[board.ResourceBrick])
	}
	if player.Resources[board.ResourceWheat] != 2 {
		t.Fatalf("Expected 2 wheat unchanged, got %d", player.Resources[board.ResourceWheat])
	}
	if player.Resources[board.ResourceOre] != 1 {
		t.Fatalf("Expected 1 ore after trade, got %d", player.Resources[board.ResourceOre])
	}
}

func TestBankTradeInvalidOffer(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceBrick)
	}
	for i := 0; i < 2; i++ {
		player.AddResource(board.ResourceWheat)
	}

	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseTradeOffer); !ok {
		t.Fatalf("Should be in offer phase, got: %T", game.phase)
	}

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

	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseTradeSelectReceive); !ok {
		t.Fatalf("Should be in receive phase, got: %T", game.phase)
	}

	sheepIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceSheep {
			sheepIndex = i
			break
		}
	}

	for i := 0; i < sheepIndex; i++ {
		game.phase.MoveCursor("down")
	}

	for i := 0; i < 2; i++ {
		game.phase.MoveCursor("right")
	}

	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseIdle); !ok {
		t.Fatalf("Should be in idle phase with error, got: %T", game.phase)
	}

	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Resources unchanged, expected 4 brick, got %d", player.Resources[board.ResourceBrick])
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

	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	tradePhase, ok := game.phase.(*phaseTradeOffer)
	if !ok {
		t.Fatalf("Should be in offer phase, got: %T", game.phase)
	}

	canceledPhase := tradePhase.Cancel()

	if _, ok := canceledPhase.(*phaseIdle); !ok {
		t.Fatalf("Should be in idle phase after cancel, got: %T", canceledPhase)
	}

	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Resources unchanged, expected 4 brick, got %d", player.Resources[board.ResourceBrick])
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

	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

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

	game.ConfirmAction(nil)

	receivePhase, ok := game.phase.(*phaseTradeSelectReceive)
	if !ok {
		t.Fatalf("Should be in receive phase, got: %T", game.phase)
	}

	canceledPhase := receivePhase.Cancel()

	newTradePhase, ok := canceledPhase.(*phaseTradeOffer)
	if !ok {
		t.Fatalf("Should be back in offer phase after cancel, got: %T", canceledPhase)
	}

	if newTradePhase.offer[board.ResourceBrick] != 4 {
		t.Fatalf("Offer preserved, expected 4 brick, got %d", newTradePhase.offer[board.ResourceBrick])
	}

	if player.Resources[board.ResourceBrick] != 4 {
		t.Fatalf("Resources unchanged, expected 4 brick, got %d", player.Resources[board.ResourceBrick])
	}
}

func TestInvalidTradeNotYetImplemented(t *testing.T) {
	game := &Game{}
	game.Start([]string{"p1", "p2", "p3"})

	game.phase = PhaseIdle(game)

	player := &game.Players[game.PlayerTurn]
	for i := 0; i < 4; i++ {
		player.AddResource(board.ResourceWood)
	}
	for i := 0; i < 2; i++ {
		player.AddResource(board.ResourceSheep)
	}

	game.phase.MoveCursor("down")
	game.ConfirmAction(nil)

	woodIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceWood {
			woodIndex = i
			break
		}
	}

	for i := 0; i < woodIndex; i++ {
		game.phase.MoveCursor("down")
	}

	for i := 0; i < 4; i++ {
		game.phase.MoveCursor("right")
	}

	game.ConfirmAction(nil)

	if _, ok := game.phase.(*phaseTradeSelectReceive); !ok {
		t.Fatalf("Should be in receive phase, got: %T", game.phase)
	}

	sheepIndex := -1
	for i, rt := range board.RESOURCE_TYPES {
		if rt == board.ResourceSheep {
			sheepIndex = i
			break
		}
	}

	for i := 0; i < sheepIndex; i++ {
		game.phase.MoveCursor("down")
	}

	for i := 0; i < 2; i++ {
		game.phase.MoveCursor("right")
	}

	game.ConfirmAction(nil)

	idlePhase, ok := game.phase.(*phaseIdle)
	if !ok {
		t.Fatalf("Should be in idle phase, got: %T", game.phase)
	}

	if idlePhase.notification != "Trade type not yet implemented" {
		t.Fatalf("Expected 'not yet implemented', got: %s", idlePhase.notification)
	}

	if player.Resources[board.ResourceWood] != 4 {
		t.Fatalf("Resources unchanged, expected 4 wood, got %d", player.Resources[board.ResourceWood])
	}
	if player.Resources[board.ResourceSheep] != 2 {
		t.Fatalf("Resources unchanged, expected 2 sheep, got %d", player.Resources[board.ResourceSheep])
	}
}
