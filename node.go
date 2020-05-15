package main

import (
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
	result, _ := context.data.Resolve(variable)
	return result
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

func (n ForNode) GetForLoopData(context Context) []ForLoopVariable {
	keys := strings.Split(n.loopArrayName, ".")
	values := make([]ForLoopVariable, 0)
	jsonparser.ArrayEach(context.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		variable := ForLoopVariable{string(dataType), value}
		values = append(values, variable)

	}, keys...)
	return values
}

func (n ForNode) Render(context Context) string {
	var rendered strings.Builder
	variables := n.GetForLoopData(context)
	l := len(variables)

	for i, variable := range variables {
		loopContext := NewForLoopContext(i, i-l, i == 0, i == l, n, n.loopVariable)
		var node Node

		for _, childNode := range n.nodes {

			switch childNode.(type) {
			case VariableNode:
				node = ForLoopVariableNode{
					childNode.(VariableNode),
					variable,
					loopContext,
				}
			default:
				node = childNode
			}
			rendered.WriteString(node.Render(context))
		}
	}

	return rendered.String()
}

type ForLoopVariableNode struct {
	node        VariableNode
	variable    ForLoopVariable
	loopContext ForLoopContext
}

func (n ForLoopVariableNode) Render(context Context) string {
	keys := strings.Split(n.node.token.content, ".")
	firstKey := keys[0]
	context.AddToContextData(n.loopContext, "forloop")

	var result string
	lookupVariable := strings.Join(keys[1:len(keys)], ".")

	// basic key look up, ('name') or forloop, vs object lookup ('person.name')
	if lookupVariable == "" {
		// try from variable context
		if r, err := n.variable.value.Resolve(firstKey); err == nil {
			result = r

			// try from forloop context
		} else if r, err := context.data.Resolve("forloop." + firstKey); err == nil {
			result = r

		} else {
			// not a for loop node, render variable node normally
			result = n.node.Render(context)
		}
	} else {
		// object lookup in variable context
		r, _ := n.variable.value.Resolve(lookupVariable)
		result = r
	}

	return result
}

type ForLoopVariable struct {
	dataType string
	value    ContextData
}

type ForLoopContext struct {
	Counter      int     `json:"counter"`
	Counter0     int     `json:"counter0"`
	Revcounter   int     `json:"revcounter"`
	Revcounter0  int     `json:"revcounter0"`
	First        bool    `json:"first"`
	Last         bool    `json:"last"`
	Parent       ForNode `json:"parent"`
	LoopVariable string  `json:"loopVariable"`
}

func NewForLoopContext(counter0 int, revcounter0 int, first bool, last bool, parent ForNode, loopVariable string) ForLoopContext {
	return ForLoopContext{
		counter0 - 1,
		counter0,
		revcounter0 - 1,
		revcounter0,
		first,
		last,
		parent,
		loopVariable,
	}
}
