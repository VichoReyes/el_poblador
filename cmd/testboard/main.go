package main

import (
	"el_poblador/board"
	"fmt"
)

func main() {
	fmt.Println()
	fmt.Println("=== Tile Rendering Test ===")

	// Test different tile types
	tiles := []board.Tile{
		{Terrain: board.TerrainWood, DiceNumber: 6},
		{Terrain: board.TerrainBrick, DiceNumber: 8},
		{Terrain: board.TerrainOre, DiceNumber: 12},
		{Terrain: board.TerrainWheat, DiceNumber: 4},
		{Terrain: board.TerrainSheep, DiceNumber: 10},
		{Terrain: board.TerrainDesert, DiceNumber: 0}, // Desert has no dice
	}

	for i, tile := range tiles {
		fmt.Printf("Tile %d (%s, dice: %d):\n", i+1, tile.Terrain.String(), tile.DiceNumber)
		rendered := tile.RenderTile(false)
		for i, line := range rendered {
			if i == 0 || i == len(rendered)-1 {
				fmt.Printf(" %s\n", line)
			} else {
				fmt.Printf("%s\n", line)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== Board Component Test ===")
	fmt.Println("Testing board printing functionality:")
	fmt.Println()

	board := board.NewChaoticBoard()
	lines := board.Print(nil)
	for _, line := range lines {
		fmt.Println(line)
	}
}
