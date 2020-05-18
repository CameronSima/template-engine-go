package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockNode(t *testing.T) {
	var testContext = `{"names": [{"name": "johnny", "email": "jbone@gmail.com", "l": [1, 2, 3]}, {"name": "cameron", "email": "boner@gmail.com", "l": [4, 5]}}]}`
	var testSource = `
	{% for n in names %}
		<p>is first? {{ first }}</p>
		<p>{{ n.name }}</p>
		<p>{{ n.email }}</p>

	 {% endfor %}`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0))
	forNode := nodes[1].(ForNode)
	strippedResult := strings.Replace(forNode.Render(c), "\n", "", -1)
	strippedResult = strings.Replace(strippedResult, " ", "", -1)
	assert.Equal(t, strippedResult, "\t\t<p>isfirst?true</p>\t\t<p>johnny</p>\t\t<p>jbone@gmail.com</p>\t\t\t<p>isfirst?false</p>\t\t<p>cameron</p>\t\t<p>boner@gmail.com</p>\t")
}

func TestNestedBlockNode(t *testing.T) {
	var testContext = `{"names": [{"name": "johnny", "email": "jbone@gmail.com", "l": [1, 2, 3]}, {"name": "cameron", "email": "boner@gmail.com", "l": [4, 5]}}]}`
	var testSource = `
	{% for n in names %}
		<p>is first? {{ first }}</p>
		<p>{{ n.name }}</p>
		<p>{{ n.email }}</p>

		{% for num in n.l %}
			<p>{{ num }}</p>
		{% endfor %}

	 {% endfor %}`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0))
	fNode := nodes[1].(ForNode)
	nestedFNode := fNode.nodes[7].(ForNode)
	result := RenderNodeList(nodes, c)

	assert.Equal(t, nestedFNode.token.content, "for num in n.l")
	assert.Equal(t, nestedFNode.loopVariable, "num")
	assert.Equal(t, nestedFNode.loopArrayName, "n.l")
	assert.Equal(t, strings.Contains(result, "<p>1</p>"), true)
	assert.Equal(t, strings.Contains(result, "<p>2</p>"), true)
	assert.Equal(t, strings.Contains(result, "<p>3</p>"), true)
	assert.Equal(t, strings.Contains(result, "<p>4</p>"), true)
	assert.Equal(t, strings.Contains(result, "<p>5</p>"), true)
}

func TestUrlNode(t *testing.T) {
	var testContext = `{"urls": [{"name": "home", "pattern": "/home"}], "http_host": "localhost:8000"}`
	var testSource = `
	<div>
		<a href={% url "home" %} />
	 </div>`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0))

	fmt.Println(RenderNodeList(nodes, c))

}
