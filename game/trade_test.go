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

// setupTestGame creates a minimal game for trade testing
func setupTestGame() *Game {
	g := &Game{}
	g.Start([]string{"Alice", "Bob", "Charlie"})
	return g
}

func TestTradeOffer_canTake_Success(t *testing.T) {
	g := setupTestGame()

	// Give player 0 resources to offer
	g.Players[0].Resources[board.ResourceWood] = 2
	g.Players[0].Resources[board.ResourceBrick] = 1

	// Give player 1 resources to accept the trade
	g.Players[1].Resources[board.ResourceWheat] = 1
	g.Players[1].Resources[board.ResourceSheep] = 1

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1, // offer to anyone
		Offering: map[board.ResourceType]int{
			board.ResourceWood:  2,
			board.ResourceBrick: 1,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
			board.ResourceSheep: 1,
		},
		Status: OfferIsPending,
	}

	result := offer.canTake(1, g)
	if result != CanTakeTrue {
		t.Errorf("Expected CanTakeTrue, got %v", result)
	}
}

func TestTradeOffer_canTake_NotEnoughResources(t *testing.T) {
	g := setupTestGame()

	// Player 1 doesn't have enough resources
	g.Players[1].Resources[board.ResourceWheat] = 0
	g.Players[1].Resources[board.ResourceSheep] = 1

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 2,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
			board.ResourceSheep: 1,
		},
		Status: OfferIsPending,
	}

	result := offer.canTake(1, g)
	if result != CanTakeNotEnoughResources {
		t.Errorf("Expected CanTakeNotEnoughResources, got %v", result)
	}
}

func TestTradeOffer_canTake_ExcludedPlayer(t *testing.T) {
	g := setupTestGame()

	// Give player 1 resources but offer is targeted to player 2
	g.Players[1].Resources[board.ResourceWheat] = 1
	g.Players[1].Resources[board.ResourceSheep] = 1

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  2, // offer specifically to player 2
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 2,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsPending,
	}

	result := offer.canTake(1, g) // player 1 tries to accept
	if result != CanTakeExcludedPlayer {
		t.Errorf("Expected CanTakeExcludedPlayer, got %v", result)
	}

	// But player 2 should be able to accept (if they have resources)
	g.Players[2].Resources[board.ResourceWheat] = 1
	result = offer.canTake(2, g)
	if result != CanTakeTrue {
		t.Errorf("Expected CanTakeTrue for targeted player, got %v", result)
	}
}

func TestTradeOffer_canTake_Obsolete(t *testing.T) {
	g := setupTestGame()

	g.Players[1].Resources[board.ResourceWheat] = 1

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 1,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsCompleted, // already completed
	}

	result := offer.canTake(1, g)
	if result != CanTakeObsolete {
		t.Errorf("Expected CanTakeObsolete, got %v", result)
	}
}

func TestTradeOffer_canTake_Ambiguous(t *testing.T) {
	g := setupTestGame()

	g.Players[1].Resources[board.ResourceWheat] = 5

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceInvalid: 1, // ambiguous offer
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsPending,
	}

	result := offer.canTake(1, g)
	if result != CanTakeIsAmbiguous {
		t.Errorf("Expected CanTakeIsAmbiguous, got %v", result)
	}
}

func TestTradeOffer_executeTrade_Success(t *testing.T) {
	g := setupTestGame()

	// Setup initial resources
	g.Players[0].Resources[board.ResourceWood] = 3
	g.Players[0].Resources[board.ResourceBrick] = 2
	g.Players[0].Resources[board.ResourceWheat] = 0
	g.Players[0].Resources[board.ResourceSheep] = 0

	g.Players[1].Resources[board.ResourceWood] = 1
	g.Players[1].Resources[board.ResourceBrick] = 0
	g.Players[1].Resources[board.ResourceWheat] = 2
	g.Players[1].Resources[board.ResourceSheep] = 3

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceWood:  2,
			board.ResourceBrick: 1,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
			board.ResourceSheep: 2,
		},
		Status: OfferIsPending,
	}

	success := offer.executeTrade(1, g)
	if !success {
		t.Fatal("Trade execution failed")
	}

	// Check offerer (player 0) resources after trade
	if g.Players[0].Resources[board.ResourceWood] != 1 { // 3 - 2 = 1
		t.Errorf("Offerer wood: expected 1, got %d", g.Players[0].Resources[board.ResourceWood])
	}
	if g.Players[0].Resources[board.ResourceBrick] != 1 { // 2 - 1 = 1
		t.Errorf("Offerer brick: expected 1, got %d", g.Players[0].Resources[board.ResourceBrick])
	}
	if g.Players[0].Resources[board.ResourceWheat] != 1 { // 0 + 1 = 1
		t.Errorf("Offerer wheat: expected 1, got %d", g.Players[0].Resources[board.ResourceWheat])
	}
	if g.Players[0].Resources[board.ResourceSheep] != 2 { // 0 + 2 = 2
		t.Errorf("Offerer sheep: expected 2, got %d", g.Players[0].Resources[board.ResourceSheep])
	}

	// Check acceptor (player 1) resources after trade
	if g.Players[1].Resources[board.ResourceWood] != 3 { // 1 + 2 = 3
		t.Errorf("Acceptor wood: expected 3, got %d", g.Players[1].Resources[board.ResourceWood])
	}
	if g.Players[1].Resources[board.ResourceBrick] != 1 { // 0 + 1 = 1
		t.Errorf("Acceptor brick: expected 1, got %d", g.Players[1].Resources[board.ResourceBrick])
	}
	if g.Players[1].Resources[board.ResourceWheat] != 1 { // 2 - 1 = 1
		t.Errorf("Acceptor wheat: expected 1, got %d", g.Players[1].Resources[board.ResourceWheat])
	}
	if g.Players[1].Resources[board.ResourceSheep] != 1 { // 3 - 2 = 1
		t.Errorf("Acceptor sheep: expected 1, got %d", g.Players[1].Resources[board.ResourceSheep])
	}

	// Check offer status
	if offer.Status != OfferIsCompleted {
		t.Errorf("Expected offer status OfferIsCompleted, got %v", offer.Status)
	}
}

