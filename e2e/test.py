import os
import shutil


def main():
    os.system('go build -o phantom.so -buildmode=c-shared src/main.go src/template.go src/context.go src/lexer.go src/node.go src/parser.go src/token.go src/utils.go')
    shutil.copyfile('phantom.so', 'e2e/template_engine_test_app/env/Lib/site-packages/template-engine-go/phantom.so')
    shutil.copyfile('main.py', 'e2e/template_engine_test_app/env/Lib/site-packages/template-engine-go/main.py')

if __name__ == '__main__':
    main()