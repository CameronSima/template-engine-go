from ctypes import *
import os
import json
from django.template.engine import Engine
from django.template.backends.base import BaseEngine
from django.urls import get_resolver


dirname = os.path.dirname(__file__)
filename = os.path.join(dirname, "./phantom.so")
lib = cdll.LoadLibrary(filename)
lib.render.restype = c_char_p


class Phantom(BaseEngine):
    app_dirname = 'templates'

    def __init__(self, params):
        pass

    def get_template(self, template_name):
        return Template(self.app_dirname + "/" + template_name)


class Template:
    def __init__(self, template_name):
        self.template_name = template_name

    def render(self, context, request):

        print("\n\nSETTINGS\n\n")
        from django.conf import settings
        print(settings.STATIC_URL)
        print(settings.__dict__)

        #context = json.dumps({ **request.user, **context })
        #context = json.dumps(context)
        context = prepare_context(request, context)


        return lib.render(
            self.template_name.encode('utf-8'),
            context.encode('utf-8')
        )

def prepare_context(request, context):
    result = {
       # 'user': request.user.__dict__,
        'urls': get_urls(),
        'cookies': request.COOKIES,
        'http_host': request.get_host()
    }
    return json.dumps({**result, **context})


def get_urls():
    urls = []
    for u in get_resolver().url_patterns:
        urls.append({
            'name': u.name,
            'pattern': u.pattern._route,
            #'pattern_regex': str(u.pattern.regex)
        })
    return urls



