# template-engine-go
HTML templating backend for Django written in go.

# build and run:
`go build -o phantom.so -buildmode=c-shared main.go template.go context.go lexer.go node.go parser.go token.go utils.go && python main.py`

### Unit test:
`go test`

#### End to end test:
`python e2e/test.py`

Builds the project and adds it to the virtual env of a small Django test app.
Includes Django and other dependencies (Windows).
