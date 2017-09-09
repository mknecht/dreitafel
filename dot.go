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
	for idx, edge := range diagram.edges {
		if edge.GetEdgeType() == EdgeTypeRead {
			log.Debug("Adding read access!")
			read := edge.(*BipartiteEdge)
			print(fmt.Sprintf("%v -> %v [arrowhead=vee];", read.storage.title, read.actor.title))
		} else if edge.GetEdgeType() == EdgeTypeWrite {
			write := edge.(*BipartiteEdge)
			log.Debug("Adding write access!")
			print(fmt.Sprintf("%v -> %v  [arrowhead=vee];", write.actor.title, write.storage.title))
		} else if edge.GetEdgeType() == EdgeTypeChannel {
			channel := edge.(*Channel)
			log.Debug("Adding channel!")
			print(fmt.Sprintf("ch%v[label=\"\", shape=circle, width=0.2]", idx))
			print(fmt.Sprintf("%v ->  ch%v [arrowhead=none];", channel.first.title, idx))
			print(fmt.Sprintf("ch%v ->  %v [arrowhead=none];", idx, channel.second.title))
		}

	}
	print("} // end digraph\n")
}
