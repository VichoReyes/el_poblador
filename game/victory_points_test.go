package game

import (
	"el_poblador/board"
	"testing"
)

func TestPlayerVictoryPointsFromSettlements(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Player starts with 0 victory points
	if points := player.VictoryPoints(g); points != 0 {
		t.Fatalf("expected 0 victory points initially, got %d", points)
	}

	// Place a settlement manually
	coord, ok := board.NewCrossCoord(2, 4)
	if !ok {
		t.Fatal("expected valid cross coordinate")
	}
	g.Board.SetSettlement(coord, 0)

	// Should now have 1 victory point
	if points := player.VictoryPoints(g); points != 1 {
		t.Fatalf("expected 1 victory point after settlement, got %d", points)
	}
}

func TestPlayerVictoryPointsFromCities(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Place settlement and upgrade to city
	coord, ok := board.NewCrossCoord(2, 4)
	if !ok {
		t.Fatal("expected valid cross coordinate")
	}
	g.Board.SetSettlement(coord, 0)
	g.Board.UpgradeToCity(coord, 0)

	// Should have 2 victory points (1 from settlement + 1 additional from city upgrade = 2 total)
	if points := player.VictoryPoints(g); points != 2 {
		t.Fatalf("expected 2 victory points after city, got %d", points)
	}
}

func TestPlayerVictoryPointsFromDevelopmentCards(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Add victory point development cards
	player.HiddenDevCards = append(player.HiddenDevCards, DevCardVictoryPoint)
	player.HiddenDevCards = append(player.HiddenDevCards, DevCardVictoryPoint)

	// Should have 2 victory points from dev cards
	if points := player.VictoryPoints(g); points != 2 {
		t.Fatalf("expected 2 victory points from dev cards, got %d", points)
	}
}

func TestPlayerVisibleVictoryPointsHidesDevCards(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Add settlement (visible) and victory point dev card (hidden)
	coord, ok := board.NewCrossCoord(2, 4)
	if !ok {
		t.Fatal("expected valid cross coordinate")
	}
	g.Board.SetSettlement(coord, 0)
	player.HiddenDevCards = append(player.HiddenDevCards, DevCardVictoryPoint)

	// Total victory points should be 2
	if points := player.VictoryPoints(g); points != 2 {
		t.Fatalf("expected 2 total victory points, got %d", points)
	}
}

func TestPlayerVisibleVictoryPointsShowsPlayedDevCards(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Add victory point dev cards - one hidden, one played
	player.HiddenDevCards = append(player.HiddenDevCards, DevCardVictoryPoint)
	player.PlayedDevCards = append(player.PlayedDevCards, DevCardVictoryPoint)

	// Total victory points should be 2
	if points := player.VictoryPoints(g); points != 2 {
		t.Fatalf("expected 2 total victory points, got %d", points)
	}
}

func TestCheckGameEndReturnsNilWhenNoWinner(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	// No one has 10 points yet
	winner := g.CheckGameEnd()
	if winner != nil {
		t.Fatal("expected no winner initially")
	}

	// Add some points but not enough to win
	coord, ok := board.NewCrossCoord(2, 4)
	if !ok {
		t.Fatal("expected valid cross coordinate")
	}
	g.Board.SetSettlement(coord, 0)
	
	winner = g.CheckGameEnd()
	if winner != nil {
		t.Fatal("expected no winner with only 1 point")
	}
}

func TestCheckGameEndReturnsWinnerAt10Points(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	player := &g.Players[0]

	// Give player 10 victory points through dev cards
	for i := 0; i < 10; i++ {
		player.HiddenDevCards = append(player.HiddenDevCards, DevCardVictoryPoint)
	}

	winner := g.CheckGameEnd()
	if winner == nil {
		t.Fatal("expected winner with 10 points")
	}
	if winner != player {
		t.Fatal("expected first player to be the winner")
	}
}

func TestCheckGameEndReturnsFirstPlayerToReach10Points(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	// Give second player 10 points
	for i := 0; i < 10; i++ {
		g.Players[1].HiddenDevCards = append(g.Players[1].HiddenDevCards, DevCardVictoryPoint)
	}

	winner := g.CheckGameEnd()
	if winner == nil {
		t.Fatal("expected winner with 10 points")
	}
	if winner != &g.Players[1] {
		t.Fatal("expected second player to be the winner")
	}
}

func TestCountSettlementsAndCities(t *testing.T) {
	g := &Game{}
	g.Start([]string{"A", "B", "C"})

	// Place settlements for different players
	coord1, _ := board.NewCrossCoord(2, 4)
	coord2, _ := board.NewCrossCoord(2, 6)
	coord3, _ := board.NewCrossCoord(3, 5)
	
	g.Board.SetSettlement(coord1, 0) // Player 0
	g.Board.SetSettlement(coord2, 0) // Player 0
	g.Board.SetSettlement(coord3, 1) // Player 1

	// Upgrade one settlement to city
	g.Board.UpgradeToCity(coord1, 0)

	// Test settlement counting
	if count := g.Board.CountSettlements(0); count != 2 {
		t.Fatalf("expected player 0 to have 2 settlements, got %d", count)
	}
	if count := g.Board.CountSettlements(1); count != 1 {
		t.Fatalf("expected player 1 to have 1 settlement, got %d", count)
	}

	// Test city counting
	if count := g.Board.CountCities(0); count != 1 {
		t.Fatalf("expected player 0 to have 1 city, got %d", count)
	}
	if count := g.Board.CountCities(1); count != 0 {
		t.Fatalf("expected player 1 to have 0 cities, got %d", count)
	}
}
