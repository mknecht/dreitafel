package dreitafel

type Generator interface {
	Generate()
}

type DotGenerator struct {
}

func (DotGenerator) Generate(diagram *FmcBlockDiagram) (string, error) {
	return diagram.actor.Name, nil
}
