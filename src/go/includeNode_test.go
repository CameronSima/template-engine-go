package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeNode(t *testing.T) {

	// test twice so we know the caching works
	var testSource = `
	{% include '../test_templates/test_include.html' %}
	{% include '../test_templates/test_include.html' %}
	<p>the end</p>`

	c := NewContext("{}")
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0))
	strippedResult := strings.Replace(RenderNodeList(nodes, &c), "\n", "", -1)
	strippedResult = strings.Replace(RenderNodeList(nodes, &c), " ", "", -1)
	assert.Equal(t, strippedResult, "\n\t<p>hi</p>\n\t<p>hi</p>\n\t<p>theend</p>")
}
