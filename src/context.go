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
	functionCalls  []PythonNode
	pythonFuncs    map[string]int
}

func NewContext(source string) Context {
	return Context{
		[]byte(source),
		make(map[string]Node),
		make([]PythonNode, 0),
		make(map[string]int),
	}
}

func (c *Context) AddFunctionCall(n PythonNode) {
	c.functionCalls = append(c.functionCalls, n)
}

func (c *Context) AddRenderContext(key string, node Node) {
	c.render_context[key] = node
}

func (c *Context) AddToContextData(d interface{}, key string) {

	var data []byte
	var err error
	switch d.(type) {
	case ContextData:
		data = d.(ContextData)
	default:
		data, err = json.Marshal(d)
	}

	if err != nil {
		fmt.Println("ERROR MARSHALLING")
	}

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

func (c ContextData) Resolve(variable string) (string, error) {
	keys := strings.Split(variable, ".")
	byteArr, _, _, err := jsonparser.Get(c, keys...)

	if err != nil {
		return "", err
	}
	// TODO: use t (type) to return typed variable
	return string(byteArr), err
}

func (c Context) HasLibrary(libName string) bool {
	//_, hasLib := c.pythonFuncs[libName]
	//return hasLib
	_, err := c.data.Resolve(libName)
	return err == nil
}
