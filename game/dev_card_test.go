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
}

func TestMonopolyCard(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give Alice a monopoly card
	game.players[0].hiddenDevCards = []DevCard{DevCardMonopoly}


	// Give players some resources  
	game.players[0].AddResource(board.ResourceWheat) // Alice (current player)
	game.players[1].AddResource(board.ResourceWheat) // Bob
	game.players[1].AddResource(board.ResourceWheat) // Bob (2 total)
	game.players[2].AddResource(board.ResourceWheat) // Charlie
	game.players[2].AddResource(board.ResourceOre)   // Charlie (different resource)

	// Create development card phase and select monopoly
	devCardPhase := PhasePlayDevelopmentCard(game, PhaseIdle(game))
	devCardPhaseImpl := devCardPhase.(*phasePlayDevelopmentCard)
	
	devCardPhaseImpl.selected = 0 // First card is monopoly

	// Execute dev card phase - should play monopoly card and transition to monopoly phase
	monopolyPhase := devCardPhaseImpl.Confirm()

	// Verify we got monopoly phase
	monopolyPhaseImpl, ok := monopolyPhase.(*phaseMonopoly)
	if !ok {
		t.Fatal("Should transition to monopoly phase")
	}

	// Record initial state
	initialAliceWheat := game.players[0].resources[board.ResourceWheat]
	initialBobWheat := game.players[1].resources[board.ResourceWheat]
	initialCharlieWheat := game.players[2].resources[board.ResourceWheat]

	// Find wheat option index
	wheatIndex := -1
	for i, resourceType := range board.RESOURCE_TYPES {
		if resourceType == board.ResourceWheat {
			wheatIndex = i
			break
		}
	}
	if wheatIndex == -1 {
		t.Fatal("Could not find wheat in resource types")
	}

	// Set selection to wheat and execute monopoly
	monopolyPhaseImpl.selected = wheatIndex
	nextPhase := monopolyPhaseImpl.Confirm()

	// Verify Alice collected all wheat
	expectedAliceWheat := initialAliceWheat + initialBobWheat + initialCharlieWheat
	if game.players[0].resources[board.ResourceWheat] != expectedAliceWheat {
		t.Errorf("Alice should have %d wheat, got %d", expectedAliceWheat, game.players[0].resources[board.ResourceWheat])
	}

	// Verify other players lost their wheat
	if game.players[1].resources[board.ResourceWheat] != 0 {
		t.Error("Bob should have no wheat after monopoly")
	}
	if game.players[2].resources[board.ResourceWheat] != 0 {
		t.Error("Charlie should have no wheat after monopoly")
	}

	// Verify the monopoly card was played
	if len(game.players[0].hiddenDevCards) != 0 {
		t.Errorf("Monopoly card should be removed from hidden cards, hidden count: %d, cards: %v", len(game.players[0].hiddenDevCards), game.players[0].hiddenDevCards)
	}
	if len(game.players[0].playedDevCards) != 1 {
		t.Errorf("Monopoly card should be added to played cards, played count: %d, cards: %v", len(game.players[0].playedDevCards), game.players[0].playedDevCards)
	}

	// Verify we return to idle phase with notification
	if _, ok := nextPhase.(*phaseIdle); !ok {
		t.Error("Monopoly should return to idle phase")
	}
}

