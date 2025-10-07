package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Total args: %d\n", len(os.Args))
	fmt.Printf("Program name: %s\n\n", os.Args[0])

	if len(os.Args) > 1 {
		fmt.Println("Arguments:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("  [%d] %q (len=%d)\n", i, arg, len(arg))
		}
	} else {
		fmt.Println("No arguments provided")
	}
}
