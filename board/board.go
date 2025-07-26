package board

// Board represents the game board
type Board struct {
	tiles map[TileCoord]Tile
	// roads and settlements are indexed by player id
	roads        map[PathCoord]int
	settlements  map[CrossCoord]int
	playerRender func(int, string) string
}

func (b *Board) ValidCrossCoord() CrossCoord {
	if coord, ok := NewCrossCoord(2, 4); ok {
		return coord
	}
	panic("valid cross coord wrong")
}

func (b *Board) CanPlaceSettlement(coord CrossCoord) bool {
	if _, ok := b.settlements[coord]; ok {
		return false
	}
	for _, neighbor := range coord.Neighbors() {
		if _, ok := b.settlements[neighbor]; ok {
			return false
		}
	}
	return true
}

func (b *Board) AdjacentTiles(coord CrossCoord) []Tile {
	tileCoords := coord.adjacentTileCoords()
	tiles := make([]Tile, len(tileCoords))
	for i, tileCoord := range tileCoords {
		tiles[i] = b.tiles[tileCoord]
	}
	return tiles
}

func (b *Board) SetSettlement(coord CrossCoord, playerId int) bool {
	if !b.CanPlaceSettlement(coord) {
		return false
	}
	b.settlements[coord] = playerId
	return true
}

func (b *Board) SetRoad(coord PathCoord, playerId int) {
	b.roads[coord] = playerId
}
