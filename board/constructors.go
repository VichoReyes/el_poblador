package board

import "math/rand/v2"

// NewDesertBoard creates a new board of only desert tiles
func NewDesertBoard() *Board {
	board := &Board{
		tiles:       make(map[TileCoord]Tile),
		roads:       make(map[PathCoord]int),
		settlements: make(map[CrossCoord]int),
		playerRender: func(_ int, content string) string {
			return content
		},
	}
	// brute force all tile coords
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			coord, valid := NewTileCoord(x, y)
			if valid {
				board.tiles[coord] = Tile{Terrain: Desierto, DiceNumber: 0}
			}
		}
	}
	return board
}

// NewChaoticBoard creates a new board with random tiles
func NewChaoticBoard() *Board {
	board := &Board{
		tiles:       make(map[TileCoord]Tile),
		roads:       make(map[PathCoord]int),
		settlements: make(map[CrossCoord]int),
		playerRender: func(_ int, content string) string {
			return content
		},
	}
	// brute force all tile coords
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			crossCoord, valid := NewCrossCoord(x, y)
			if valid && rand.IntN(4) == 0 {
				playerId := rand.IntN(4)
				board.settlements[crossCoord] = playerId
				neighbors := crossCoord.Neighbors()
				board.roads[NewPathCoord(crossCoord, neighbors[rand.IntN(len(neighbors))])] = playerId
			}
			tileCoord, valid := NewTileCoord(x, y)
			if valid {
				terrain := TerrainType(rand.IntN(6))
				dice := rand.IntN(11) + 2
				if terrain == Desierto {
					dice = 0
				}
				board.tiles[tileCoord] = Tile{Terrain: terrain, DiceNumber: dice}
			}
		}
	}
	return board
}

func NewLegalBoard(playerRender func(int, string) string) *Board {
	diceNumbers := []int{2, 3, 3, 4, 4, 5, 5, 6, 6, 8, 8, 9, 9, 10, 10, 11, 11, 12}
	terrains := []TerrainType{
		Bosque, Bosque, Bosque, Bosque,
		Arcilla, Arcilla, Arcilla,
		Montaña, Montaña, Montaña,
		Plantación, Plantación, Plantación, Plantación,
		Pasto, Pasto, Pasto, Pasto,
		Desierto,
	}
	rand.Shuffle(len(diceNumbers), func(i, j int) {
		diceNumbers[i], diceNumbers[j] = diceNumbers[j], diceNumbers[i]
	})
	rand.Shuffle(len(terrains), func(i, j int) {
		terrains[i], terrains[j] = terrains[j], terrains[i]
	})
	board := &Board{
		tiles:        make(map[TileCoord]Tile),
		roads:        make(map[PathCoord]int),
		settlements:  make(map[CrossCoord]int),
		playerRender: playerRender,
	}
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			tileCoord, valid := NewTileCoord(x, y)
			if !valid {
				continue
			}
			terr := terrains[0]
			terrains = terrains[1:]
			if terr == Desierto {
				board.tiles[tileCoord] = Tile{Terrain: terr, DiceNumber: 0}
			} else {
				board.tiles[tileCoord] = Tile{Terrain: terr, DiceNumber: diceNumbers[0]}
				diceNumbers = diceNumbers[1:]
			}
		}
	}
	return board
}
