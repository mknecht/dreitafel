package main

import (
	"dreitafel"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: pattern to parse\n")
		os.Exit(2)
	}

	dreitafel.CompileToSvgToStdout(os.Args[1])
}
