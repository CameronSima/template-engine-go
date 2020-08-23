from django.apps import AppConfig

class PhantomConfig(AppConfig):
    name = 'phantom'
    verbose_name = 'Phantom template engine'

    def ready(self):
        pass