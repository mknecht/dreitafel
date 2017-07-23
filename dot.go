package dreitafel

import (
	"fmt"
	"reflect"
	"sync"
)

var header = `# Generated with Dreitafel
# https://github.com/mknecht/dreitafel
`

type Generator interface {
	Generate()
}

type DotGenerator struct {
}

func (DotGenerator) Generate(diagrams chan *FmcBlockDiagram, errors chan error, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var diagram *FmcBlockDiagram

	for diagram = <-diagrams; diagram != nil; diagram = <-diagrams {
		print := func(txt string) {
			fmt.Printf("        %v\n", txt)
		}
		fmt.Println(header)
		fmt.Printf("digraph \"%v\" {\n", diagram.title)
		print(`# horizontal layout`)
		print(`label="\G";`)
		print("rankdir=LR;")

		for _, node := range diagram.nodes {
			if reflect.TypeOf(node) == reflect.TypeOf(&Actor{}) {
				print(fmt.Sprintf("%v[shape=box];", node.(*Actor).title))
			}
		}
		for _, node := range diagram.nodes {
			if reflect.TypeOf(node) == reflect.TypeOf(&Storage{}) {
				print(fmt.Sprintf("%v[shape=box,style=rounded];", node.(*Storage).title))
			}
		}
		fmt.Printf("} // end digraph\n")
	}

}
