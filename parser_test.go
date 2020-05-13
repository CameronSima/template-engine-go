package main

// import (
// 	"fmt"
// 	"testing"
// )

// func TestParser(t *testing.T) {

// 	var testString = `
// 	{% extends 'test_base.html' %}
// 	{% block content %}
// 	<div>
// 		<span>hi-1</span>
// 		<div>{{ username }}</div>
//             <p>hi-2</p>
// 	</div>
// 	 {% endblock %}

// 	 <div>hi-3</div>

// 	 {% block footer %}
// 	 <div>footer</div>
// 	 {% endblock %}`

// 	c := NewContext("{}")
// 	parser := NewParser(testString, &c)
// 	nodes := parser.Parse(make([]string, 0), 0, len(parser.tokens))

// 	for _, n := range nodes {
// 		fmt.Println("*****************")
// 		fmt.Println(n)
// 	}

// }
