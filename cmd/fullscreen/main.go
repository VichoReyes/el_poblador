package main

import (
	"el_poblador/game"
	"fmt"
)

func main() {
	g := &game.Game{}
	g.Start([]string{"Fred", "George", "Harold", "Ivor"})
	// 2 rounds of placing settlements and roads
	for i := 0; i < 2*4; i++ {
		g.MoveCursorToPlaceSettlement()
		g.ConfirmAction(nil) // place settlement
		g.ConfirmAction(nil) // place road
	}
	g.ConfirmAction(nil) // roll dice

	fmt.Println(g.Print(80, 32, nil, 0, 0))
}
