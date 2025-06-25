package board

import "fmt"

// PrintBoard prints the game board made of ASCII hexagons
func PrintBoard() {
	printHexagon()
	fmt.Println()
	printHexagonRow()
}

// printHexagon prints a single ASCII hexagon
func printHexagon() {
	fmt.Println("  /‾‾‾\\")
	fmt.Println(" /     \\")
	fmt.Println("|       |")
	fmt.Println(" \\     /")
	fmt.Println("  \\___/")
}

// printHexagonRow prints a row of connected hexagons (placeholder)
func printHexagonRow() {
	// TODO: Implement proper hexagon grid layout
	fmt.Println("TODO: Implement hexagon grid layout")
	fmt.Println("For now, just a single hexagon:")
	printHexagon()
}

// GetBoardSize returns the dimensions of the board
func GetBoardSize() (int, int) {
	// TODO: Return actual board dimensions
	return 5, 5 // placeholder
}
