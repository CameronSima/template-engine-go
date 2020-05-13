package main

import (
	"C"
	"fmt"
)

//export render
func render(templateName *C.char, context *C.char) *C.char {
	source := ReadTemplate(C.GoString(templateName))
	c := C.GoString(context)
	template := NewTemplate(source, NewContext(c))
	result := template.Render()

	fmt.Println("RESULT")
	fmt.Println(result)
	return C.CString(result)
}

func main() {}
