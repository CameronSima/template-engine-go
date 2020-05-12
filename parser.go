package main

import (
	"strings"
)

type Parser struct {
	tokens []Token
	id     string
	//tagFuncs map[string]TagFunc
}

func (p *Parser) Parse(parseUntil []string, start int, end int, context *Context) []Node {

	nodes := make([]Node, 0)
	for i, token := range p.tokens[start:end] {

		switch token.tokenType {
		case TOKEN_VAR:
			node := VariableNode{token}
			nodes = append(nodes, node)

		case TOKEN_TEXT:

			if token.content != " endblock " {
				node := TextNode{token}
				nodes = append(nodes, node)
			}

		case TOKEN_BLOCK:
			bits := strings.Split(token.content, " ")
			command := bits[0]

			if Contains(parseUntil, command) {
				return nodes
			}

			node := GetBlockScopedNode(p, token, command, i, context)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (p *Parser) prepend(token Token) {
	p.tokens = append([]Token{token}, p.tokens...)
}

// func NewParser(tokens []Token) Parser {

// 	return Parser{tokens}

// }

// func NewParser(tokens []Token) Parser {
// 	tagFuncs := map[string]TagFunc{
// 		"block": Block,
// 	}
// 	p := Parser{tokens, tagFuncs}
// 	return p
// }
