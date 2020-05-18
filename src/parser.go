package main

import (
	"strings"
)

type Parser struct {
	tokens    TokenStack
	skipUntil int
	context   *Context
}

func (p *Parser) Parse(parseUntil []string) []Node {
	nodes := make([]Node, 0)

	for p.tokens.IsEmpty == false {
		token, _ := p.tokens.NextToken()

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
				p.tokens.PrependToken(token)
				return nodes
			}
			node := p.GetBlockScopedNode(token, command)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (p *Parser) GetBlockScopedNode(token Token, command string) Node {
	var node Node

	switch command {
	case "block":
		nodeList := p.Parse([]string{"endblock"})
		node = NewBlockNode(token, nodeList, p.context)
	case "for":
		nodeList := p.Parse([]string{"endfor"})
		node = NewForNode(token, nodeList, p.context)
	case "extends":
		node = NewExtendsNode(token, p.context)
	case "url":
		node = UrlNode{token}
	case "static":
		node = StaticNode{token}
	case "csrftoken":
		node = CsrfNode{}
	default:
		node = BlankNode{}
	}
	return node
}

func NewParser(source string, context *Context) Parser {
	tokens := NewLexer(source).Tokenize()
	return Parser{NewTokenStack(tokens), 0, context}
}
