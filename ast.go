// syntax tree for FMC block diagram source
package dreitafel

type TokenType int

const (
	TokenTypeActor TokenType = iota
	TokenTypeStorage
	TokenTypeRightAccess
	TokenTypeLineEnd
)

type Token interface {
	GetTokenType() TokenType
}

// In a pseudo-visual language,
// the tokens are the visual atoms,
// not the textual atoms.
// This means, an opening bracket is not regarded as meaningful,
// but the actor structure is, e.g. “[” is not, but “[Engine]” is.
type BaseToken struct {
	tokenType TokenType
	// we start with the horizontal case
	from int
	to   int
}

func (token *BaseToken) GetTokenType() TokenType {
	return token.tokenType
}

type Node struct {
	BaseToken
	title string
}
