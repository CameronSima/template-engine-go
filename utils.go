package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

func RenderNodeList(nodeList []Node, context Context) string {
	var rendered strings.Builder

	for _, node := range nodeList {
		rendered.WriteString(node.Render(context))	
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
		if v == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}

type TokenStack struct {
	lock    *sync.Mutex // may want to add threading later
	tokens  []Token
	IsEmpty bool
}

func NewTokenStack(tokens []Token) TokenStack {
	s := TokenStack{
		&sync.Mutex{},
		tokens,
		len(tokens) < 0,
	}
	s.reverse()
	return s
}

func (s *TokenStack) NextToken() (Token, error) {
	return s.pop()
}

func (s *TokenStack) PrependToken(token Token) {
	s.push(token)
}

func (s *TokenStack) pop() (Token, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.tokens)
	if l == 1 {
		s.IsEmpty = true
	} else if l == 0 {
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

func (s *TokenStack) reverse() {
	for i, j := 0, len(s.tokens)-1; i < j; i, j = i+1, j-1 {
		s.tokens[i], s.tokens[j] = s.tokens[j], s.tokens[i]
	}
}
