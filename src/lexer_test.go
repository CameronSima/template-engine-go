package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestTokenize(t *testing.T) {

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

// 	l := NewLexer(testString)
// 	tokens := l.Tokenize()

// 	t2 := tokens[1]

// 	t4 := tokens[3]
// 	t5 := tokens[4]
// 	t6 := tokens[5]
// 	//t7 := tokens[6]
// 	t8 := tokens[7]
// 	//t9 := tokens[8]
// 	t10 := tokens[9]
// 	t11 := tokens[10]
// 	t12 := tokens[11]

// 	assert.Equal(t, len(tokens), 12)
// 	assert.Equal(t, t2.tokenType, 2)
// 	assert.Equal(t, t2.content, `extends 'test_base.html'`)

// 	assert.Equal(t, t4.tokenType, 2)
// 	assert.Equal(t, t4.content, `block content`)

// 	assert.Equal(t, t5.tokenType, 0)

// 	assert.Equal(t, t6.tokenType, 1)
// 	assert.Equal(t, t6.content, `username`)

// 	assert.Equal(t, t8.tokenType, 2)
// 	assert.Equal(t, t8.content, `endblock`)

// 	assert.Equal(t, t10.tokenType, 2)
// 	assert.Equal(t, t10.content, `block footer`)

// 	assert.Equal(t, t11.tokenType, 0)

// 	assert.Equal(t, t12.tokenType, 2)
// 	assert.Equal(t, t12.content, `endblock`)
// }

func TestTokenizeBaseTemplate(t *testing.T) {
	var testString = `
	<div>
		<h1>Main title</h1>
	{% block content %}{% endblock %}
	
		<p>some text</p>
	</div>
	
	{% block footer %}{% endblock %}`

	l := NewLexer(testString)
	tokens := l.Tokenize()

	t1 := tokens[0]
	t2 := tokens[1]
	//t3 := tokens[2]
	t4 := tokens[3]

	// Main title content
	assert.Equal(t, t1.tokenType, 0)
	// block content tag
	assert.Equal(t, t2.tokenType, 2)
	// endblock tag
	assert.Equal(t, t4.tokenType, 2)

}
