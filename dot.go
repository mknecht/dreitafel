package dreitafel

import (
	"fmt"
	"sync"
)

type Generator interface {
	Generate()
}

type DotGenerator struct {
}

func (DotGenerator) Generate(diagrams chan *FmcBlockDiagram, errors chan error, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var diagram *FmcBlockDiagram

	for diagram = <-diagrams; diagram != nil; diagram = <-diagrams {
		fmt.Println(diagram.title)
		fmt.Println(" === ")
		fmt.Println("Actors: ")
		for _, node := range diagram.nodes {
			fmt.Printf("• %v\n", node.(*Actor).title)
		}
	}

}
