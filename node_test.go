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
	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))
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

		{% for num in names.l %}
			<p>{{ num }}</p>
		{% endfor %}

	 {% endfor %}`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))
	forNode := nodes[1].(ForNode)
	fmt.Println(forNode.Render(c))
}
