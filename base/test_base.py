import unittest
from base import Phantom

print("hi")
class TestBase(unittest.TestCase):

    def test_render(self):
        c = {'username': 'cameron'}
        p = Phantom({})
        p.app_dirname = './'
        t = p.get_template('test_render.html')
        r = t.render(c, None)
        self.assertEqual(' '.join(r.split()), '<div class="main"> <h1>Hello world</h1> </div>')

    def test_post_process(self):
        c = {'param1': 'cameron'}

        p = Phantom({})
        p.app_dirname = './'
        t = p.get_template('test_post_process.html')
        r = t.render(c, None)
        print(r)
        print(t)
        self.assertEqual(' '.join(r.split()), '<div> cameronparam2 </div>')
    


if __name__ == '__main__':
    unittest.main()