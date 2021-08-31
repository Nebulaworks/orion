# -*- coding: utf-8 -*-


from distutils.core import Command

from setuptools import find_packages, setup


class DisabledCommands(Command):
    user_options = []

    def initialize_options(self):
        raise Exception("This command is disabled")

    def finalize_options(self):
        raise Exception("This command is disabled")


with open("README.md") as f:
    readme = f.read()

# Sometimes requirements in setup.py are supposed to be just whats required to run
# the application. This this case thats teh same as requirements.txt
# For more info: https://packaging.python.org/discussions/install-requires-vs-requirements/
requirements = [line.strip() for line in open("requirements.txt").readlines()]

# Version here doesnt matter much since we are not
# installing this outside of our repo or shipping
# to pypi
setup(
    version="0.2.0",
    name="divergence",
    description="Git Branch Divergence Report",
    long_description=readme,
    author="NWI Engineering",
    author_email="engineering@nebulaworks.com",
    url="https://github.com/Nebulaworks/orion",
    packages=find_packages(exclude=("tests", "docs")),
    install_requires=requirements,
    cmdclass={"register": DisabledCommands, "upload": DisabledCommands},
    scripts=['divergence.py'],
    entry_points={'console_scripts': ['divergence = divergence:main']},
)
