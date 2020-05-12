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
