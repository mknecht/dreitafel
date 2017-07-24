package dreitafel

import (
	"fmt"
	"reflect"
	"sync"
)

var header = `# Generated with Dreitafel
# https://github.com/mknecht/dreitafel
`

type DotGenerator interface {
	GenerateDot()
}

func GenerateDot(diagrams chan DotGenerator, errors chan error, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for diagram := <-diagrams; diagram != nil; diagram = <-diagrams {
		diagram.GenerateDot()
	}

}

func (diagram *FmcBlockDiagram) GenerateDot() {
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
