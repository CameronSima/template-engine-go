import unittest
from .base_renderer import Phantom


class TestBase(unittest.TestCase):

    # def test_render(self):
    #     c = {'username': 'cameron'}
    #     p = Phantom({})
    #     p.app_dirname = './base'
    #     t = p.get_template('test_render.html')
    #     r = t.render(c, None)
    #     self.assertEqual(' '.join(r.split()), '<div class="main"> <h1>Hello world</h1> </div>')

    def test_post_process_python_func(self):
        def test_func(param1, param2):
            return param1 + param2

        c = {'param1': 'cameron', 'custom_func': test_func }

        p = Phantom({})
        p.app_dirname = './base'
        t = p.get_template('test_post_process.html')
        r = t.render(c, None)
        self.assertEqual(' '.join(r.split()), '<div> cameronparam2 </div>')

    def test_post_process_callable_in_context(self):
        class TestObject:
            @property
            def callable_property(self):
                return {'nested': 'callable_value'}

        c = {'callable': TestObject()}
        p = Phantom({})
        p.app_dirname = './base'
        t = p.get_template('test_post_process_callable.html')
        r = t.render(c, None)
        print(r)
        print(t)
        self.assertEqual(' '.join(r.split()), '<div> callable_value </div>')


if __name__ == '__main__':
    unittest.main()