# -*- coding: utf-8 -*-

from distutils.core import Command

from setuptools import find_packages, setup


class DisabledCommands(Command):
    user_options = []

    def initialize_options(self):
        raise Exception("This command is disabled")

    def finalize_options(self):
        raise Exception("This command is disabled")


requirements = [line.strip() for line in open("requirements.txt").readlines()]

setup(
    version="1.0.0",
    name="gen_release_notes",
    description="Easily generate Markdown friendly release notes based on a Git repository",
    author="NWI Engineering",
    author_email="engineering@nebulaworks.com",
    url="https://github.com/Nebulaworks/orion",
    packages=find_packages(exclude=("tests")),
    include_package_data=True,
    install_requires=requirements,
    cmdclass={"register": DisabledCommands, "upload": DisabledCommands},
    entry_points={"console_scripts": ["genreleasenotes = generate:main"]},
    scripts=["generate.py"],
)
