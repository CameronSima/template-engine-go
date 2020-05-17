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

func CreateToken(tokenStr string, isTag bool, lineNumber int) Token {
	var token Token

	if isTag {
		tokenStr = strings.TrimSpace(tokenStr)
		openingTag := tokenStr[0:2]
		content := tokenStr[2:len(tokenStr)-2]
		content = strings.TrimSpace(content)

		switch openingTag {
		case VARIABLE_TAG_START:
			token = Token{TOKEN_VAR, content, lineNumber}

		case BLOCK_TAG_START:
			token = Token{TOKEN_BLOCK, content, lineNumber}

		case COMMENT_TAG_START:
			token = Token{TOKEN_COMMENT, "", lineNumber}
		}
	} else {
		token = Token{TOKEN_TEXT, tokenStr, lineNumber}
	}
	return token
}
