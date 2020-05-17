# template-engine-go
HTML templating backend for Django written in go.

run:
go build -o phantom.so -buildmode=c-shared main.go template.go context.go lexer.go node.go parser.go token.go utils.go && python main.py
