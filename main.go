package main

import (
	"C"
	"fmt"
)

//export render
func render(templateName *C.char, context *C.char) *C.char {
	source := ReadTemplate(C.GoString(templateName))
	template := NewTemplate(source)
	c := C.GoString(context)
	result := template.Render(NewContext(c))

	fmt.Println("RESULT")
	fmt.Println(result)
	return C.CString(result)
}

func main() {}
