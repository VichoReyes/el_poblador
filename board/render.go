package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Blink(true)
		lines[midLine].WriteString(style.Render(" ○ "))
	} else if hasSettlement {
		_, isCity := board.cityUpgrades[coord]
		if isCity {
			lines[midLine].WriteString(board.playerRender(settlementOwner, "███"))
		} else {
			lines[midLine].WriteString(board.playerRender(settlementOwner, "▲▲▲"))
		}
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
