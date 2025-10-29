package game

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestGobSerialization(t *testing.T) {
	// Create a new game
	game := &Game{}
	game.Start([]string{"Alice", "Bob", "Charlie"})

	// Advance the game a bit (similar to fullscreen example)
	// Place first settlement and road for Alice
	game.MoveCursorToPlaceSettlement()
	game.ConfirmAction(nil)
	game.ConfirmAction(nil) // Confirm road placement

	// Place first settlement and road for Bob
	game.MoveCursorToPlaceSettlement()
	game.ConfirmAction(nil)
	game.ConfirmAction(nil)

	// Place first settlement and road for Charlie
	game.MoveCursorToPlaceSettlement()
	game.ConfirmAction(nil)
	game.ConfirmAction(nil)

	// Serialize the game to gob
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(game)
	if err != nil {
		t.Fatalf("Failed to encode game to gob: %v", err)
	}

	// Verify we got some data
	gobData := buf.Bytes()
	if len(gobData) == 0 {
		t.Fatal("Gob data is empty")
	}

	// Deserialize the gob back into a game struct
	var deserializedGame Game
	decoder := gob.NewDecoder(bytes.NewReader(gobData))
	err = decoder.Decode(&deserializedGame)
	if err != nil {
		t.Fatalf("Failed to decode gob to game: %v", err)
	}

	// Verify the deserialized game has the same data
	if len(deserializedGame.Players) != len(game.Players) {
		t.Errorf("Player count mismatch: got %d, want %d", len(deserializedGame.Players), len(game.Players))
	}

	for i, player := range game.Players {
		deserializedPlayer := deserializedGame.Players[i]
		if deserializedPlayer.Name != player.Name {
			t.Errorf("Player %d name mismatch: got %s, want %s", i, deserializedPlayer.Name, player.Name)
		}
		if deserializedPlayer.Color != player.Color {
			t.Errorf("Player %d color mismatch: got %+v, want %+v", i, deserializedPlayer.Color, player.Color)
		}
	}

	// Verify board state
	if deserializedGame.Board == nil {
		t.Fatal("Deserialized board is nil")
	}

	// Check that settlements were preserved
	if len(deserializedGame.Board.Settlements) != len(game.Board.Settlements) {
		t.Errorf("Settlement count mismatch: got %d, want %d", len(deserializedGame.Board.Settlements), len(game.Board.Settlements))
	}

	// Check that roads were preserved
	if len(deserializedGame.Board.Roads) != len(game.Board.Roads) {
		t.Errorf("Road count mismatch: got %d, want %d", len(deserializedGame.Board.Roads), len(game.Board.Roads))
	}

	// Verify player turn
	if deserializedGame.PlayerTurn != game.PlayerTurn {
		t.Errorf("PlayerTurn mismatch: got %d, want %d", deserializedGame.PlayerTurn, game.PlayerTurn)
	}

	// Verify dev card deck
	if len(deserializedGame.DevCardDeck) != len(game.DevCardDeck) {
		t.Errorf("DevCardDeck count mismatch: got %d, want %d", len(deserializedGame.DevCardDeck), len(game.DevCardDeck))
	}

	// Verify action log
	if len(deserializedGame.ActionLog) != len(game.ActionLog) {
		t.Errorf("ActionLog count mismatch: got %d, want %d", len(deserializedGame.ActionLog), len(game.ActionLog))
	}

	t.Logf("Successfully serialized and deserialized game with %d bytes of gob data", len(gobData))
}
