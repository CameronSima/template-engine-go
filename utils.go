package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func RenderNodeList(nodeList []Node, context Context) string {
	var rendered strings.Builder

	for _, node := range nodeList {

		if node == nil {
			rendered.WriteString("NIL NODE")
		} else {
			rendered.WriteString(node.Render(context))
		}

	}
	return rendered.String()
}

func ReadTemplate(templateName string) string {
	templateName = strings.Trim(templateName, `"'`)
	absPath, _ := filepath.Abs(templateName)
	templateBytes, err := ioutil.ReadFile(absPath)

	if err != nil {
		fmt.Println(err)
	}
	return string(templateBytes)
}

func Contains(l []string, s string) bool {
	for _, v := range l {
		if v == strings.Replace(s, " ", "", -1) {
			return true
		}
	}
	return false
}

type TokenStack struct {
	lock sync.Mutex // may want to add threading later
	tokens []Token
}

func NewTokenStack(tokens []Token) {
	return {sync.Mutex{}, sort.Reverse(tokens)}
}

func (s *TokenStack) NextToken() (Token, error) {
	return s.tokens.pop()
}

func (s *TokenStack) PrependToken(token Token) {
	s.push(token)
}

func (s *TokenStack) pop() (Token, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.tokens)
	if l == 0 {
		return Token{0, "", 0}, errors.New("Stack is empty")
	} 

	result := s.tokens[l-1]
	s.tokens = s.tokens[:l-1]
	return result, nil
}

func (s *TokenStack) push(t Token) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.tokens = append(s.tokens, t)
}