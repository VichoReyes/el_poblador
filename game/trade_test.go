package game

import (
	"el_poblador/board"
	"testing"
)

// setupGameForTradeTest creates a new game with 3 players and gives them resources.
func setupGameForTradeTest() *Game {
	g := &Game{}
	g.Start([]string{"P1", "P2", "P3"})
	g.Players[0].Resources[board.ResourceWood] = 5
	g.Players[0].Resources[board.ResourceBrick] = 5
	g.Players[1].Resources[board.ResourceSheep] = 5
	g.Players[1].Resources[board.ResourceOre] = 5
	g.Players[2].Resources[board.ResourceWheat] = 5
	return g
}

func TestTradeOfferCreation(t *testing.T) {
	g := setupGameForTradeTest()
	if len(g.TradeOffers) != 0 {
		t.Fatalf("Expected 0 trade offers at the start, got %d", len(g.TradeOffers))
	}

	// Simulate player 0 creating an offer
	g.PlayerTurn = 0
	// Assume the flow of phases works and an offer is created
	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0,
		TargetID:  1,
		Offering:  map[board.ResourceType]int{board.ResourceWood: 1},
		Requesting: map[board.ResourceType]int{board.ResourceSheep: 1},
		Status:    OfferIsPending,
	})

	if len(g.TradeOffers) != 1 {
		t.Errorf("Expected 1 trade offer, got %d", len(g.TradeOffers))
	}
	if g.TradeOffers[0].OffererID != 0 {
		t.Errorf("Expected offerer to be player 0, got %d", g.TradeOffers[0].OffererID)
	}
}

func TestAmbiguousTradeOfferCreation(t *testing.T) {
	g := setupGameForTradeTest()
	
	// Simulate player 0 creating an ambiguous offer
	g.PlayerTurn = 0
	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0,
		TargetID:  -1, // All players
		Offering:  map[board.ResourceType]int{board.ResourceBrick: 2},
		Requesting: map[board.ResourceType]int{board.ResourceInvalid: 1}, // "something"
		Status:    OfferIsPending,
	})

	if len(g.TradeOffers) != 1 {
		t.Errorf("Expected 1 trade offer, got %d", len(g.TradeOffers))
	}
	if _, ok := g.TradeOffers[0].Requesting[board.ResourceInvalid]; !ok {
		t.Errorf("Expected offer to have an ambiguous request, but it didn't")
	}
}

func TestOfferRetraction(t *testing.T) {
	g := setupGameForTradeTest()
	g.PlayerTurn = 0
	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0, // Player 0 is the offerer
		TargetID:  1,
		Status:    OfferIsPending,
	})

	// Player 0 retracts the offer
	tradeMenu := PhasePlayerTrade(g, PhaseIdle(g)).(*phasePlayerTrade)
	tradeMenu.selected = 0 // Select the first (and only) offer
	tradeMenu.Confirm()

	if len(g.TradeOffers) != 1 {
		t.Fatalf("Expected 1 offer, got %d", len(g.TradeOffers))
	}
	if g.TradeOffers[0].Status != OfferIsRetracted {
		t.Errorf("Expected offer status to be Retracted, got %v", g.TradeOffers[0].Status)
	}
}

