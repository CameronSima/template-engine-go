package main

import (
	"fmt"
	"testing"
)

func TestBlockNode(t *testing.T) {
	var testContext = `{"names": [{"name": "johnny", "email": "jbone@gmail.com"}, {"name": "cameron", "email": "boner@gmail.com"}}]}`
	var testSource = `
	{% for n in names %}
		<p>{{ n.name }}</p>
		<p>{{ n.email }}</p>
	 {% endfor %}`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))
	forNode := nodes[1].(ForNode)
	fmt.Println(forNode.Render(c))
}
