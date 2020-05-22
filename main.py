from ctypes import *
import os
import json
from base import Phantom as PhantomBase
from base import Template as BaseTemplate
from django.conf import settings
from django.template.engine import Engine
from django.template.backends.base import BaseEngine
from django.urls import get_resolver


class Phantom(BaseEngine, PhantomBase):
    pass


class Template(BaseTemplate):
    def prepare_context(self, request, context):
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