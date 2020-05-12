package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {

	var testString = `
	<div>
		<span>hi</span>
		<div> {{ user }} </div> 
		<p>{% hell0 %} </p>
	</div>`

	l := NewLexer(testString)
	tokens := l.Tokenize()

	assert.Equal(t, len(tokens), 5)

}
