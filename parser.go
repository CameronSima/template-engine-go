package main

import (
	"strings"
)

type Parser struct {
	tokens    []Token
	skipUntil int
	context   *Context
}

func (p *Parser) Parse(parseUntil []string, start int, end int) []Node {
	shouldSkip := false
	nodes := make([]Node, 0)

	for i, token := range p.tokens[start:end] {
		shouldSkip = p.skipUntil != 0 && i < p.skipUntil

		if shouldSkip {
			continue
		}

		switch token.tokenType {
		case TOKEN_VAR:
			node := VariableNode{token}
			nodes = append(nodes, node)

		case TOKEN_TEXT:
			node := TextNode{token}
			nodes = append(nodes, node)

		case TOKEN_BLOCK:
			bits := strings.Split(token.content, " ")
			command := bits[0]

			if Contains(parseUntil, command) {
				p.skipUntil = start + i
				return nodes
			}
			node := p.GetBlockScopedNode(token, command, i)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (p *Parser) GetBlockScopedNode(token Token, command string, currentLine int) Node {
	var node Node

	switch command {
	case "block":
		nodeList := p.Parse([]string{"endblock"}, currentLine+1, len(p.tokens))
		node = NewBlockNode(token, nodeList, p.context)

	case "for":
		nodeList := p.Parse([]string{"endfor"}, currentLine+1, len(p.tokens))
		node = NewForNode(token, nodeList, p.context)

	case "extends":
		node = NewExtendsNode(token, p.context)
	default:
		node = BlankNode{}
	}
	return node
}

func NewParser(source string, context *Context) Parser {
	tokens := NewLexer(source).Tokenize()
	return Parser{tokens, 0, context}
}
