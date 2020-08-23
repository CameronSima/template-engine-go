package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnresolvedNode(t *testing.T) {

	var testSource = `
	<p>the beginning</p>
	{% python_callable.property %}
	<p>the end</p>`

	c := NewContext(`{}`)
	c.pythonFuncs["custom_func"] = 1
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0))
	strippedResult := strings.Replace(RenderNodeList(nodes, c), "\n", "", -1)
	strippedResult = strings.Replace(RenderNodeList(nodes, c), " ", "", -1)
	assert.Equal(t, strippedResult, "\n\t<p>thebeginning</p>\n\t{}\n\t<p>theend</p>")
	//assert.Equal(t, c.pythonFuncs["custom_func"], 1)

	fmt.Println(c.unresolvedVariables)
}
