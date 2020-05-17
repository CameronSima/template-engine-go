package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	var testString = `
	{% extends 'test_base.html' %}
	{% block content %}
	<div>
		<span>hi-1</span>
		<div>{{ username }}</div>
            <p>hi-2</p>
	</div>
	 {% endblock %}

	 <div>hi-4</div>

	 {% block footer %}
	 <div>footer</div>
	 {% endblock %}`

	c := NewContext(`{"username": "cameron"}`)
	parser := NewParser(testString, &c)
	nodes := parser.Parse(make([]string, 0))

	for _, n := range nodes {
		fmt.Println("*****************")
		fmt.Println(n)
	}

	fmt.Println(RenderNodeList(nodes, c))

	for _, t := range parser.tokens.tokens {
		fmt.Println("TOKEN")
		fmt.Println(t)
	}

}
