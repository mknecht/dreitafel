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
	EdgeTypeUnknown EdgeType = iota
	EdgeTypeRead
	EdgeTypeWrite
	EdgeTypeChannel
)

type FmcEdge interface {
	GetEdgeType() EdgeType
}

// bipartite graph
type FmcBaseEdge struct {
	edgeType EdgeType
}

func (edge *FmcBaseEdge) GetEdgeType() EdgeType {
	return edge.edgeType
}

// bipartite graph
type BipartiteEdge struct {
	FmcBaseEdge
	actor   *Actor
	storage *Storage
}

// channel contains a storage: actor <-> actor
type Channel struct {
	FmcBaseEdge
	first  *Actor
	second *Actor
}
