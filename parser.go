package dreitafel

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

type TokenType int

const (
	ActorToken TokenType = iota
)

type Parser struct {
	source string
	pos    int
	col    int
	row    int
}

type Expression struct {
}

func Parse(source string) (*FmcBlockDiagram, error) {
	var err error

	diagram := &FmcBlockDiagram{}

	parser := Parser{}
	parser.row = 1
	parser.source = source

	parser.consumeWhitespaceIfAny()
	diagram.actor, err = parser.consumeActor()
	if err != nil {
		return nil, err
	}

	parser.consumeWhitespaceIfAny()
	err = parser.consumeRightwardsAccess()
	if err != nil {
		return nil, err
	}
	diagram.edge = &Edge{actor: diagram.actor, storage: diagram.storage}

	parser.consumeWhitespaceIfAny()
	diagram.storage, err = parser.consumeStorage()
	if err != nil {
		return nil, err
	}
	parser.consumeWhitespaceIfAny()
	if err != nil {
		return nil, err
	}
	return diagram, nil
}

func (parser *Parser) consumeStorage() (*StorageNode, error) {
	storage := StorageNode{}
	storage.xpos = parser.col

	err := parser.consumeChar('(')
	if err != nil {
		return nil, err
	}
	parser.consumeWhitespaceIfAny()

	storage.Name, err = parser.consumeName()
	if err != nil {
		return nil, err
	}

	parser.consumeWhitespaceIfAny()
	err = parser.consumeChar(')')
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func (parser *Parser) consumeActor() (*ActorNode, error) {
	actor := ActorNode{}
	actor.xpos = parser.col

	err := parser.consumeChar('[')
	if err != nil {
		return nil, err
	}
	parser.consumeWhitespaceIfAny()

	actor.Name, err = parser.consumeName()
	if err != nil {
		return nil, err
	}

	parser.consumeWhitespaceIfAny()
	err = parser.consumeChar(']')
	if err != nil {
		return nil, err
	}

	return &actor, nil
}

func (parser *Parser) consumeName() (string, error) {
	// REDOME
	// Good enough for now, but don't re-invent the wheel.
	// Eventually consider existing definitions for
	// acceptable names, such as XML element names,
	// or host names
	exp := `\w+`

	name := regexp.MustCompile(exp).FindString(parser.source[parser.pos:])
	if name == "" {
		return "", parser.errorf("Name does not match expression: %v", exp)
	}

	parser.pos += len(name)
	parser.col += len(name)
	fmt.Printf("Consumed name %q until %v\n", name, parser.pos)
	return name, nil

}

func (parser *Parser) consumeChar(char rune) error {
	runeValue, width := utf8.DecodeRuneInString(parser.source[parser.pos:])
	if runeValue != char {
		return parser.errorf("Expected char: %q", char)
	}

	parser.col += width
	parser.pos += width
	fmt.Printf("Consumed char until %v\n", parser.pos)
	return nil
}

func (parser *Parser) consumeWhitespaceIfAny() {
	exp := `\s*`

	match := regexp.MustCompile(exp).FindString(parser.source[parser.pos:])

	if match == "" {
		return
	}

	parser.pos += len(match)
	parser.col += len(match)
	fmt.Printf("Consumed whitespace until %v\n", parser.pos)
}

func (parser *Parser) consumeRightwardsAccess() error {
	exp := `-+>`

	match := regexp.MustCompile(exp).FindString(parser.source[parser.pos:])

	if match == "" {
		return parser.errorf("rightwards access is not matching regex %q", exp)
	}

	parser.pos += len(match)
	parser.col += len(match)
	fmt.Printf("Consumed rightwards access until %v\n", parser.pos)
	return nil
}

func (parser *Parser) errorf(format string, a ...interface{}) error {
	return fmt.Errorf("%v:%v (pos %v) "+format, parser.row, parser.col, parser.pos, a)
}
