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
	GenerateDot(print func(txt string))
}

func GenerateDot(diagrams chan DotGenerator, dot chan *string, errors chan error, wg *sync.WaitGroup) {
	defer close(dot)
	defer wg.Done()

	print := func(txt string) {
		dot <- &txt
	}

	for diagram := range diagrams {
		if diagram == nil {
			continue
		}
		log.Debugf("Got diagram: %v", diagram)
		diagram.GenerateDot(print)
		log.Debugf("Done generating diagram: %v", diagram)
	}
	log.Debug("Done generating diagrams")
}

func (diagram *FmcBlockDiagram) GenerateDot(print func(txt string)) {
	print(header)
	print(fmt.Sprintf("digraph \"%v\" {", diagram.title))
	print(``)
	print(`# horizontal layout`)
	print("rankdir=LR;")
	print("splines=ortho;")
	print("nodesep=0.8;")
	print("arrowhead=vee;")

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
			log.Debug("Adding read access!")
			edgestr = fmt.Sprintf("%v -> %v [arrowhead=vee];", edge.storage.title, edge.actor.title)
		} else {
			log.Debug("Adding write access!")
			edgestr = fmt.Sprintf("%v -> %v  [arrowhead=vee];", edge.actor.title, edge.storage.title)
		}

		print(edgestr)
	}
	print("} // end digraph\n")
}
