package game

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"
)

func TestSaveLoadRoundtrip(t *testing.T) {
	// Create and start a game
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Advance game state
	game.phase = PhaseDiceRoll(game)
	game.LogAction("Test action")

	// Save to file
	filename := "test_save.gob"
	defer os.Remove(filename)

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(game); err != nil {
		t.Fatalf("Failed to encode game: %v", err)
	}

	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Load from file
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var loadedGame Game
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&loadedGame); err != nil {
		t.Fatalf("Failed to decode game: %v", err)
	}

	// Restore phase
	loadedGame.phase = PhaseDiceRoll(&loadedGame)

	// Verify state
	if len(loadedGame.Players) != 3 {
		t.Errorf("Expected 3 players, got %d", len(loadedGame.Players))
	}

	// Verify players are preserved (names may be shuffled)
	playerNames := make(map[string]bool)
	for _, p := range loadedGame.Players {
		playerNames[p.Name] = true
	}
	expectedNames := []string{"Alice", "Bob", "Charlie"}
	for _, name := range expectedNames {
		if !playerNames[name] {
			t.Errorf("Expected to find player %s, but didn't", name)
		}
	}

	if len(loadedGame.ActionLog) == 0 {
		t.Error("Expected action log to be preserved")
	}

	if loadedGame.phase == nil {
		t.Error("Phase should be restored to PhaseDiceRoll")
	}
}
