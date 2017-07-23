package dreitafel

import (
	"fmt"
	"regexp"
	"sync"
)

// REDOME
// Good enough for now, but don't re-invent the wheel.
// Eventually consider existing definitions for
// acceptable names, such as XML element names,
// or host names
var titleExp = `[\p{L}\d_]+`

type Lexer struct {
	line   *string
	lineNo int

	col int
	row int

	recognized   chan<- Token
	unrecognized chan<- error
}

func KeepParsing(lines <-chan *string, diagrams chan<- *FmcBlockDiagram, errors chan<- error, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	tokens := make(chan Token)

	go tokenize(lines, tokens, errors)

	buildDiagram(tokens, diagrams, errors)
}

func tokenize(lines <-chan *string, recognizedTokens chan<- Token, errors chan<- error) {
	var line *string

	lineNo := 0
	for line = <-lines; line != nil; line = <-lines {
		lineNo++
		lexer := Lexer{row: lineNo, recognized: recognizedTokens, unrecognized: errors, line: line, lineNo: lineNo}
		lexer.tokenizeLine()
	}
	recognizedTokens <- nil // terminate the consumer as well
	fmt.Println("Tokenizer: done")
}

func (lexer *Lexer) tokenizeLine() {
	defer func() {
		lexer.recognized <- &BaseToken{from: 0, to: 0, tokenType: TokenTypeLineEnd}
	}()

	for {
		var token Token
		// keep recognizing actor, storage, connection, whitespace & EOL
		// on error, skip rest of line

		fmt.Printf("%v <= %v\n", lexer.col, len(*lexer.line))
		if lexer.isdone() {
			return
		}

		if lexer.acceptWhitespaceIfAny() {
			continue
		}

		token = lexer.acceptActor()
		if token == nil {
			token = lexer.acceptStorage()
		}
		if token == nil {
			token = lexer.acceptRightwardsAccess()
		}
		if token == nil {
			lexer.unrecognized <- fmt.Errorf("Don't know what that is: %v", (*lexer.line)[lexer.col:])
			lexer.col = len(*lexer.line)
		} else {
			fmt.Printf("Recognized until %v/%v: %q\n", lexer.col, len(*lexer.line), token)
			lexer.recognized <- token
		}
	}
}

func buildDiagram(tokens <-chan Token, diagrams chan<- *FmcBlockDiagram, errors chan<- error) {
	diagram := FmcBlockDiagram{title: "My first diagram"}

	for token := <-tokens; token != nil; token = <-tokens {

		fmt.Println(token)
		if token.GetTokenType() == TokenTypeActor {
			node := token.(*Node)
			actor := Actor{}
			actor.FmcBaseNode.id = node.title
			actor.FmcBaseNode.title = node.title
			diagram.nodes = append(diagram.nodes, &actor)
		} else if token.GetTokenType() == TokenTypeStorage {
			node := token.(*Node)
			storage := Storage{}
			storage.FmcBaseNode.id = node.title
			storage.FmcBaseNode.title = node.title
			diagram.nodes = append(diagram.nodes, &storage)
		}
		// I want pattern matching. :(

		// • actor/storage
		//   -> emit actor
		//   -> if remembering connection:
		//      (handle non-bipartite graph)
		//      -> emit connection from previous node to this one
		//   -> memorize node & reset connection
		// • connection
		//   -> if not remembering actor:
		//     -> emit syntax error
		//     -> skip until next element
		//   -> memorize connection
		//   -> go to beginning of loop
		// •
		// :error -> skip until next actor
	}

	diagrams <- &diagram
}

func (lexer *Lexer) isdone() bool {
	return lexer.col >= len(*lexer.line)
}

func (lexer *Lexer) acceptStorage() Token {
	exp := `\(\s*(` + titleExp + `)\s*\)`

	matches := regexp.MustCompile(exp).FindSubmatch([]byte((*lexer.line)[lexer.col:]))

	if matches == nil {
		return nil
	}

	title := string(matches[1])

	node := Node{title: title}
	node.tokenType = TokenTypeStorage
	node.from = lexer.col + len(title)
	node.to = lexer.col
	lexer.col += len(string(matches[0]))
	return &node
}

func (lexer *Lexer) acceptActor() Token {
	exp := `\[\s*(` + titleExp + `)\s*\]`

	matches := regexp.MustCompile(exp).FindSubmatch([]byte((*lexer.line)[lexer.col:]))

	if matches == nil {
		return nil
	}

	title := string(matches[1])

	node := Node{title: title}
	node.tokenType = TokenTypeActor
	node.from = lexer.col + len(title)
	node.to = lexer.col
	lexer.col += len(string(matches[0]))
	return &node
}

func (lexer *Lexer) acceptWhitespaceIfAny() bool {
	exp := `\s*`

	match := regexp.MustCompile(exp).FindString((*lexer.line)[lexer.col:])

	if match == "" {
		return false
	}

	lexer.col += len(match)
	fmt.Printf("Ate whitespace until %v\n", lexer.col)
	return true
}

func (lexer *Lexer) acceptRightwardsAccess() Token {
	exp := `-+>`

	match := regexp.MustCompile(exp).FindString((*lexer.line)[lexer.col:])

	if match == "" {
		return nil
	}

	lexer.col += len(match)
	fmt.Printf("Ate rightwards access until %v\n", lexer.col)
	return &BaseToken{tokenType: TokenTypeRightAccess, from: lexer.col - len(match), to: lexer.col}
}

func (lexer *Lexer) errorf(format string, a ...interface{}) error {
	return fmt.Errorf("%v:%v "+format, lexer.row, lexer.col, a)
}
