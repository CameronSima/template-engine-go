package main

import (
	"fmt"
	"strings"
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
	value := context.Resolve(variable)
	return value
}

type TextNode struct {
	token Token
}

func (n TextNode) Render(context Context) string {
	return n.token.content
}

type BlockNode struct {
	token Token
	nodes []Node
}

func (n BlockNode) Render(context Context) string {
	var result string
	bits := strings.Split(n.token.content, " ")
	parameter := bits[1]

	fmt.Println("BLOCK NODE PARAMETER")
	fmt.Println(parameter)

	if nodeList, exists := context.GetRenderContext(parameter); exists {
		return RenderNodeList(nodeList, context)
	} else {
		context.AddRenderContext(parameter, n.nodes)
		return "testng"
	}
	return result
}

type ExtendsNode struct {
	token Token
}

func (n ExtendsNode) Render(context Context) string {
	bits := strings.Split(n.token.content, " ")
	parameter := bits[1]
	templateSource := ReadTemplate(parameter)

	lexer := NewLexer(templateSource)
	parser := Parser{lexer.Tokenize(), "extends parser"}
	template := Template{parser, templateSource}
	//template := NewTemplate(templateSource)
	r := template.Render(context)

	fmt.Println("EXTENBDS NODE")
	fmt.Println(r)
	return r
}

func GetBlockScopedNode(p *Parser, token Token, command string, currentLine int) Node {
	var node Node

	switch command {
	case "block":
		nodeList := p.Parse([]string{"endblock"}, currentLine+1, len(p.tokens))
		node = BlockNode{token, nodeList}

	case "extends":
		node = ExtendsNode{token}
	default:
		node = BlankNode{}

	}
	return node

}

// func RenderBlock(context Context, n Node, command string, variable string) string {
// 	if node, exists := context.GetRenderContext(variable); exists {
// 		fmt.Println("PULLING CONTEXT")

// 		fmt.Println(node.Render(context))
// 		return node.Render(context)
// 	} else {
// 		fmt.Println("ADDING TO CONTEXT")

// 		context.AddRenderContext(variable, n)
// 		return ""
// 	}
// }

// func RenderExtends(context Context, variable string) string {
// 	templateSource := ReadTemplate(variable)

// 	fmt.Println("EXTENDS NODE RENDER")
// 	fmt.Println(templateSource)

// 	lexer := NewLexer(templateSource)
// 	parser := Parser{lexer.Tokenize(), "extends parser"}
// 	template := Template{parser, templateSource}
// 	//template := NewTemplate(templateSource)
// 	return template.Render(context)
// }
