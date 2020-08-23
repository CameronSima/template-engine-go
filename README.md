# template-engine-go
HTML templating backend for Django written in go.

# build and run:
`go build -o phantom.so -buildmode=c-shared src/go/main.go src/go/template.go src/go/context.go src/go/lexer.go src/go/node.go src/go/parser.go src/go/token.go src/go/utils.go src/go/pythonNode.go src/go/includeNode.go && python src/py/main.py`

### Unit test:
`go test`

#### End to end test:
`python e2e/test.py`

Builds the project and adds it to the virtual env of a small Django test app.
Includes Django and other dependencies (Windows).
