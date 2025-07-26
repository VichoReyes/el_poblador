package board

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
)

// TerrainType represents the different types of terrain in Catan
type TerrainType int

const (
	Bosque TerrainType = iota
	Arcilla
	Montaña
	Plantación
	Pasto
	Desierto
)

// String returns the string representation of the terrain type
func (t TerrainType) String() string {
	switch t {
	case Bosque:
		return "Bosque"
	case Arcilla:
		return "Arcilla"
	case Montaña:
		return "Montaña"
	case Plantación:
		return "Plantación"
	case Pasto:
		return "Pasto"
	case Desierto:
		return "Desierto"
	default:
		return "Desconocido"
	}
}

// Tile represents a hexagonal tile on the Catan board
type Tile struct {
	Terrain    TerrainType
	DiceNumber int // 2-12, 0 for desert (no dice)
}

// RenderTile returns a 5-element array of strings representing the tile
// looks like this:
//
//		/‾‾‾‾‾‾\
//		/  B  \
//	    2
//		\      /
//		\______/
//
// where B is the terrain abbreviation and 2 is the dice number
func (tile *Tile) RenderTile(isCursor bool) [5]string {
	terrainAbbrev := tile.getTerrainAbbrev()

	diceStr := ""
	if tile.DiceNumber > 0 {
		diceStr = fmt.Sprintf("%d", tile.DiceNumber)
	}

	code := tile.getTerrainColor()
	if isCursor {
		code = fmt.Sprintf("%s\033[48;5;15m", code)
	}
	endCode := "\033[0m"
	lines := [5]string{
		fmt.Sprintf("%s/‾‾‾‾‾‾\\%s", code, endCode),
		fmt.Sprintf("%s/  %s  \\%s", code, terrainAbbrev, endCode),
		fmt.Sprintf("%s    %2s    %s", code, diceStr, endCode),
		fmt.Sprintf("%s\\        /%s", code, endCode),
		fmt.Sprintf("%s\\______/%s", code, endCode),
	}
	return lines
}

// getTerrainAbbrev returns a short abbreviation for the terrain type
func (tile *Tile) getTerrainAbbrev() string {
	switch tile.Terrain {
	case Bosque:
		return "BOSQ"
	case Arcilla:
		return "ARCI"
	case Montaña:
		return "MONT"
	case Plantación:
		return "PLAN"
	case Pasto:
		return "PAST"
	case Desierto:
		return "DESI"
	default:
		return "????"
	}
}

// getTerrainColor returns the terminal 256 color code for the terrain color
func (tile *Tile) getTerrainColor() string {
	var colorNumber int
	switch tile.Terrain {
	case Bosque:
		// dark green
		colorNumber = 2
	case Arcilla:
		// red
		colorNumber = 1
	case Montaña:
		// dark gray
		colorNumber = 8
	case Plantación:
		// yellow
		colorNumber = 11
	case Pasto:
		// green
		colorNumber = 10
	case Desierto:
		// brown
		colorNumber = 3
	default:
		colorNumber = 0
	}
	return fmt.Sprintf("\033[38;5;%dm", colorNumber)
}

// Board represents the game board
type Board struct {
	tiles map[TileCoord]Tile
	// roads and settlements are indexed by player id
	roads        map[PathCoord]int
	settlements  map[CrossCoord]int
	playerRender func(int, string) string
}

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

func (b *Board) ValidCrossCoord() CrossCoord {
	if coord, ok := NewCrossCoord(2, 4); ok {
		return coord
	}
	panic("valid cross coord wrong")
}

// PrintBoard prints the game board made of ASCII hexagons
func (b *Board) Print(cursor interface{}) []string {
	// there will be 31 lines (5 * 5 + 6 for the roads)
	lines := make([]strings.Builder, 31)
	sidePadding(lines)

	for x := 0; x <= 5; x++ {
		for y := 0; y <= 10; y++ {
			coord, valid := NewCrossCoord(x, y)
			if valid {
				renderCrossing(b, lines, coord, cursor)
			}
		}
	}

	sidePadding(lines)

	renderedLines := []string{}
	for i := range lines {
		renderedLines = append(renderedLines, lines[i].String())
	}
	return renderedLines
}

