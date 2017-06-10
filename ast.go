package dreitafel

type AstNode struct {
	xpos int
}

type ActorNode struct {
	AstNode
	Name string
}

type StorageNode struct {
	AstNode
	Name string
}
