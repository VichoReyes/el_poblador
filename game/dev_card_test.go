package game

import (
	"el_poblador/board"
	"strings"
	"testing"
)

func TestDevelopmentCardPurchase(t *testing.T) {
	// Create a new game
	game := &Game{}
	game.Start([]string{"Player1", "Player2", "Player3"})

	// Give player 0 enough resources to buy a development card
	player := &game.players[0]
	player.AddResource(board.ResourceWheat)
	player.AddResource(board.ResourceOre)
	player.AddResource(board.ResourceSheep)

	// Check initial state
	initialDeckSize := len(game.devCardDeck)
	initialPlayerCards := player.TotalDevCards()

	// Verify player can buy a development card
	if !player.CanBuyDevelopmentCard() {
		t.Error("Player should be able to buy development card")
	}

	// Buy a development card
	if !player.BuyDevelopmentCard() {
		t.Error("Player should be able to consume resources for development card")
	}

	// Draw a card from the game deck
	card := game.DrawDevelopmentCard()
	if card == nil {
		t.Error("Should be able to draw a card from the deck")
	}

	// Add card to player's hidden deck
	player.hiddenDevCards = append(player.hiddenDevCards, *card)

	// Verify final state
	finalDeckSize := len(game.devCardDeck)
	finalPlayerCards := player.TotalDevCards()

	if finalDeckSize != initialDeckSize-1 {
		t.Errorf("Deck size should decrease by 1, got %d, expected %d", finalDeckSize, initialDeckSize-1)
	}

	if finalPlayerCards != initialPlayerCards+1 {
		t.Errorf("Player cards should increase by 1, got %d, expected %d", finalPlayerCards, initialPlayerCards+1)
	}

	// Verify player no longer has resources to buy another card
	if player.CanBuyDevelopmentCard() {
		t.Error("Player should not be able to buy another development card")
	}
}

func TestDevelopmentCardPlay(t *testing.T) {
	// Create a new game
	game := &Game{}
	game.Start([]string{"Player1", "Player2", "Player3"})

	// Give player 0 a development card
	player := &game.players[0]
	player.hiddenDevCards = []DevCard{DevCardKnight}

	// Check initial state
	initialHidden := len(player.hiddenDevCards)
	initialPlayed := len(player.playedDevCards)

	// Play the card
	if !player.PlayDevCard(DevCardKnight) {
		t.Error("Should be able to play development card")
	}

	// Verify final state
	finalHidden := len(player.hiddenDevCards)
	finalPlayed := len(player.playedDevCards)

	if finalHidden != initialHidden-1 {
		t.Errorf("Hidden cards should decrease by 1, got %d, expected %d", finalHidden, initialHidden-1)
	}

	if finalPlayed != initialPlayed+1 {
		t.Errorf("Played cards should increase by 1, got %d, expected %d", finalPlayed, initialPlayed+1)
	}

	// Try to play a card that's not in hidden deck
	if player.PlayDevCard(DevCardKnight) {
		t.Error("Should not be able to play a card that's not in hidden deck")
	}
}

func TestKnightCardUsage(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give the first player a knight card
	player := &game.players[0]
	player.hiddenDevCards = append(player.hiddenDevCards, DevCardKnight)

	// Test that the player has a knight card
	if !player.HasKnightCard() {
		t.Fatal("Player should have a knight card")
	}

	// Test that the dice roll phase always shows the knight option
	dicePhase := PhaseDiceRoll(game)
	if phaseWithMenu, ok := dicePhase.(PhaseWithMenu); ok {
		menu := phaseWithMenu.Menu()
		if !strings.Contains(menu, "Play Knight") {
			t.Fatal("Dice roll phase should always show Play Knight option")
		}
	}

	// Test that the idle phase shows the development card option
	idlePhase := PhaseIdle(game)
	if phaseWithMenu, ok := idlePhase.(PhaseWithMenu); ok {
		menu := phaseWithMenu.Menu()
		if !strings.Contains(menu, "Play Development Card") {
			t.Fatal("Idle phase should show Play Development Card option")
		}
	}

	// Test that the development card phase shows knight options when available
	devCardPhase := PhasePlayDevelopmentCard(game, idlePhase)
	if phaseWithMenu, ok := devCardPhase.(PhaseWithMenu); ok {
		menu := phaseWithMenu.Menu()
		if !strings.Contains(menu, "Knight") {
			t.Fatal("Development card phase should show Knight option when player has knight cards")
		}
	}

	// Test that the knight card is consumed when played
	initialCardCount := len(player.hiddenDevCards)
	if !player.PlayDevCard(DevCardKnight) {
		t.Fatal("Should be able to play knight card")
	}

	if len(player.hiddenDevCards) != initialCardCount-1 {
		t.Fatal("Knight card should be removed from hidden cards when played")
	}

	if len(player.playedDevCards) != 1 {
		t.Fatal("Knight card should be added to played cards when played")
	}

	// Test that the player no longer has knight cards
	if player.HasKnightCard() {
		t.Fatal("Player should not have knight cards after playing one")
	}
}
