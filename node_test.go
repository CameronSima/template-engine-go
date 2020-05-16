package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestBlockNode(t *testing.T) {
// 	var testContext = `{"names": [{"name": "johnny", "email": "jbone@gmail.com", "l": [1, 2, 3]}, {"name": "cameron", "email": "boner@gmail.com", "l": [4, 5]}}]}`
// 	var testSource = `
// 	{% for n in names %}
// 		<p>is first? {{ first }}</p>
// 		<p>{{ n.name }}</p>
// 		<p>{{ n.email }}</p>

// 	 {% endfor %}`

// 	c := NewContext(testContext)
// 	parser := NewParser(testSource, &c)
// 	nodes := parser.Parse(make([]string, 0))
// 	forNode := nodes[1].(ForNode)
// 	strippedResult := strings.Replace(forNode.Render(c), "\n", "", -1)
// 	strippedResult = strings.Replace(strippedResult, " ", "", -1)
// 	assert.Equal(t, strippedResult, "\t\t<p>isfirst?true</p>\t\t<p>johnny</p>\t\t<p>jbone@gmail.com</p>\t\t\t<p>isfirst?false</p>\t\t<p>cameron</p>\t\t<p>boner@gmail.com</p>\t")
// }

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

	assert.Equal(t, nestedFNode.token.content, "for num in n.l")
	assert.Equal(t, nestedFNode.loopVariable, "num")
	assert.Equal(t, nestedFNode.loopArrayName, "n.l")

	fmt.Println("Child NODE")
	fmt.Println(nestedFNode)
	fmt.Println(RenderNodeList(nodes, c))

	// c.AddToContextData(`{"name": "johnny", "email": "jbone@gmail.com", "l": ["1", "2", "3"]}`, "n")
	// fmt.Println("Rendering nested for loop variable:")
	// fmt.Println(nestedFNode.Render(c))

}
