import os
import subprocess
import shutil

path = os.path.abspath(os.path.join(os.path.dirname(__file__),".."))
phantom_binary_location = os.path.normpath(os.path.join(path, 'phantom.so'))
dependancy_location = os.path.normpath(os.path.join(path, 'e2e/template_engine_test_app/env/Lib/site-packages/template_engine_go'))

if not os.path.exists(dependancy_location):
    try:
        os.makedirs(dependancy_location)
    except OSError as exc: # Guard against race condition
        if exc.errno != errno.EEXIST:
            raise

# build /py subdirectory
if not os.path.exists(os.path.normpath(os.path.join(dependancy_location, 'py'))):
    try:
        os.makedirs(os.path.normpath(os.path.join(dependancy_location, 'py')))
    except OSError as exc: # Guard against race condition
        if exc.errno != errno.EEXIST:
            raise

def main():
    # build go project
    os.system('go build -o phantom.so -buildmode=c-shared src/go/main.go src/go/template.go src/go/context.go src/go/lexer.go src/go/node.go src/go/parser.go src/go/token.go src/go/utils.go src/go/pythonNode.go src/go/includeNode.go')
    shutil.copy2(phantom_binary_location, os.path.normpath(os.path.join(dependancy_location, 'phantom.so')))
    shutil.copy2(os.path.normpath(os.path.join(path, 'src/py/main.py')), os.path.normpath(os.path.join(dependancy_location, 'py/main.py')))

    # run app
    # python_virtualenv_bin = 'C:\\Users\\PcGamer\\dev\\template_engine_go\\e2e\\template_engine_test_app\\env\\Scripts\\pip3.7.exe'
    # subprocess.Popen([python_virtualenv_bin, 'python ' + os.path.normpath(os.path.join(path, 'e2e/template_engine_test_app/test_app/manage.py'))])
    # os.system('python ' + os.path.normpath(os.path.join(path, 'e2e/template_engine_test_app/test_app/manage.py')) + ' manage.py runserver')

if __name__ == '__main__':
    main()