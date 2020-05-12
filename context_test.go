package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	var testContext = `
	{
		"user": {
			"username": "cameron", 
			"email": "cjsima@gmail.com",
			"profile": {
				"num_friends": 3
			}
		}
	}
	`
	c := NewContext(testContext)

	value := c.Resolve("user.username")
	value2 := c.Resolve("user.profile.num_friends")

	assert.Equal(t, "cameron", value)
	assert.Equal(t, "3", value2)

}
