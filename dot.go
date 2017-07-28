package dreitafel

import (
	"fmt"
	"reflect"
	"sync"

	log "github.com/Sirupsen/logrus"
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
	print(``)
	print(`# horizontal layout`)
	print("rankdir=LR;")
	print("splines=ortho;")
	print("nodesep=0.8;")

	print("")
	print(`# Actors`)
	for _, node := range diagram.nodes {
		if reflect.TypeOf(node) == reflect.TypeOf(&Actor{}) {
			print(fmt.Sprintf("%v[shape=box];", node.(*Actor).title))
		}
	}
	print("")
	print(`# Storages`)
	for _, node := range diagram.nodes {
		if reflect.TypeOf(node) == reflect.TypeOf(&Storage{}) {
			print(fmt.Sprintf("%v[shape=box,style=rounded];", node.(*Storage).title))
		}
	}
	print("")
	print(`# Accesses`)
	for _, edge_ := range diagram.edges {
		var edgestr string
		edge := edge_.(*FmcBaseEdge)

		if edge.edgeType == EdgeTypeRead {
			log.Debug("Read access!")
			edgestr = fmt.Sprintf("%v -> %v;", edge.storage.title, edge.actor.title)
		} else {
			log.Debug("Write access!")
			edgestr = fmt.Sprintf("%v -> %v;", edge.actor.title, edge.storage.title)
		}

		print(edgestr)
	}
	fmt.Printf("} // end digraph\n")
}
