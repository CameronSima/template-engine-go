from ctypes import *
import os
import json

dirname = os.path.dirname(__file__)
filename = os.path.join(dirname, "../../../phantom.so")
lib = cdll.LoadLibrary(filename)
lib.render.restype = c_char_p
    


class Phantom:
    app_dirname = 'templates'

    def __init__(self, params):
        pass

    def get_template(self, template_name):
        return Template(self.app_dirname + "/" + template_name)


class Template:
    def __init__(self, template_name):
        self.template_name = template_name
        self.python_funcs = self.get_python_funcs()

    def render(self, context, request):
        self.original_context = context
        context = self.prepare_context(context, request)
        result = lib.render(
            self.template_name.encode('utf-8'),
            context.encode('utf-8')
        )
        result = json.loads(result)

        print("RESULT IN RENDER")
        print(result)
        return self.post_process(result)

    def post_process(self, result):
        formatted_str = result['result']
        python_nodes = result['functionCalls']
        unresolved_vars = result['unresolvedVariables']
        all_pre_processed = python_nodes + unresolved_vars

        """
        Anything that wasn't resolvable from context (dict values, arrays, etc)
        we'll try to resolve here as a python callable, either from our dict of 
        python functions or the python copy of context (i.e. django ORM object methods).
        """
        processed_values = []

        for node in all_pre_processed:
            print("unprocessed NODE")
            print(node)

            if 'token' in node:
                # is a regular unresolved variable, not a function
                value = safeget(self.original_context, *node['token']['Content'].split('.'))
                print("VALUE")
                print(value)
                processed_values.append(value)
                
            else:
                #func = self.python_funcs[node['functionName']]
                func = self.original_context[node['functionName']]
                processed_values.append(func(*node['parameters']))
        return formatted_str.format(*processed_values)

    def get_python_funcs(self):
        return {}

    def prepare_context(self, context, request=None):
        func_names = {funcName: 1 for funcName, _ in self.python_funcs.items()}
        context = {k: v for k, v in context.items() if is_serializable(v)}
        return json.dumps({**context, **func_names})

def is_serializable(data):
    try:
        json.dumps(data)
        return True
    except (TypeError, OverflowError):
        return False

def safeget(dct, *keys):
    for key in keys:
        try:
            dct = dct[key]
        except TypeError:
            # dont think this even exists in Django, might delte it
            dct = getattr(dct, key.replace('()', ''))

        except KeyError:
            return None
    return dct