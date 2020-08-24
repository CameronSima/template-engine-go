package main

import (
	"fmt"
	"strings"
)

// Custom tags: {% my_python_function param1 param2 %}
type PythonNode struct {
	token        Token    `json:"-"`
	FunctionName string   `json:"functionName"`
	Parameters   []string `json:"parameters"`
}

func NewPythonNode(token Token, context *Context) Node {
	bits := strings.Split(token.Content, " ")
	functionName := bits[0]
	params := bits[1:]
	resolvedParams := make([]string, 0)

	// resolve params from context if they're variables and not raw strings: "value"
	var resolved string
	for _, p := range params {
		if strings.Contains(p, `"`) {
			resolved = strings.Trim(p, `"`)
		} else {
			resolved, _ = context.data.Resolve(p)
		}
		resolvedParams = append(resolvedParams, resolved)
	}
	n := PythonNode{token, functionName, resolvedParams}

	//if context.HasLibrary(functionName) {
	context.AddFunctionCall(n)
	//}
	return n
}

func (n PythonNode) Render(context *Context) string {
	return "{}"
}

func (n PythonNode) String() string {
	return fmt.Sprintf("Type: PythonNode \nFunction name: %v \nParameters: %v", n.FunctionName, n.Parameters)
}
