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

var matcher = fmt.Sprintf("(%s.*?%s|%s.*?%s|%s.*?%s)",
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
	lineNumber := 1

	for _, match := range matches {
		end = match[0]

		if match[1] != 0 {
			content := l.template[beg:end]
			lineNumber += strings.Count(content, "\n")
			token := CreateToken(content, false, 0)
			tokens = append(tokens, token)
		}

		beg = match[1]
		tag := l.template[match[0]:match[1]]
		lineNumber += strings.Count(tag, "\n")
		token := CreateToken(tag, true, lineNumber)
		tokens = append(tokens, token)
	}

	if end != len(l.template) {
		piece := l.template[beg:]
		isTag := strings.Contains(piece, "{")
		t := CreateToken(piece, isTag, lineNumber)

		fmt.Println("END")
		fmt.Println(t)

		tokens = append(tokens, t)
	}

	fmt.Println("TOKENS:")
	fmt.Println(tokens)
	return tokens
}

type Lexer struct {
	re       *regexp.Regexp
	template string
}

func NewLexer(template string) Lexer {
	return Lexer{
		regexp.MustCompile(matcher),
		template,
	}
}
