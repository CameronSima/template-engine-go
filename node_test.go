package main

import (
	"fmt"
	"testing"

	"github.com/buger/jsonparser"
)

func TestBlockNode(t *testing.T) {
	var testContext = `{"names": ["cameron", "bob"]}`
	var testSource = `
	{% for name in names %}
		<p>{{ name }}</p>
	 {% endfor %}`

	c := NewContext(testContext)
	parser := NewParser(testSource, &c)
	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))
	forNode := nodes[1].(ForNode)
	fmt.Println(forNode)
	fmt.Println(forNode.loopArrayName)
	fmt.Println(forNode.Render(c))
	//f := c.Resolve("names")

	jsonparser.ArrayEach(c.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		fmt.Println(string(value))
		fmt.Println(jsonparser.Get(value))
	}, "names")

}
