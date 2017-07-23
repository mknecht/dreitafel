// model for FMC block diagram
package dreitafel

type FmcBlockDiagram struct {
	title string
	nodes []FmcNode
	edges []FmcEdge
}

type FmcNode interface {
}

type FmcBaseNode struct {
	id    string
	title string
}

type Actor struct {
	FmcBaseNode
}

type Storage struct {
	FmcBaseNode
}

type FmcEdge interface {
}

// bipartite graph
type FmcBaseEdge struct {
	actor   *Actor
	storage *Storage
}
