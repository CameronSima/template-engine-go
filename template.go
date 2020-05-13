package main

type Template struct {
	parser  Parser
	source  string
	context *Context
}

func (t Template) Render() string {

	nodes := t.parser.Parse(make([]string, 0), 0, len(t.parser.tokens))

	// fmt.Println("Context")
	// fmt.Println(t.context.render_context)
	return RenderNodeList(nodes, *t.context)
}

func NewTemplate(source string, context Context) Template {
	parser := NewParser(source, &context)
	return Template{parser, source, &context}
}