func TestYearOfPlentyCard(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give Alice a year of plenty card
	game.players[0].hiddenDevCards = []DevCard{DevCardYearOfPlenty}

	// Record initial resource counts
	initialWheat := game.players[0].resources[board.ResourceWheat]
	initialOre := game.players[0].resources[board.ResourceOre]

	// Create development card phase and select year of plenty
	devCardPhase := PhasePlayDevelopmentCard(game, PhaseIdle(game))
	devCardPhaseImpl := devCardPhase.(*phasePlayDevelopmentCard)
	devCardPhaseImpl.selected = 0 // First card is year of plenty

	// Execute dev card phase - should play card and transition to year of plenty phase
	yearOfPlentyPhase := devCardPhaseImpl.Confirm()

	// Verify we got year of plenty phase
	yearOfPlentyPhaseImpl, ok := yearOfPlentyPhase.(*phaseYearOfPlenty)
	if !ok {
		t.Fatal("Should transition to year of plenty phase")
	}

	// Verify the card was played
	if len(game.players[0].hiddenDevCards) != 0 {
		t.Error("Year of Plenty card should be removed from hidden cards")
	}
	if len(game.players[0].playedDevCards) != 1 {
		t.Error("Year of Plenty card should be added to played cards")
	}

	// Find wheat and ore indices
	wheatIndex := -1
	oreIndex := -1
	for i, resourceType := range board.RESOURCE_TYPES {
		if resourceType == board.ResourceWheat {
			wheatIndex = i
		} else if resourceType == board.ResourceOre {
			oreIndex = i
		}
	}
	if wheatIndex == -1 || oreIndex == -1 {
		t.Fatal("Could not find wheat or ore in resource types")
	}

	// Select first resource (wheat)
	yearOfPlentyPhaseImpl.selected = wheatIndex
	phase2 := yearOfPlentyPhaseImpl.Confirm()

	// Should still be in year of plenty phase waiting for second resource
	yearOfPlentyPhaseImpl2, ok := phase2.(*phaseYearOfPlenty)
	if !ok {
		t.Fatal("Should still be in year of plenty phase after first selection")
	}

	// Verify help text shows what was selected
	helpText := yearOfPlentyPhaseImpl2.HelpText()
	if !strings.Contains(helpText, "Trigo") { // Wheat in Spanish
		t.Errorf("Help text should show selected wheat, got: %s", helpText)
	}

	// Select second resource (ore)
	yearOfPlentyPhaseImpl2.selected = oreIndex
	finalPhase := yearOfPlentyPhaseImpl2.Confirm()

	// Should transition to idle phase
	if _, ok := finalPhase.(*phaseIdle); !ok {
		t.Error("Should return to idle phase after selecting both resources")
	}

	// Verify player received both resources
	finalWheat := game.players[0].resources[board.ResourceWheat]
	finalOre := game.players[0].resources[board.ResourceOre]

	if finalWheat != initialWheat+1 {
		t.Errorf("Player should have gained 1 wheat, got %d, expected %d", finalWheat, initialWheat+1)
	}
	if finalOre != initialOre+1 {
		t.Errorf("Player should have gained 1 ore, got %d, expected %d", finalOre, initialOre+1)
	}
}

func TestRoadBuildingCard(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give Alice a road building card
	game.players[0].hiddenDevCards = []DevCard{DevCardRoadBuilding}

	// Create development card phase and select road building
	devCardPhase := PhasePlayDevelopmentCard(game, PhaseIdle(game))
	devCardPhaseImpl := devCardPhase.(*phasePlayDevelopmentCard)
	devCardPhaseImpl.selected = 0 // First card is road building

	// Execute dev card phase - should play card and transition to road building phase
	roadBuildingPhase := devCardPhaseImpl.Confirm()

	// Verify we got road start phase (for free road building)
	roadStartPhaseImpl, ok := roadBuildingPhase.(*phaseRoadStart)
	if !ok {
		t.Fatal("Should transition to road start phase for road building")
	}

	// Verify the card was played
	if len(game.players[0].hiddenDevCards) != 0 {
		t.Error("Road Building card should be removed from hidden cards")
	}
	if len(game.players[0].playedDevCards) != 1 {
		t.Error("Road Building card should be added to played cards")
	}

	// Verify help text indicates first free road
	helpText := roadStartPhaseImpl.HelpText()
	if !strings.Contains(helpText, "first") {
		t.Errorf("Help text should indicate first road, got: %s", helpText)
	}
	if !strings.Contains(helpText, "free") {
		t.Errorf("Help text should indicate free road, got: %s", helpText)
	}

	// Note: Full road placement testing would require setting up the board with valid positions
	// which is complex and would involve placing initial settlements first.
	// The basic card playing mechanism is tested above.
}

func TestBuyDevelopmentCardThroughBuildingPhase(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give Alice enough resources to buy a development card
	player := &game.players[0]
	player.AddResource(board.ResourceWheat)
	player.AddResource(board.ResourceOre)
	player.AddResource(board.ResourceSheep)

	// Record initial state
	initialDeckSize := len(game.devCardDeck)
	initialPlayerCards := player.TotalDevCards()

	// Create building phase
	buildingPhase := PhaseBuilding(game, PhaseIdle(game))
	buildingPhaseImpl := buildingPhase.(*phaseBuilding)

	// Select development card option (index 3)
	buildingPhaseImpl.selected = 3

	// Confirm the purchase
	nextPhase := buildingPhaseImpl.Confirm()

	// Should return to previous phase (idle)
	if _, ok := nextPhase.(*phaseIdle); !ok {
		t.Error("Should return to idle phase after buying development card")
	}

	// Check that the development card was actually added to the player
	finalPlayerCards := game.players[0].TotalDevCards() // Use game.players[0] to get the actual player
	if finalPlayerCards != initialPlayerCards+1 {
		t.Errorf("Player should have gained 1 development card, got %d, expected %d", finalPlayerCards, initialPlayerCards+1)
	}

	// Check that the deck size decreased
	finalDeckSize := len(game.devCardDeck)
	if finalDeckSize != initialDeckSize-1 {
		t.Errorf("Deck size should decrease by 1, got %d, expected %d", finalDeckSize, initialDeckSize-1)
	}

	// Check that resources were consumed
	if player.CanBuyDevelopmentCard() {
		t.Error("Player should not be able to buy another development card after consuming resources")
	}
}
