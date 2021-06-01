import os
import divergence

def set_git_repo_path():
    os.environ["GIT_REPO_PATH"] = "/null/path/"


def test_divergence():
    set_git_repo_path()

    report = divergence.main()

    assert type(report) is list