// takes responsibility for the crossing and whatever is to its right
// right-up and right-down paths
func renderCrossing(board *Board, lines []strings.Builder, coord CrossCoord, cursor interface{}) {
	// print crossing
	midLine := coord.Y * 3
	settlementOwner, hasSettlement := board.settlements[coord]
	hasCursor := false
	if c, ok := cursor.(CrossCoord); ok && c == coord {
		hasCursor = true
	}
	if hasCursor {
		lines[midLine].WriteString("\033[5;33m ○ \033[0m")
	} else if hasSettlement {
		lines[midLine].WriteString(board.playerRender(settlementOwner, "▲▲▲"))
	} else {
		lines[midLine].WriteString("   ")
	}

	// print right side
	if (coord.X+coord.Y)%2 == 1 {
		up, valid := coord.Up()
		if valid {
			path := NewPathCoord(coord, up)
			roadOwner, hasRoad := board.roads[path]
			if hasRoad {
				lines[midLine-2].WriteString(board.playerRender(roadOwner, "//"))
				lines[midLine-1].WriteString(board.playerRender(roadOwner, "//"))
			} else {
				lines[midLine-2].WriteString("  ")
				lines[midLine-1].WriteString("  ")
			}
		}
		down, valid := coord.Down()
		if valid {
			path := NewPathCoord(coord, down)
			roadOwner, hasRoad := board.roads[path]
			if hasRoad {
				lines[midLine+1].WriteString(board.playerRender(roadOwner, "\\\\"))
				lines[midLine+2].WriteString(board.playerRender(roadOwner, "\\\\"))
			} else {
				lines[midLine+1].WriteString("  ")
				lines[midLine+2].WriteString("  ")
			}
		}
		tileCoord, valid := NewTileCoord(coord.X, coord.Y)
		if valid {
			hasCursor := false
			if c, ok := cursor.(TileCoord); ok && c == tileCoord {
				hasCursor = true
			}
			tile := board.tiles[tileCoord]
			renderedTile := tile.RenderTile(hasCursor)
			lines[midLine-2].WriteString(renderedTile[0])
			lines[midLine-1].WriteString(renderedTile[1])
			lines[midLine].WriteString(renderedTile[2])
			lines[midLine+1].WriteString(renderedTile[3])
			lines[midLine+2].WriteString(renderedTile[4])
		}
	} else {
		right, valid := coord.Right()
		if !valid {
			return
		}
		pathCoord := NewPathCoord(coord, right)
		roadOwner, hasRoad := board.roads[pathCoord]
		if hasRoad {
			lines[midLine].WriteString(board.playerRender(roadOwner, " ==== "))
		} else {
			lines[midLine].WriteString("      ")
		}
	}
}

func sidePadding(lines []strings.Builder) {
	// fake paths, tiles and crossings spaces
	top := []int{3 + 6 + 3 + 10, 2 + 8 + 2 + 10, 2 + 10 + 2 + 8, 3 + 10, 2 + 10, 2 + 8}
	// left padding with virtual tiles would go
	// 2, 2, 1, 0, 1, 2, repeat
	pattern := []int{2, 2, 1, 0, 1, 2}
	for i := range lines {
		base := pattern[i%len(pattern)]
		if i < len(top) {
			base += top[i]
		}
		if len(lines)-i-1 < len(top) {
			base += top[len(lines)-i-1]
		}
		lines[i].WriteString(strings.Repeat(" ", base))
	}

}

func terminalLength(s string) int {
	n := 0
	skip := false
	for _, c := range s {
		if c == '\033' {
			skip = true
		} else if c == 'm' {
			skip = false
		} else if !skip {
			n++
		}
	}
	return n
}

// CrossCoord represents the coordinates of an intersection point
// where three hexagons meet.
type CrossCoord struct {
	X, Y int
}

