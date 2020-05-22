package main

import (
	"C"
	"fmt"
)
import "encoding/json"

type Result struct {
	NeedsPostProcessing bool         `json:"needsPostProcessing"`
	FunctionCalls       []PythonNode `json:"functionCalls"`
	Result              string       `json:"result"`
}

func (r Result) Stringify() string {
	result, _ := json.Marshal(r)
	return string(result)
}

//export render
func render(templateName *C.char, context *C.char) *C.char {
	source := ReadTemplate(C.GoString(templateName))
	c := C.GoString(context)
	template := NewTemplate(source, NewContext(c))
	//result := template.Render()
	resultStr := template.Render()
	result := Result{
		len(template.context.functionCalls) > 0,
		template.context.functionCalls,
		resultStr,
	}

	fmt.Println("RESULT")
	fmt.Println(result.Stringify())
	return C.CString(result.Stringify())
	//return C.CString(result)
}

func main() {}