func TestTradeOffer_executeTrade_Failure_NotEnoughResources(t *testing.T) {
	g := setupTestGame()

	// Player 0 has resources to offer
	g.Players[0].Resources[board.ResourceWood] = 2

	// Player 1 doesn't have enough resources to accept
	g.Players[1].Resources[board.ResourceWheat] = 0

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 2,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsPending,
	}

	success := offer.executeTrade(1, g)
	if success {
		t.Error("Trade should have failed due to insufficient resources")
	}

	// Verify resources didn't change
	if g.Players[0].Resources[board.ResourceWood] != 2 {
		t.Error("Offerer resources should not have changed")
	}
	if g.Players[1].Resources[board.ResourceWheat] != 0 {
		t.Error("Acceptor resources should not have changed")
	}

	// Verify status didn't change
	if offer.Status != OfferIsPending {
		t.Error("Offer status should still be pending")
	}
}

func TestTradeOffer_executeTrade_Failure_AlreadyCompleted(t *testing.T) {
	g := setupTestGame()

	g.Players[0].Resources[board.ResourceWood] = 2
	g.Players[1].Resources[board.ResourceWheat] = 1

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  -1,
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 2,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsCompleted,
	}

	success := offer.executeTrade(1, g)
	if success {
		t.Error("Trade should have failed because offer is already completed")
	}

	// Verify resources didn't change
	if g.Players[0].Resources[board.ResourceWood] != 2 {
		t.Error("Offerer resources should not have changed")
	}
	if g.Players[1].Resources[board.ResourceWheat] != 1 {
		t.Error("Acceptor resources should not have changed")
	}
}

func TestTradeOffer_executeTrade_TargetedOffer(t *testing.T) {
	g := setupTestGame()

	// Setup resources
	g.Players[0].Resources[board.ResourceWood] = 3
	g.Players[1].Resources[board.ResourceWheat] = 2
	g.Players[2].Resources[board.ResourceWheat] = 2

	offer := TradeOffer{
		OffererID: 0,
		TargetID:  2, // specifically to player 2
		Offering: map[board.ResourceType]int{
			board.ResourceWood: 2,
		},
		Requesting: map[board.ResourceType]int{
			board.ResourceWheat: 1,
		},
		Status: OfferIsPending,
	}

	// Player 1 tries to accept (should fail)
	success := offer.executeTrade(1, g)
	if success {
		t.Error("Player 1 should not be able to accept offer targeted to player 2")
	}

	// Player 2 accepts (should succeed)
	success = offer.executeTrade(2, g)
	if !success {
		t.Fatal("Player 2 should be able to accept the targeted offer")
	}

	// Verify the trade happened between player 0 and player 2
	if g.Players[0].Resources[board.ResourceWood] != 1 { // 3 - 2 = 1
		t.Errorf("Offerer wood: expected 1, got %d", g.Players[0].Resources[board.ResourceWood])
	}
	if g.Players[0].Resources[board.ResourceWheat] != 1 { // 0 + 1 = 1
		t.Errorf("Offerer wheat: expected 1, got %d", g.Players[0].Resources[board.ResourceWheat])
	}
	if g.Players[2].Resources[board.ResourceWood] != 2 { // 0 + 2 = 2
		t.Errorf("Acceptor wood: expected 2, got %d", g.Players[2].Resources[board.ResourceWood])
	}
	if g.Players[2].Resources[board.ResourceWheat] != 1 { // 2 - 1 = 1
		t.Errorf("Acceptor wheat: expected 1, got %d", g.Players[2].Resources[board.ResourceWheat])
	}

	// Player 1's resources should be unchanged
	if g.Players[1].Resources[board.ResourceWheat] != 2 {
		t.Error("Player 1 resources should not have changed")
	}
}
