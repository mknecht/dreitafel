// model for FMC block diagram
package dreitafel

type FmcBlockDiagram struct {
	title string
	nodes []FmcNode
	edges []FmcEdge
}

func (diagram *FmcBlockDiagram) String() string {
	return diagram.title
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

type EdgeType int

const (
	EdgeTypeRead EdgeType = iota
	EdgeTypeWrite
)

type FmcEdge interface {
}

// bipartite graph
type FmcBaseEdge struct {
	edgeType EdgeType
	actor    *Actor
	storage  *Storage
}
