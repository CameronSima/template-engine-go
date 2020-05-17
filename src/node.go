package main

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type Node interface {
	Render(context Context) string
	String() string
}

type BlankNode struct{}

func (n BlankNode) Render(context Context) string { return "" }
func (n BlankNode) String() string                { return "Type: BlankNode" }

type VariableNode struct {
	token Token
}

func (n VariableNode) Render(context Context) string {
	variable := strings.Replace(n.token.content, " ", "", -1)
	result, _ := context.data.Resolve(variable)
	return result
}

func (n VariableNode) String() string {
	return fmt.Sprintf("Type: VariableNode, Token: %v", n.token)
}

type TextNode struct {
	token Token
}

func (n TextNode) Render(context Context) string {
	stripped := strings.Replace(n.token.content, " ", "", -1)
	if stripped == "" || stripped == "\n" {
		return stripped
	}
	return n.token.content
}

func (n TextNode) String() string {
	return fmt.Sprintf("Type: TextNode, Token: %v", n.token)
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

func (n BlockNode) String() string {
	return fmt.Sprintf("Type: BlockNode, Token: %v, Children: %v", n.token, n.nodes)
}

type ExtendsNode struct {
	token Token
	nodes []Node
}

func NewExtendsNode(token Token, context *Context) ExtendsNode {
	bits := strings.Split(token.content, " ")
	parameter := bits[1]
	templateSource := ReadTemplate(parameter)
	parser := NewParser(templateSource, context)
	nodes := parser.Parse(make([]string, 0))
	return ExtendsNode{token, nodes}
}

func (n ExtendsNode) Render(context Context) string {
	return RenderNodeList(n.nodes, context)
}

func (n ExtendsNode) String() string {
	return fmt.Sprintf("Type: ExtendsNode, Token: %v, Children: %v", n.token, n.nodes)
}

type ForNode struct {
	token         Token
	nodes         []Node
	loopVariable  string
	loopArrayName string
	loopContext ForLoopContext
	parent        Node
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
		NewForLoopContext(-1, -1, false, false, "", make([]byte, 0)),
		BlankNode{},
	}
}

func (n ForNode) Render(context Context) string {
	var rendered strings.Builder
	variables := n.GetForLoopData(context.data)	
	l := len(variables)

	for i, variable := range variables {
		n.loopContext = NewForLoopContext(i, i-l, i == 0, i == l, n.loopVariable, variable.value)
		context.AddToContextData(variable.value, n.loopVariable)

		var node Node
		for _, childNode := range n.nodes {

			switch childNode.(type) {
			// wrap variable node with loop context
			case VariableNode:
				node = ForLoopVariableNode{
					childNode.(VariableNode),
					variable,
					n,
				}

			default:
				node = childNode
			}
			rendered.WriteString(node.Render(context))
		}
	}
	return rendered.String()
}

func (n ForNode) String() string {
	return fmt.Sprintf("Type: ForNode\n Token: %v\n Parent: %v\n Children: %v\n", n.token, n.parent, n.nodes)
}

func (n ForNode) GetForLoopData(data ContextData) []ForLoopVariable {
	keys := strings.Split(n.loopArrayName, ".")
	values := make([]ForLoopVariable, 0)

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		variable := ForLoopVariable{string(dataType), value}
		values = append(values, variable)
	}, keys...)

	return values
}

type ForLoopVariableNode struct {
	node        VariableNode
	variable    ForLoopVariable
	forLoop ForNode
}

func (n ForLoopVariableNode) Render(context Context) string {
	keys := strings.Split(n.node.token.content, ".")
	firstKey := keys[0]
	context.AddToContextData(n.forLoop.loopContext, "forloop")
	
	var result string
	lookupVariable := strings.Join(keys[1:len(keys)], ".")

	// basic key look up, ('name') or forloop, vs object lookup ('person.name')
	if lookupVariable == "" {

		// try from context
		if r, err := context.data.Resolve(firstKey); err == nil {
			result = r
		// try from forloop context
		} else if r, err := context.data.Resolve("forloop." + firstKey); err == nil {
			result = r

		} else if r, err := context.data.Resolve(firstKey); err == nil {
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

func (n ForLoopVariableNode) String() string {
	return fmt.Sprintf("Type: ForLoopVariableNode, Variable: %v, InnerNode: %v", n.variable, n.node)

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
	LoopVariable string  `json:"loopVariable"`
	Data ContextData `json:"data"`
	//CurrentVariable ForLoopVariable `json:"currentVariable"`
}

func NewForLoopContext(counter0 int, revcounter0 int, first bool, last bool, loopVariable string, data []byte) ForLoopContext {
	return ForLoopContext{
		counter0 - 1,
		counter0,
		revcounter0 - 1,
		revcounter0,
		first,
		last,
		loopVariable,
		data,
	}
}

// type UrlNode struct {
// 	token Token
// 	viewName string

// }
