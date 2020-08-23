package main

import (
	"encoding/json"
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
	value string
}

func NewVariableNode(token Token, context *Context) VariableNode {
	return VariableNode{token, ""}
}

func (n VariableNode) Render(context Context) string {
	//return n.value
	variable := strings.Replace(n.token.content, " ", "", -1)
	if result, err := context.data.Resolve(variable); err == nil {
		return result
	} else {
		return ""
	}
}

func (n VariableNode) String() string {
	return fmt.Sprintf("Type: VariableNode, Token: %v", n.token)
}

type UnresolvedVariableNode struct {
	variableNode Node
}

func (n UnresolvedVariableNode) Render(context Context) string {
	context.AddUnresolvedVariable(n.variableNode.(VariableNode))
	return "{}"
}

func (n UnresolvedVariableNode) String() string {
	return fmt.Sprintf("Type: UnresolvedVariableNode, VariableNode: %v", n.variableNode)
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
	isKeyValuePair bool
	token          Token
	nodes          []Node
	loopVariable   string
	loopArrayName  string
	loopContext    ForLoopContext
	parent         Node
}

func NewForNode(token Token, parsedNodes []Node, context *Context) ForNode {
	var loopVariable string
	var loopArrayName string
	isKeyValuePair := false
	bits := strings.Split(token.content, " ")

	if len(bits) == 5 {
		isKeyValuePair = true
		loopVariable = bits[1] + bits[2]
		loopArrayName = bits[4]
	} else {
		loopVariable = bits[1]
		loopArrayName = bits[3]
	}

	return ForNode{
		isKeyValuePair,
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

			fmt.Println(childNode)

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
	lookupKeys := strings.Split(n.loopArrayName, ".")
	results := make([]ForLoopVariable, 0)

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		variable := ForLoopVariable{string(dataType), value, offset}
		results = append(results, variable)
	}, lookupKeys...)
	return results
}

type ForLoopVariableNode struct {
	node     VariableNode
	variable ForLoopVariable
	forLoop  ForNode
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

		} else {
			// not a for loop node, render variable node normally
			result = n.node.Render(context)
		}
	} else {
		// object lookup in variable context
		if r, err := context.data.Resolve(n.node.token.content); err == nil {
			result = r
		} else {
			unresolvedNode := UnresolvedVariableNode{n.node}
			result = unresolvedNode.Render(context)
		}
	}

	return result
}

func (n ForLoopVariableNode) String() string {
	return fmt.Sprintf("Type: ForLoopVariableNode, Variable: %v, InnerNode: %v", n.variable, n.node)

}

type ForLoopVariable struct {
	dataType string
	value    ContextData
	index    int
}

type ForLoopContext struct {
	Counter      int         `json:"counter"`
	Counter0     int         `json:"counter0"`
	Revcounter   int         `json:"revcounter"`
	Revcounter0  int         `json:"revcounter0"`
	First        bool        `json:"first"`
	Last         bool        `json:"last"`
	LoopVariable string      `json:"loopVariable"`
	Data         ContextData `json:"data"`
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

type UrlNode struct {
	token Token
}

func (n UrlNode) Render(context Context) string {
	var pattern string
	bits := strings.Split(n.token.content, " ")
	viewName := strings.Trim(bits[1], `"'`)
	http_host, _ := context.data.Resolve("http_host")

	jsonparser.ArrayEach(context.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var urlDef UrlDefinition
		err = json.Unmarshal(value, &urlDef)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(urlDef)

		if urlDef.Name == viewName {
			pattern = urlDef.Pattern
		}
	}, "urls")
	return fmt.Sprintf(`"%s%s"`, http_host, pattern)
}

func (n UrlNode) String() string {
	return n.token.content
}

type UrlDefinition struct {
	Name    string
	Pattern string
}

type StaticNode struct {
	token Token
}

func (n StaticNode) Render(context Context) string {
	bits := strings.Split(n.token.content, " ")
	url := strings.Trim(bits[1], `"'`)
	http_host, _ := context.data.Resolve("http_host")
	static_url, _ := context.data.Resolve("static_url")
	return fmt.Sprintf(`"%s%s%s"`, http_host, static_url, url)
}

func (n StaticNode) String() string {
	return "Type: StaticNode \n" + n.token.content
}

type CsrfNode struct{}

func (n CsrfNode) Render(context Context) string {
	csrfToken, _ := context.data.Resolve("cookies.CSRF_TOKEN")
	return fmt.Sprintf(`<input type="hidden" name="csrfmiddlewaretoken" value="%s">`, csrfToken)
}

func (n CsrfNode) String() string {
	return "Type: CsrfNode \n"
}
