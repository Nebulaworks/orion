# -*- coding: utf-8 -*-

from setuptools import setup, find_packages

# Sometimes requirements in setup.py are supposed to be just whats required to run
# the application. This this case thats teh same as requirements.txt
# For more info: https://packaging.python.org/discussions/install-requires-vs-requirements/
requirements = [line.strip() for line in open("requirements.txt").readlines()]

# Version here doesnt matter much since we are not
# installing this outside of our repo or shipping
# to pypi
setup(
    version='0.1.0',
    name='flasksample',
    description='flask sample app',
    author='NWI Engineering',
    author_email='engineering@nebulaworks.com',
    url='https://github.com/Nebulaworks/orion',
    packages=find_packages(exclude=('tests', 'docs')),
    entry_points={'console_scripts': ['flasksample = flasksample.app:main']},
    install_requires=requirements
)
