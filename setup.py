import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="phantom_templates_django_cameron_sima", # Replace with your own username
    version="0.0.1",
    author="Cameron Sima",
    author_email="cjsima@gmail.com",
    description="Fast Django template renderer written in Golang",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/CameronSima/template-engine-go",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)