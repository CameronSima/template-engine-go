package main

import (
	"fmt"
	"regexp"
	"strings"
)

const FILTER_SEPARATOR = "|"
const FILTER_ARGUMENT_SEPARATOR = ":"
const VARIABLE_ATTRIBUTE_SEPARATOR = "."
const BLOCK_TAG_START = "{%"
const BLOCK_TAG_END = "%}"
const VARIABLE_TAG_START = "{{"
const VARIABLE_TAG_END = "}}"
const COMMENT_TAG_START = "{#"
const COMMENT_TAG_END = "#}"
const TRANSLATOR_COMMENT_MARK = "Translators"
const SINGLE_BRACE_START = "{"
const SINGLE_BRACE_END = "}"

var matcher = fmt.Sprintf("[\\%s\\%s\\%s\\%s\\%s\\%s]+",
	BLOCK_TAG_START,
	BLOCK_TAG_END,
	VARIABLE_TAG_START,
	VARIABLE_TAG_END,
	COMMENT_TAG_START,
	COMMENT_TAG_END)

func (l Lexer) Tokenize() []Token {
	matches := l.re.FindAllStringIndex(l.template, -1)
	tokens := make([]Token, 0)
	beg := 0
	end := 0
	inTag := false
	lineNumber := 1
	taggedToken := make([]string, 0, 2)

	for _, match := range matches {
		end = match[0]
		if match[1] != 0 {
			content := l.template[beg:end]
			lineNumber += strings.Count(content, "\n")

			if inTag {
				taggedToken = append(taggedToken, content)
			} else {
				t := CreateToken([]string{content}, inTag, lineNumber)
				tokens = append(tokens, t)
			}
		}

		beg = match[1]
		tag := l.template[match[0]:match[1]]
		taggedToken = append(taggedToken, tag)

		if len(taggedToken) == 3 {
			t := CreateToken(taggedToken, inTag, lineNumber)
			tokens = append(tokens, t)
			taggedToken = make([]string, 0, 2)
		}
		inTag = !inTag
	}

	if end != len(l.template) {
		t := CreateToken([]string{l.template[beg:]}, inTag, lineNumber)
		tokens = append(tokens, t)
	}

	return tokens
}

type Lexer struct {
	re       *regexp.Regexp
	template string
}

func NewLexer(template string) Lexer {
	l := Lexer{
		regexp.MustCompile(matcher),
		template,
	}

	return l
}
