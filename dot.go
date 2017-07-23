package dreitafel

import (
	"fmt"
	"reflect"
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
			if reflect.TypeOf(node) == reflect.TypeOf(&Actor{}) {
				fmt.Printf("• %v\n", node.(*Actor).title)
			}
		}
		fmt.Println("Storages: ")
		for _, node := range diagram.nodes {
			if reflect.TypeOf(node) == reflect.TypeOf(&Storage{}) {
				fmt.Printf("• %v\n", node.(*Storage).title)
			}
		}
	}

}