// TileCoord represents the coordinates of a hexagonal tile.
type TileCoord struct {
	X, Y int
}

// PathCoord represents a path/edge between two intersections.
type PathCoord struct {
	From, To CrossCoord
}

// NewCrossCoord creates a intersection coordinate and returns whether it is valid
func NewCrossCoord(x, y int) (CrossCoord, bool) {
	coord := CrossCoord{X: x, Y: y}
	if coord.IsInBounds() {
		return coord, true
	}
	return CrossCoord{}, false
}

// IsInBounds checks if the coordinate is within the bounds of the board
// currently only for 3-4 player games
func (c CrossCoord) IsInBounds() bool {
	x := c.X
	y := c.Y
	// left edge
	if x < 0 {
		return false
	}
	// top left edge
	if x+y < 2 {
		return false
	}
	// top right edge
	if x-y > 3 {
		return false
	}
	// right edge
	if x > 5 {
		return false
	}
	// bottom right edge
	if x+y > 13 {
		return false
	}
	// bottom left edge
	if y-x > 8 {
		return false
	}
	return true
}

func (c CrossCoord) Up() (CrossCoord, bool) {
	return NewCrossCoord(c.X, c.Y-1)
}

func (c CrossCoord) Down() (CrossCoord, bool) {
	return NewCrossCoord(c.X, c.Y+1)
}

func (c CrossCoord) Left() (CrossCoord, bool) {
	if (c.X+c.Y)%2 == 0 {
		return CrossCoord{}, false
	}
	return NewCrossCoord(c.X-1, c.Y)
}

func (c CrossCoord) Right() (CrossCoord, bool) {
	if (c.X+c.Y)%2 == 1 {
		return CrossCoord{}, false
	}
	return NewCrossCoord(c.X+1, c.Y)
}

func (c CrossCoord) Neighbors() []CrossCoord {
	// TODO: use Up, Down, Left, Right methods
	var potential []CrossCoord
	if (c.X+c.Y)%2 == 0 {
		potential = []CrossCoord{
			{X: c.X + 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	} else {
		potential = []CrossCoord{
			{X: c.X - 1, Y: c.Y},
			{X: c.X, Y: c.Y - 1},
			{X: c.X, Y: c.Y + 1},
		}
	}
	neighbors := []CrossCoord{}
	for _, p := range potential {
		if p.IsInBounds() {
			neighbors = append(neighbors, p)
		}
	}
	return neighbors
}

// NewTileCoord creates a new tile coordinate and returns whether it is valid
func NewTileCoord(x, y int) (TileCoord, bool) {
	if (x+y)%2 == 0 {
		return TileCoord{}, false
	}
	_, leftOk := NewCrossCoord(x, y)
	_, acrossOk := NewCrossCoord(x+1, y)
	if leftOk && acrossOk {
		return TileCoord{X: x, Y: y}, true
	}
	return TileCoord{}, false
}

// NewPathCoord creates a new path coordinate between two intersections
// panics if from and to are not neighbors
func NewPathCoord(from, to CrossCoord) PathCoord {
	fromNeighbors := from.Neighbors()
	if !slices.Contains(fromNeighbors, to) {
		panic("from and to are not neighbors")
	}
	// Ensure canonical ordering (ascending)
	if (from.X < to.X) || (from.X == to.X && from.Y < to.Y) {
		return PathCoord{From: from, To: to}
	}
	return PathCoord{From: to, To: from}
}

// String returns the string representation of an intersection coordinate
func (ic CrossCoord) String() string {
	return fmt.Sprintf("(%d,%d)", ic.X, ic.Y)
}

// String returns the string representation of a tile coordinate
func (tc TileCoord) String() string {
	return fmt.Sprintf("[%d,%d]", tc.X, tc.Y)
}

// String returns the string representation of a path coordinate
func (pc PathCoord) String() string {
	return fmt.Sprintf("%s-%s", pc.From.String(), pc.To.String())
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
