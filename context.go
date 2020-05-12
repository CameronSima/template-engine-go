package main

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type Context struct {
	data           []byte
	render_context map[string][]Node
}

func (c *Context) AddRenderContext(key string, nodes []Node) {
	c.render_context[key] = nodes
}

func (c Context) GetRenderContext(key string) ([]Node, bool) {
	if nodes, found := c.render_context[key]; found {
		return nodes, found
	}
	return make([]Node, 0, 0), false
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

func NewContext(source string) Context {
	return Context{
		[]byte(source),
		make(map[string][]Node),
	}
}
