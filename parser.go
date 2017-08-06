package dreitafel

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	log "github.com/Sirupsen/logrus"
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

func KeepParsing(lines <-chan *string, diagrams chan<- *FmcBlockDiagram, errorsChan chan<- error) {

	tokens := make(chan Token)

	go tokenize(lines, tokens, errorsChan)

	buildDiagram(tokens, diagrams, errorsChan)
}

func tokenize(lines <-chan *string, recognizedTokens chan<- Token, errorsChan chan<- error) {
	defer close(recognizedTokens) // terminate the consumer as well

	lineNo := 0
	for line := range lines {
		log.Debugf("Parsing line: '%v' @%v", *line, line)
		lineNo++
		lexer := Lexer{row: lineNo, recognized: recognizedTokens, unrecognized: errorsChan, line: line, lineNo: lineNo}
		lexer.tokenizeLine()
	}
	log.Debug("Tokenizer: done")
}

func (lexer *Lexer) tokenizeLine() {
	defer func() {
		lexer.recognized <- &BaseToken{from: 0, to: 0, tokenType: TokenTypeLineEnd}
	}()

	for {
		var token Token
		// keep recognizing actor, storage, connection, whitespace & EOL
		// on error, skip rest of line

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
			token = lexer.acceptLeftwardsAccess()
		}
		if token == nil {
			lexer.unrecognized <- lexer.errorf("Don't know what that is: %v", (*lexer.line)[lexer.col:])
			lexer.col = len(*lexer.line)
		} else {
			log.Debugf("Recognized until %v/%v: %v (%q)", lexer.col, len(*lexer.line), reflect.TypeOf(token), token)
			lexer.recognized <- token
		}
	}
}

func buildDiagram(tokens <-chan Token, diagrams chan<- *FmcBlockDiagram, errorsChan chan<- error) {
	defer close(diagrams)

	plog := log.New()
	diagram := FmcBlockDiagram{title: "My first diagram"}

	var prevNode FmcNode
	var prevToken Token

	resetPrev := func() { prevNode = nil; prevToken = nil }
	setPrev := func(n FmcNode, t Token) { prevNode = n; prevToken = t }

	for token := range tokens {

		plog.Debug(token)
		switch token.GetTokenType() {
		case TokenTypeActor:
			node := token.(*Node)
			actor := Actor{}
			actor.FmcBaseNode.id = node.title
			actor.FmcBaseNode.title = node.title
			diagram.nodes = append(diagram.nodes, &actor)

			if prevToken != nil {
				switch prevToken.GetTokenType() {
				case TokenTypeLeftAccess, TokenTypeRightAccess:
					edge := prevNode.(*FmcBaseEdge)
					if edge.actor != nil {
						errorsChan <- errors.New("Syntax error: Bipartite graph must connect Actor to Storage, not Actor to Actor directly.")
						resetPrev()
						continue
					}
					edge.actor = &actor
					diagram.edges = append(diagram.edges, edge)

					// other tokens before are fine:
					// a “multiple expressions” line
				}
			}

			setPrev(&actor, token)
		case TokenTypeStorage:
			node := token.(*Node)
			storage := Storage{}
			storage.FmcBaseNode.id = node.title
			storage.FmcBaseNode.title = node.title
			diagram.nodes = append(diagram.nodes, &storage)

			if prevToken != nil {
				switch prevToken.GetTokenType() {
				case TokenTypeLeftAccess, TokenTypeRightAccess:
					edge := prevNode.(*FmcBaseEdge)
					if edge.storage != nil {
						errorsChan <- errors.New("Syntax error: Bipartite graph must connect Actor to Storage, not Storage to Storage directly.")
						resetPrev()
						continue
					}
					edge.storage = &storage
					diagram.edges = append(diagram.edges, edge)

					// other tokens before are fine:
					// a “multiple expressions” line
				}
			}

			setPrev(&storage, token)
		case TokenTypeLeftAccess, TokenTypeRightAccess:
			if prevToken == nil {
				errorsChan <- errors.New("Syntax error: Dangling access")
				resetPrev()
				continue
			}

			edge := FmcBaseEdge{}
			switch prevToken.GetTokenType() {
			case TokenTypeActor:
				edge.actor = prevNode.(*Actor)
				if token.GetTokenType() == TokenTypeLeftAccess {
					edge.edgeType = EdgeTypeRead
					log.Debug("Read access!")
				} else {
					log.Debug("Write access!")
					edge.edgeType = EdgeTypeWrite
				}
				setPrev(&edge, token)
			case TokenTypeStorage:
				edge.storage = prevNode.(*Storage)
				if token.GetTokenType() == TokenTypeLeftAccess {
					log.Debug("Write access!")
					edge.edgeType = EdgeTypeWrite
				} else {
					log.Debug("Read access!")
					edge.edgeType = EdgeTypeRead
				}
				setPrev(&edge, token)
			default:
				errorsChan <- errors.New("Syntax error: Read Access must connect Actor and Storage")
				resetPrev()
			}
		}
	}

	diagrams <- &diagram
	log.Debug("Diagram builder done.")
}

func (lexer *Lexer) isdone() bool {
	return lexer.col >= len(*lexer.line)
}

func (lexer *Lexer) acceptStorage() Token {
	exp := `^\(\s*(` + titleExp + `)\s*\)`

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
	exp := `^\[\s*(` + titleExp + `)\s*\]`

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
	exp := `^\s*`

	match := regexp.MustCompile(exp).FindString((*lexer.line)[lexer.col:])

	if match == "" {
		return false
	}

	lexer.col += len(match)
	log.Debugf("Ate whitespace until %v", lexer.col)
	return true
}

func (lexer *Lexer) acceptRightwardsAccess() Token {
	exp := `^-+>`

	match := regexp.MustCompile(exp).FindString((*lexer.line)[lexer.col:])

	if match == "" {
		return nil
	}

	lexer.col += len(match)
	log.Debugf("Ate rightwards access until %v", lexer.col)
	return &BaseToken{tokenType: TokenTypeRightAccess, from: lexer.col - len(match), to: lexer.col}
}

func (lexer *Lexer) acceptLeftwardsAccess() Token {
	exp := `^<-+`

	match := regexp.MustCompile(exp).FindString((*lexer.line)[lexer.col:])

	if match == "" {
		return nil
	}

	lexer.col += len(match)
	log.Debugf("Ate leftwards access until %v", lexer.col)
	return &BaseToken{tokenType: TokenTypeLeftAccess, from: lexer.col - len(match), to: lexer.col}
}

func (lexer *Lexer) errorf(format string, a ...interface{}) error {
	return fmt.Errorf("%v:%v "+format, lexer.row, lexer.col, a)
}
