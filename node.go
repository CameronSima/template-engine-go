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
	return context.data.Resolve(variable)
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

		var node Node
		for _, childNode := range n.nodes {
			switch childNode.(type) {
			case VariableNode:
				loopContext := NewForLoopContext(i, i-l, i == 0, i == l, n, n.loopVariable)
				context.AddToContextData(loopContext, "forloop")
				node = ForLoopVariableNode{
					childNode.(VariableNode),
					variable,
					n.loopVariable,
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
	node         VariableNode
	variable     ForLoopVariable
	loopVariable string
}

func (n ForLoopVariableNode) Render(context Context) string {
	keys := strings.Split(n.node.token.content, ".")
	v := keys[0]
	lookupVariable := strings.Join(keys[1:len(keys)], ".")

	if v == n.loopVariable {
		if lookupVariable == "" {
			return string(n.variable.value)
		}
		return n.variable.value.Resolve(lookupVariable)
	} else {
		return n.node.Render(context)
	}
}

type ForLoopVariable struct {
	dataType string
	value    ContextData
}

type ForLoopContext struct {
	counter      int
	counter0     int
	revcoounter  int
	revcounter0  int
	first        bool
	last         bool
	parent       ForNode
	loopVariable string
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
