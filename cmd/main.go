package main

import (
	"dreitafel"
	"fmt"
	"os"
)

func parse(input string) (*dreitafel.FmcBlockDiagram, error) {
	diagram, err := dreitafel.Parse(input)
	if err != nil {
		return nil, err
	}
	return diagram, nil
}

func toDot(diagram *dreitafel.FmcBlockDiagram) (string, error) {
	generator := dreitafel.DotGenerator{}
	return generator.Generate(diagram)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: pattern to parse\n")
		os.Exit(2)
	}

	diagram, err := parse(os.Args[1])
	fmt.Println("====PARSE RESULT====")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	} else {
		fmt.Println(diagram)
	}

	result, err := toDot(diagram)
	fmt.Println("====DOT RESULT====")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	} else {
		fmt.Println(result)
	}
}
