package main

type Template struct {
	parser Parser
	source string
}

func (t Template) Render(context Context) string {
	nodes := t.parser.Parse(make([]string, 0), 0, len(t.parser.tokens))
	return RenderNodeList(nodes, context)
}

func NewTemplate(source string) Template {
	lexer := NewLexer(source)
	parser := Parser{lexer.Tokenize(), "main parser"}
	return Template{parser, source}
}
