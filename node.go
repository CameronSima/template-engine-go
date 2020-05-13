package main

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type Node interface {
	Render(context Context) string
}

type BlankNode struct{}

func (n BlankNode) Render(context Context) string { return "" }

type VariableNode struct {
	token Token
}

func (n VariableNode) Render(context Context) string {
	variable := strings.Replace(n.token.content, " ", "", -1)
	return context.Resolve(variable)
}

type TextNode struct {
	token Token
}

func (n TextNode) Render(context Context) string {
	stripped := strings.Replace(n.token.content, " ", "", -1)
	if stripped == "" {
		return stripped
	}
	return n.token.content

}

// ------------- Block-typed nodes ---------------

// base block-node
type BlockNode struct {
	token       Token
	nodes       []Node
	placeholder bool
}

func NewBlockNode(token Token, parsedNodes []Node, context *Context) BlockNode {
	bits := strings.Split(token.content, " ")
	parameter := bits[1]
	var newNode BlockNode

	if _, exists := context.GetRenderContext(parameter); exists {
		newNode = BlockNode{token, parsedNodes, false}
	} else {
		newNode = BlockNode{token, parsedNodes, true}

	}
	context.AddRenderContext(parameter, newNode)
	return newNode
}

func (n BlockNode) Render(context Context) string {
	var result string
	bits := strings.Split(n.token.content, " ")
	parameter := bits[1]

	if n.placeholder == true {
		node, _ := context.GetRenderContext(parameter)
		result = RenderNodeList(node.(BlockNode).nodes, context)

	} else {
		result = ""
	}
	return result
}

type ExtendsNode struct {
	token Token
	nodes []Node
}

func (n ExtendsNode) Render(context Context) string {
	return RenderNodeList(n.nodes, context)
}

func NewExtendsNode(token Token, context *Context) ExtendsNode {
	bits := strings.Split(token.content, " ")
	parameter := bits[1]
	templateSource := ReadTemplate(parameter)
	parser := NewParser(templateSource, context)
	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))
	return ExtendsNode{token, nodes}
}

type ForNode struct {
	token         Token
	nodes         []Node
	loopVariable  string
	loopArrayName string
}

func NewForNode(token Token, parsedNodes []Node, context *Context) ForNode {
	bits := strings.Split(token.content, " ")
	loopVariable := bits[1]
	loopArrayName := bits[3]
	return ForNode{
		token,
		parsedNodes,
		loopVariable,
		loopArrayName,
	}
}

func (n ForNode) Render(context Context) string {
	renderedNodes := make([]Node, 0)
	keys := strings.Split(n.loopArrayName, ".")

	for _, node := n.nodes {
		switch type(node) {
		case VariableNode:

		
		}
	}

	index := 0
	jsonparser.ArrayEach(context.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fmt.Println(string(value))

		node := n.nodes[0].(VariableNode)
		node.Render()


	}, keys...)

	return ""
}
