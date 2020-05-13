package main

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type Context struct {
	data           []byte
	render_context map[string]Node
}

func (c *Context) AddRenderContext(key string, node Node) {
	c.render_context[key] = node
}

func (c Context) GetRenderContext(key string) (Node, bool) {
	if node, found := c.render_context[key]; found {
		return node, found
	}
	return BlankNode{}, false
}

func (c Context) Resolve(variable string) string {
	keys := strings.Split(variable, ".")
	byteArr, _, _, err := jsonparser.Get(c.data, keys...)

	if err != nil {
		fmt.Println("Error resolving variable: " + variable)
		fmt.Println(err)
	}
	// TODO: use t (type) to return typed variable
	return string(byteArr)
}

// func (c Context) ResolveArray(variable string) []string {
// 	values := make([]string, 0)
// 	keys := strings.Split(variable, ".")

// 	jsonparser.ArrayEach(c.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 		fmt.Println(jsonparser.Get(value)
// 	}, keys...)

// }

func NewContext(source string) Context {
	return Context{
		[]byte(source),
		make(map[string]Node),
	}
}
