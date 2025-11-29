package game

import (
	"el_poblador/board"
	"testing"
)

func isCrossAdjacentToTile(c board.CrossCoord, t board.TileCoord) bool {
	// replicate board.CrossCoord.adjacentTileCoords logic locally for testing
	var potentials []board.TileCoord
	if (c.X+c.Y)%2 == 0 {
		potentials = []board.TileCoord{
			{X: c.X - 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	} else {
		potentials = []board.TileCoord{
			{X: c.X, Y: c.Y},
			{X: c.X - 1, Y: c.Y - 1},
			{X: c.X - 1, Y: c.Y + 1},
		}
	}
	for _, p := range potentials {
		if candidate, ok := board.NewTileCoord(p.X, p.Y); ok {
			if candidate == t {
				return true
			}
		}
	}
	return false
}

// set up a simple scenario where current player plays Knight, places robber, and steals
func TestKnightPlayLeadsToRobberPlacementAndPotentialSteal(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	current := &g.Players[g.PlayerTurn]
	// Ensure current player has a Knight
	current.HiddenDevCards = append(current.HiddenDevCards, DevCardKnight)

	// Give next player some resources to be stealable
	victimIdx := (g.PlayerTurn + 1) % len(g.Players)
	g.Players[victimIdx].AddResource(board.ResourceWheat)
	g.Players[victimIdx].AddResource(board.ResourceBrick)

	// Also ensure victim has a settlement adjacent to a known tile we will pick
	// Choose a valid tile and one adjacent cross for victim
	tile, ok := board.NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid tile coordinate (2,3)")
	}
	// Find a cross adjacent to tile to assign to victim
	found := false
	for x := 0; x <= 5 && !found; x++ {
		for y := 0; y <= 10 && !found; y++ {
			c, ok := board.NewCrossCoord(x, y)
			if !ok {
				continue
			}
			if isCrossAdjacentToTile(c, tile) {
				if g.Board.CanPlaceSettlement(c) {
					g.Board.SetSettlement(c, victimIdx)
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatal("failed to assign victim settlement adjacent to chosen tile")
	}

	// Start from dice roll phase and select Play Knight (index 1)
	g.phase = PhaseDiceRoll(g)
	// Move selection down once to "Play Knight"
	g.MoveCursor("down", nil)
	g.ConfirmAction(nil)

	// We should now be in place-robber phase
	if _, ok := g.phase.(*phasePlaceRobber); !ok {
		t.Fatalf("expected phasePlaceRobber, got %T", g.phase)
	}

	// Directly set the tile cursor to our chosen tile and confirm placement
	pr := g.phase.(*phasePlaceRobber)
	pr.tileCoord = tile
	g.ConfirmAction(nil)

	// Since victim has resources and is adjacent, we should now be in steal phase
	steal, ok := g.phase.(*phaseStealCard)
	if !ok {
		t.Fatalf("expected phaseStealCard, got %T", g.phase)
	}

	// Select the victim if multiple players are available
	// Move cursor until selected player matches victimIdx; cap iterations
	for i := 0; i < len(steal.stealablePlayers)*2; i++ {
		if steal.stealablePlayers[steal.selected].Name == g.Players[victimIdx].Name {
			break
		}
		g.MoveCursor("down", nil)
	}

	beforeVictim := g.Players[victimIdx].TotalResources()
	beforeThief := current.TotalResources()

	g.ConfirmAction(nil)

	afterVictim := g.Players[victimIdx].TotalResources()
	afterThief := current.TotalResources()

	if !(afterVictim == beforeVictim-1 && afterThief == beforeThief+1) {
		t.Fatalf("expected thief +1 and victim -1 resources, got thief %d->%d, victim %d->%d", beforeThief, afterThief, beforeVictim, afterVictim)
	}
}

func TestMoveTileCursorHelpers(t *testing.T) {
	// Pick a valid starting tile
	// Pick a valid tile known from coordinates tests; ensure x+y is odd and within bounds
	start, ok := board.NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid starting tile")
	}

	// Right should move to (x+1, y+1)
	r, rok := moveTileCursor(start, "right")
	if !rok {
		t.Fatal("expected right move to be ok")
	}
	if r.X != start.X+1 || r.Y != start.Y+1 {
		t.Fatalf("right move unexpected: got %v from %v", r, start)
	}

	// Left from r should return to start
	l, lok := moveTileCursor(r, "left")
	if !lok || l != start {
		t.Fatalf("left move should return to start, got %v from %v", l, r)
	}

	// Up should decrease Y by 2
	u, uok := moveTileCursor(start, "up")
	if !uok {
		t.Fatal("expected up move to be ok")
	}
	if u.X != start.X || u.Y != start.Y-2 {
		t.Fatalf("up move unexpected: got %v from %v", u, start)
	}

	// Down should increase Y by 2 and return to start from u
	d, dok := moveTileCursor(u, "down")
	if !dok || d != start {
		t.Fatalf("down move should return to start, got %v from %v", d, u)
	}
}

func TestRobberCannotBePlacedOnSameTile(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	// Place robber on initial tile
	initialTile, ok := board.NewTileCoord(2, 3)
	if !ok {
		t.Fatal("expected valid tile coordinate (2,3)")
	}

	g.Board.PlaceRobber(initialTile)

	// Create a place robber phase and set cursor to the same tile
	phase := PhasePlaceRobber(g, PhaseIdle(g)).(*phasePlaceRobber)
	phase.tileCoord = initialTile

	// Try to confirm placement on same tile
	newPhase := phase.Confirm()

	// Should still be in place robber phase with error message
	if prPhase, ok := newPhase.(*phasePlaceRobber); !ok {
		t.Fatalf("expected to stay in phasePlaceRobber, got %T", newPhase)
	} else if prPhase.invalid == "" {
		t.Fatal("expected invalid message when trying to place robber on same tile")
	} else if prPhase.invalid != "Robber cannot be moved to the same tile it's already on" {
		t.Fatalf("expected specific error message, got: %s", prPhase.invalid)
	}

	// Verify robber didn't move
	currentPos := g.Board.GetRobber()
	if currentPos != initialTile {
		t.Fatalf("robber should not have moved from %v, but got %v", initialTile, currentPos)
	}
}
