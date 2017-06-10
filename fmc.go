package dreitafel

type FmcBlockDiagram struct {
	actor   *ActorNode
	storage *StorageNode

	edge *Edge
}

type Edge struct {
	actor   *ActorNode
	storage *StorageNode
}
