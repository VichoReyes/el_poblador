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
	width, height := board.GetBoardSize()
	fmt.Printf("Board size: %dx%d\n", width, height)

	fmt.Println()
	fmt.Println("=== End Test ===")
}
