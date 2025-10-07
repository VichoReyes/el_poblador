package board

import "math/rand/v2"

// NewDesertBoard creates a new board of only desert tiles
func NewDesertBoard() *Board {
	board := &Board{
		Tiles:        make(map[TileCoord]Tile),
		Roads:        make(map[PathCoord]int),
		Settlements:  make(map[CrossCoord]int),
		CityUpgrades: make(map[CrossCoord]int),
		PlayerColors: make(map[int]int),
	}
	// brute force all tile coords
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			coord, valid := NewTileCoord(x, y)
			if valid {
				board.Tiles[coord] = Tile{Terrain: TerrainDesert, DiceNumber: 0}
			}
		}
	}
	return board
}

// NewChaoticBoard creates a new board with random tiles
func NewChaoticBoard() *Board {
	board := &Board{
		Tiles:        make(map[TileCoord]Tile),
		Roads:        make(map[PathCoord]int),
		Settlements:  make(map[CrossCoord]int),
		CityUpgrades: make(map[CrossCoord]int),
		PlayerColors: make(map[int]int),
	}
	// brute force all tile coords
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			crossCoord, valid := NewCrossCoord(x, y)
			if valid && rand.IntN(4) == 0 {
				playerId := rand.IntN(4)
				board.Settlements[crossCoord] = playerId
				neighbors := crossCoord.Neighbors()
				board.Roads[NewPathCoord(crossCoord, neighbors[rand.IntN(len(neighbors))])] = playerId
			}
			tileCoord, valid := NewTileCoord(x, y)
			if valid {
				terrain := TerrainType(rand.IntN(6))
				dice := rand.IntN(11) + 2
				if terrain == TerrainDesert {
					dice = 0
				}
				board.Tiles[tileCoord] = Tile{Terrain: terrain, DiceNumber: dice}
			}
		}
	}
	return board
}

func NewLegalBoard(playerColors map[int]int) *Board {
	diceNumbers := []int{2, 3, 3, 4, 4, 5, 5, 6, 6, 8, 8, 9, 9, 10, 10, 11, 11, 12}
	terrains := []TerrainType{
		TerrainWood, TerrainWood, TerrainWood, TerrainWood,
		TerrainBrick, TerrainBrick, TerrainBrick,
		TerrainOre, TerrainOre, TerrainOre,
		TerrainWheat, TerrainWheat, TerrainWheat, TerrainWheat,
		TerrainSheep, TerrainSheep, TerrainSheep, TerrainSheep,
		TerrainDesert,
	}
	rand.Shuffle(len(diceNumbers), func(i, j int) {
		diceNumbers[i], diceNumbers[j] = diceNumbers[j], diceNumbers[i]
	})
	rand.Shuffle(len(terrains), func(i, j int) {
		terrains[i], terrains[j] = terrains[j], terrains[i]
	})
	board := &Board{
		Tiles:        make(map[TileCoord]Tile),
		Roads:        make(map[PathCoord]int),
		Settlements:  make(map[CrossCoord]int),
		CityUpgrades: make(map[CrossCoord]int),
		PlayerColors: playerColors,
	}
	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			tileCoord, valid := NewTileCoord(x, y)
			if !valid {
				continue
			}
			terr := terrains[0]
			terrains = terrains[1:]
			if terr == TerrainDesert {
				board.Tiles[tileCoord] = Tile{Terrain: terr, DiceNumber: 0}
			} else {
				board.Tiles[tileCoord] = Tile{Terrain: terr, DiceNumber: diceNumbers[0]}
				diceNumbers = diceNumbers[1:]
			}
		}
	}
	return board
}
