"""
Test function to check if the methods declared in git_parser work as intended
"""

from unittest.mock import Mock, patch

import pytest
from gen_release_notes import generate
from gen_release_notes.git_parser import git_parser

sample_changed_files = [
    {"repo_owner": "someguy", "repo_name": "verycool", "pr_number": "1337"},
    {"repo_owner": "anotherguy", "repo_name": "alsoverycool", "pr_number": "123457"},
]

sample_good_remote_urls = [
    {"url": "git@github.com:someordinary/repository.git", "owner": "someordinary", "repo": "repository"}
]

sample_bad_remote_urls = [
    {"url": "git@github.com:nope/repository.git", "owner": "someordinary", "repo": "repository"},
    {"url": "git@github.com:someordinary/nope.git", "owner": "someordinary", "repo": "repository"},
    {"url": "https://github.com/someordinary/repository", "owner": "someordinary", "repo": "repository"},
]

sample_pr_information = [
    "Merge pull request #1 from mycoomrepo/someone/myissue",
    "Merge pull request #2 from mycoolrepo/someone/myissue [READY] - [issue] - This is a Pr message",
]

sample_pr_type = [
    {
        "repo_owner": "somecoolowner",
        "repo_name": "anevencoolerrepo",
        "pr_number": "1337",
        "pr_commit_message": "this is a useful message",
    },
    {
        "repo_owner": "AnOwner",
        "repo_name": "ARepoThatIshCooll",
        "pr_number": "1",
        "pr_commit_message": "This - Made - Changes to EVERYTHING",
    },
]

sample_changelog_obj = [{"logs": "this\nis\na\nlist", "start_point": "start", "end_point": "end"}]

sample_modded_dirs = [
    {
        "changed_dirs": ["env/dir_1/changes/", "env/dir_2/changes/", "env/dir_3/changes/"],
        "changed_files": [
            "env/dir_1/changes/text.txt",
            "env/dir_2/changes/code.json",
            "env/dir_3/changes/something.tf",
        ],
    },
    {
        "changed_dirs": [
            "zshing/project/somecode/",
            "blast/from/the/past/",
            "modules/another_changes/",
            "m2/was/here/",
        ],
        "changed_files": [
            "zshing/project/somecode/modules.tf",
            "zshing/project/somecode/s3.tf",
            "zshing/project/somecode/vars.tf",
            "blast/from/the/past/another.tf",
            "blast/from/the/past/something.json",
            "modules/another_changes/woah.txt",
            "modules/another_changes/txt",
            "m2/was/here/sample.txt",
            "m2/was/here/hhhhhhh.yml",
        ],
    },
]


@pytest.mark.parametrize("user_input_changed_files", sample_changed_files)
def test_get_all_changed_files(user_input_changed_files):
    """
    Test for verifying that the changed file function returns something instead of None
    """

    with patch("requests.get"):
        result = git_parser.get_all_changed_files(
            user_input_changed_files["repo_owner"],
            user_input_changed_files["repo_name"],
            user_input_changed_files["pr_number"],
        )
    assert result is not None


@pytest.mark.parametrize("user_input_good_remote_urls", sample_good_remote_urls)
def test_get_repo_information_good(user_input_good_remote_urls):
    """
    Test for verifying that the needed repo information returned is valid and has the expected results
    """

    mocked_repo = Mock()
    mocked_repo_url = Mock()

    mocked_repo.remotes = {"origin": mocked_repo_url}
    mocked_repo_url.url = user_input_good_remote_urls["url"]
    result = git_parser.get_repo_information(mocked_repo)

    assert result is not None
    assert (
        result["owner_name"] == user_input_good_remote_urls["owner"]
        and result["repo_name"] == user_input_good_remote_urls["repo"]
    )


@pytest.mark.parametrize("user_input_bad_remote_urls", sample_bad_remote_urls)
def test_get_repo_information_bad(user_input_bad_remote_urls):
    """
    Test to check when a failed attempt comes through for getting repo information
    """

    mocked_repo = Mock()
    mocked_repo_url = Mock()

    mocked_repo.remotes = {"origin": mocked_repo_url}
    mocked_repo_url.url = user_input_bad_remote_urls["url"]
    result = git_parser.get_repo_information(mocked_repo)

    assert result is not None
    assert (
        result["owner_name"] != user_input_bad_remote_urls["owner"]
        or result["repo_name"] != user_input_bad_remote_urls["repo"]
    )


@pytest.mark.parametrize("user_input_pr_information", sample_pr_information)
def test_get_pr_information(user_input_pr_information):
    """
    Test to check if the information gathered in a PR is valid and the data structures are present
    """

    result = git_parser.get_pr_information(user_input_pr_information)
    assert "pr_number" in result and "pr_title" in result


@pytest.mark.parametrize("user_input_pr_type", sample_pr_type)
def test_get_pr_type(user_input_pr_type):
    """
    Test to confirm that the PR type that is gotten from the function meets the standard
    """

    with patch.object(git_parser, "get_all_changed_files", return_value=[{"filename": "somefile"}]):
        result = git_parser.get_pr_change_type(
            user_input_pr_type["repo_owner"],
            user_input_pr_type["repo_name"],
            user_input_pr_type["pr_number"],
            user_input_pr_type["pr_commit_message"],
        )
        assert "pr_envs" in result and "pr_type" in result


@pytest.mark.parametrize("user_input_changelog_obj", sample_changelog_obj)
@patch.object(git_parser, "get_repo_information", return_value={"owner_name": "someOwner", "repo_name": "someRepo"})
@patch.object(git_parser, "get_pr_information", return_value={"pr_number": "123", "pr_title": "SomeCoolTitle"})
def test_gen_changelog_obj(mocked_pr_info, mocked_repo_info, user_input_changelog_obj):
    """
    Test to confirm that the changelog object is formatted properly
    """

    mock_repo_obj = Mock()
    mock_repo_obj.git = Mock()
    mock_repo_obj.git.log = Mock(return_value=user_input_changelog_obj["logs"])
    result = git_parser.gen_changelog_obj(
        mock_repo_obj, user_input_changelog_obj["start_point"], user_input_changelog_obj["end_point"]
    )

    assert result is not None
    assert (
        generate.dev_env_name in result
        and generate.stg_env_name in result
        and generate.prd_env_name in result
        and generate.misc_env in result
    )


@pytest.mark.parametrize("user_input_modded_dirs", sample_modded_dirs)
def test_get_all_modded_dirs(user_input_modded_dirs):
    """
    Test to get all of the modded dirs and confirm their structure
    """

    mocked_repo = Mock()
    mocked_repo.git = Mock()
    mocked_repo.git.log = Mock(return_value="\n".join(user_input_modded_dirs["changed_dirs"]))
    mocked_repo.git.diff = Mock(return_value="\n".join(user_input_modded_dirs["changed_files"]))

    results = git_parser.gen_all_modded_dirs(mocked_repo, "some_branch")
    assert (
        len(results) == len(user_input_modded_dirs["changed_dirs"])
        and sorted(user_input_modded_dirs["changed_dirs"]) == results
    )
