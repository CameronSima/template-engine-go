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

	value, _ := c.data.Resolve("user.username")
	value2, _ := c.data.Resolve("user.profile.num_friends")

	assert.Equal(t, "cameron", value)
	assert.Equal(t, "3", value2)

}
