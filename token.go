package main

import (
	"strings"
)

const (
	TOKEN_TEXT    = 0
	TOKEN_VAR     = 1
	TOKEN_BLOCK   = 2
	TOKEN_COMMENT = 3
)

type Token struct {
	tokenType int
	content   string
	//position   []int
	lineNumber int
}

func CreateToken(bits []string, inTag bool, lineNumber int) Token {
	var token Token

	if inTag {
		content := ""

		if len(bits) > 1 {
			content = strings.TrimSpace(bits[1])
		}

		openingTag := strings.Replace(bits[0], " ", "", -1)

		switch openingTag {
		case VARIABLE_TAG_START:
			token = Token{TOKEN_VAR, content, lineNumber}

		case BLOCK_TAG_START:
			token = Token{TOKEN_BLOCK, content, lineNumber}

		case COMMENT_TAG_START:
			token = Token{TOKEN_COMMENT, "", lineNumber}
		}
	} else {
		token = Token{TOKEN_TEXT, bits[0], lineNumber}
	}
	return token
}
