package main

import (
	"el_poblador/board"
	"el_poblador/game"
	"fmt"
)

func main() {
	fmt.Println("=== Action Log Demo ===")

	// Create a game
	g := &game.Game{}
	g.Start([]string{"Alice", "Bob", "Charlie"})

	// Simulate some actions to generate log entries

	// 1. Give Alice resources and buy a development card
	alice := &g.Players()[0]
	alice.AddResource(board.ResourceWheat)
	alice.AddResource(board.ResourceOre)
	alice.AddResource(board.ResourceSheep)

	// Log a development card purchase
	g.LogAction(fmt.Sprintf("%s bought a development card", alice.RenderName()))

	// 2. Log some building actions
	bob := &g.Players()[1]
	g.LogAction(fmt.Sprintf("%s built a road", bob.RenderName()))

	charlie := &g.Players()[2]
	g.LogAction(fmt.Sprintf("%s built a settlement", charlie.RenderName()))

	// 3. Log some more complex actions
	g.LogAction(fmt.Sprintf("%s played Knight", alice.RenderName()))
	g.LogAction(fmt.Sprintf("%s moved the robber", alice.RenderName()))
	g.LogAction(fmt.Sprintf("%s stole a card from %s", alice.RenderName(), bob.RenderName()))

	// 4. Log resource generation
	g.LogAction(fmt.Sprintf("%s gained 2 Madera, 1 Ladrillo from dice roll (8)", bob.RenderName()))

	// 5. Log turn change
	g.LogAction(fmt.Sprintf("Turn passed to %s", charlie.RenderName()))

	// Display the action log
	fmt.Println("\nAction Log:")
	fmt.Println("----------")
	for i, action := range g.ActionLog() {
		fmt.Printf("%2d. %s\n", i+1, action)
	}
}