package main

// import (
// 	"fmt"
// 	"testing"
// )

// func TestBlockNode(t *testing.T) {
// 	var testContext = `{"username": "cameron"}`
// 	var testSource = `
// 		{% block content %}
// 		<div>
// 			<span>hi</span>
// 			<div>{{ username }}</div>
// 			<p>hi</p>
// 		</div>
// 		{% endblock %}`

// 	lexer := NewLexer(testSource)
// 	parser := NewParser(lexer.Tokenize())
// 	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))

// 	for _, n := range nodes {
// 		fmt.Println("NODE")
// 		fmt.Println(n.Render(NewContext(testContext)))
// 	}

// 	fmt.Println(len(nodes))
// }
