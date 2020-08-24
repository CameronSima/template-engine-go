package main

import (
	"fmt"
	"strings"
)

type IncludeNode struct {
	token Token
	nodes []Node
}

func NewIncludeNode(token Token, context *Context) IncludeNode {
	bits := strings.Split(token.Content, " ")
	parameter := bits[1]

	if includeNode, exists := context.GetRenderContext(parameter); exists {
		return includeNode.(IncludeNode)
	} else {
		templateSource := ReadTemplate(parameter)
		parser := NewParser(templateSource, context)
		nodes := parser.Parse(make([]string, 0))
		includeNode := IncludeNode{token, nodes}

		// cache to prevent re-parsing if this is used in a for-loop
		context.AddRenderContext(parameter, includeNode)
		return includeNode
	}
}

func (n IncludeNode) Render(context *Context) string {
	return RenderNodeList(n.nodes, context)
}

func (n IncludeNode) String() string {
	return fmt.Sprintf("Type: IncludeNode, Token: %v, Children: %v", n.token, n.nodes)
}
