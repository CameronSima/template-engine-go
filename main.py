from ctypes import *
import json


lib = cdll.LoadLibrary("./phantom.so")
lib.render.restype = c_char_p


class Phantom:
    def get_template(self, template_name):
        return Template(template_name)


class Template:
    def __init__(self, template_name):
        self.template_name = template_name

    def render(self, context, request):
        context = json.dumps({ **request, **context })
        return lib.render(
            self.template_name.encode('utf-8'),
            context.encode('utf-8')
        )


c = {'username': 'cameron'}
r = {'GET': {'param': 'hi'}}
p = Phantom()
t = p.get_template('test.html')
t.render(c, r)


