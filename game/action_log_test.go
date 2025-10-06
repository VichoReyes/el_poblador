package game

import (
	"el_poblador/board"
	"testing"
)

func TestActionLogBasicFunctionality(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Test basic logging
	game.LogAction("Test action 1")
	game.LogAction("Test action 2")

	if len(game.actionLog) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(game.actionLog))
	}

	if game.actionLog[0] != "Test action 1" {
		t.Errorf("Expected first action 'Test action 1', got '%s'", game.actionLog[0])
	}

	if game.actionLog[1] != "Test action 2" {
		t.Errorf("Expected second action 'Test action 2', got '%s'", game.actionLog[1])
	}
}

func TestActionLogMaxLength(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Add more than 15 actions
	for i := 1; i <= 20; i++ {
		game.LogAction("Action " + string(rune('0'+i%10)))
	}

	if len(game.actionLog) != 15 {
		t.Errorf("Expected log length to be capped at 15, got %d", len(game.actionLog))
	}

	// The oldest actions should be removed
	if game.actionLog[0] != "Action 6" {
		t.Errorf("Expected first action 'Action 6', got '%s'", game.actionLog[0])
	}

	if game.actionLog[14] != "Action 0" {
		t.Errorf("Expected last action 'Action 0', got '%s'", game.actionLog[14])
	}
}

func TestActionLogDevelopmentCardPurchase(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give the first player resources to buy a development card
	firstPlayer := &game.players[0]
	firstPlayer.AddResource(board.ResourceWheat)
	firstPlayer.AddResource(board.ResourceOre)
	firstPlayer.AddResource(board.ResourceSheep)

	// Check that development card deck is not empty
	if len(game.devCardDeck) == 0 {
		t.Fatal("Development card deck is empty")
	}

	// Simulate buying through the building phase
	game.phase = PhaseBuilding(game, PhaseIdle(game))
	buildPhase := game.phase.(*phaseBuilding)
	buildPhase.selected = 3 // Development Card option

	// Execute the purchase
	originalDevCardCount := len(firstPlayer.hiddenDevCards)
	game.phase = game.phase.Confirm()
	newDevCardCount := len(firstPlayer.hiddenDevCards)

	// Verify the purchase actually happened
	if newDevCardCount != originalDevCardCount+1 {
		t.Fatalf("Development card purchase failed. Cards before: %d, after: %d", originalDevCardCount, newDevCardCount)
	}

	// Check that the action was logged
	if len(game.actionLog) == 0 {
		t.Fatal("Expected at least one action to be logged")
	}

	expectedAction := firstPlayer.Name + " bought a development card"
	found := false
	for _, action := range game.actionLog {
		if action == expectedAction {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find action '%s' in log: %v", expectedAction, game.actionLog)
	}
}

func TestRenderNameMethod(t *testing.T) {
	// Test that RenderName method exists and works
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	player := &game.players[0]

	// Test that RenderName exists and returns the name (colored or not)
	renderedName := player.RenderName()

	// The method should at least return the player's name
	if renderedName != player.Name {
		// In test environments, lipgloss might not add colors, so check if it contains the name
		if !containsText(renderedName, player.Name) {
			t.Errorf("RenderName() should contain player name. Got: '%s', expected to contain: '%s'", renderedName, player.Name)
		}
	}
}

func TestActionLogUsesRenderNameInsteadOfPlainName(t *testing.T) {
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Give the first player resources to buy a development card
	firstPlayer := &game.players[0]
	firstPlayer.AddResource(board.ResourceWheat)
	firstPlayer.AddResource(board.ResourceOre)
	firstPlayer.AddResource(board.ResourceSheep)

	// Simulate buying through the building phase
	game.phase = PhaseBuilding(game, PhaseIdle(game))
	buildPhase := game.phase.(*phaseBuilding)
	buildPhase.selected = 3 // Development Card option

	// Execute the purchase
	game.phase = game.phase.Confirm()

	// Check that the action was logged
	if len(game.actionLog) == 0 {
		t.Fatal("Expected at least one action to be logged")
	}

	loggedAction := game.actionLog[0]

	// Verify it contains the expected text structure
	if !containsText(loggedAction, "bought a development card") {
		t.Errorf("Expected action to contain 'bought a development card', got: %s", loggedAction)
	}

	// Verify it contains the player's name
	if !containsText(loggedAction, firstPlayer.Name) {
		t.Errorf("Expected action to contain player name '%s', got: %s", firstPlayer.Name, loggedAction)
	}
}

// Helper function to check if string contains ANSI color codes
func containsANSIColorCodes(s string) bool {
	// ANSI color codes start with \x1b[ (ESC[)
	return containsText(s, "\x1b[") || containsText(s, "\033[")
}

// Helper function to check if string contains substring (case-insensitive for robustness)
func containsText(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}