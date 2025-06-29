package main

import (
	"el_poblador/board"
	"fmt"
)

func main() {
	fmt.Println("=== Board Component Test ===")
	fmt.Println("Testing board printing functionality:")
	fmt.Println()

	board.PrintBoard()

	fmt.Println()
	fmt.Println("=== Tile Rendering Test ===")

	// Test different tile types
	tiles := []board.Tile{
		{Terrain: board.Bosque, DiceNumber: 6},
		{Terrain: board.Arcilla, DiceNumber: 8},
		{Terrain: board.Montaña, DiceNumber: 12},
		{Terrain: board.Plantación, DiceNumber: 4},
		{Terrain: board.Pasto, DiceNumber: 10},
		{Terrain: board.Desierto, DiceNumber: 0}, // Desert has no dice
	}

	for i, tile := range tiles {
		fmt.Printf("Tile %d (%s, dice: %d):\n", i+1, tile.Terrain.String(), tile.DiceNumber)
		rendered := tile.RenderTile()
		for i, line := range rendered {
			if i == 0 || i == 3 {
				fmt.Printf(" %s\n", line)
			} else {
				fmt.Printf("%s\n", line)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== End Test ===")
}
