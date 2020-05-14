package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type ContextData []byte

type Context struct {
	data           ContextData
	render_context map[string]Node
}

func (c *Context) AddRenderContext(key string, node Node) {
	c.render_context[key] = node
}

func (c *Context) AddToContextData(d interface{}, key string) {
	data, _ := json.Marshal(d)

	value, err := jsonparser.Set(c.data, data, key)
	if err != nil {
		fmt.Println("Error adding to context")
	}
	c.data = value
}

func (c Context) GetRenderContext(key string) (Node, bool) {
	if node, found := c.render_context[key]; found {
		return node, found
	}
	return BlankNode{}, false
}

func (c ContextData) Resolve(variable string) string {
	keys := strings.Split(variable, ".")
	byteArr, _, _, err := jsonparser.Get(c, keys...)

	if err != nil {
		fmt.Println("Error resolving variable: " + variable)
		fmt.Println("Context data: ")
		fmt.Println(string(c))
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