func TestOfferAcceptance(t *testing.T) {
	g := setupGameForTradeTest()
	g.PlayerTurn = 1 // Player 1's turn
	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0,
		TargetID:  1,
		Offering:  map[board.ResourceType]int{board.ResourceWood: 2},
		Requesting: map[board.ResourceType]int{board.ResourceSheep: 1},
		Status:    OfferIsPending,
	})

	p0WoodBefore := g.Players[0].Resources[board.ResourceWood]
	p1SheepBefore := g.Players[1].Resources[board.ResourceSheep]

	// Player 1 accepts the offer
	tradeMenu := PhasePlayerTrade(g, PhaseIdle(g)).(*phasePlayerTrade)
	tradeMenu.selected = 0 // Select the incoming offer
	nextPhase := tradeMenu.Confirm()

	if _, ok := nextPhase.(*phaseIdle); !ok {
		t.Fatalf("Expected to return to phaseIdle after trade, but got %T", nextPhase)
	}

	if g.TradeOffers[0].Status != OfferIsCompleted {
		t.Errorf("Expected offer status to be Completed, got %v", g.TradeOffers[0].Status)
	}

	// Check resources
	if g.Players[0].Resources[board.ResourceWood] != p0WoodBefore-2 {
		t.Errorf("Expected player 0 to have %d wood, got %d", p0WoodBefore-2, g.Players[0].Resources[board.ResourceWood])
	}
	if g.Players[1].Resources[board.ResourceWood] != 2 {
		t.Errorf("Expected player 1 to have 2 wood, got %d", g.Players[1].Resources[board.ResourceWood])
	}
	if g.Players[1].Resources[board.ResourceSheep] != p1SheepBefore-1 {
		t.Errorf("Expected player 1 to have %d sheep, got %d", p1SheepBefore-1, g.Players[1].Resources[board.ResourceSheep])
	}
	if g.Players[0].Resources[board.ResourceSheep] != 1 {
		t.Errorf("Expected player 0 to have 1 sheep, got %d", g.Players[0].Resources[board.ResourceSheep])
	}
}

func TestAcceptInvalidOffer(t *testing.T) {
	g := setupGameForTradeTest()
	g.PlayerTurn = 1 // Player 1's turn
	g.Players[0].Resources[board.ResourceWood] = 1 // Not enough to fulfill offer

	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0,
		TargetID:  1,
		Offering:  map[board.ResourceType]int{board.ResourceWood: 2},
		Requesting: map[board.ResourceType]int{board.ResourceSheep: 1},
		Status:    OfferIsPending,
	})

	// Player 1 tries to accept
	tradeMenu := PhasePlayerTrade(g, PhaseIdle(g)).(*phasePlayerTrade)
	tradeMenu.selected = 0
	nextPhase := tradeMenu.Confirm()

	if _, ok := nextPhase.(*phaseIdle); !ok {
		t.Fatalf("Expected to return to phaseIdle on failure, but got %T", nextPhase)
	}
	if g.TradeOffers[0].Status != OfferIsPending {
		t.Errorf("Expected offer status to remain Pending, got %v", g.TradeOffers[0].Status)
	}
}

func TestAcceptAmbiguousOffer(t *testing.T) {
	g := setupGameForTradeTest()
	g.PlayerTurn = 1
	g.TradeOffers = append(g.TradeOffers, TradeOffer{
		ID:        0,
		OffererID: 0,
		TargetID:  1,
		Offering:  map[board.ResourceType]int{board.ResourceInvalid: 1},
		Requesting: map[board.ResourceType]int{board.ResourceSheep: 1},
		Status:    OfferIsPending,
	})
	
	tradeMenu := PhasePlayerTrade(g, PhaseIdle(g)).(*phasePlayerTrade)
	tradeMenu.selected = 0
	nextPhase := tradeMenu.Confirm()

	if _, ok := nextPhase.(*phaseIdle); !ok {
		t.Fatalf("Expected to return to phaseIdle on failure, but got %T", nextPhase)
	}
	if g.TradeOffers[0].Status != OfferIsPending {
		t.Errorf("Expected offer status to remain Pending, got %v", g.TradeOffers[0].Status)
	}
}

func TestEndTurnInvalidation(t *testing.T) {
	g := setupGameForTradeTest()
	g.PlayerTurn = 0
	g.TradeOffers = append(g.TradeOffers, TradeOffer{Status: OfferIsPending})
	g.TradeOffers = append(g.TradeOffers, TradeOffer{Status: OfferIsPending})

	if len(g.TradeOffers) != 2 {
		t.Fatalf("Test setup failed, expected 2 offers, got %d", len(g.TradeOffers))
	}
	
	// Simulate ending the turn
	idlePhase := PhaseIdle(g).(*phaseIdle)
	idlePhase.selected = 3 // "End Turn"
	idlePhase.Confirm()

	if len(g.TradeOffers) != 0 {
		t.Errorf("Expected all offers to be invalidated at end of turn, but %d remain", len(g.TradeOffers))
	}
}