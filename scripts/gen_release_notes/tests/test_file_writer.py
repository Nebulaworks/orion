"""
Test function to check if the methods declared in file_writer work as intended
"""

from unittest.mock import Mock, mock_open, patch

import pytest
from gen_release_notes.file_writer import file_writer

sample_changelog_list = [
    {
        "dev": [
            {
                "type_change": "doc",
                "change_entries": [
                    "env/dev/someproject/something/nice.tf",
                    "docs/something/data.md",
                    "docs/something/another.md",
                ],
            },
            {
                "type_change": "awslambda",
                "change_entries": ["awslambda/anotheroot/changes.py", "awslamda/anothertest/changes.py"],
            },
        ]
    }
]

sample_dir_list = [
    ["dev"],
    ["dev", "tst"],
]

sample_cli_args = [
    {"next_release_tag_branch": "some_bracnch", "sprint_number": "21"},
    {"next_release_tag_branch": "another_branch", "sprint_number": "31"},
]

output_path = "sample.md"


@pytest.mark.parametrize("user_input_changelog_obj", sample_changelog_list)
@pytest.mark.parametrize("user_input_dir_list", sample_dir_list)
@pytest.mark.parametrize("user_input_cli_args", sample_cli_args)
@patch.object(file_writer.jinja2.Environment, "get_template")
def test_write_to_file(mocked_jinjya, user_input_changelog_obj, user_input_dir_list, user_input_cli_args):
    """
    Test to check if the file generation works as intended.

    Asserts if we can mock create the file; contents itself are not an issue
    """

    mock_env = Mock()
    mock_env.next_release_tag_branch = Mock(return_value=user_input_cli_args["next_release_tag_branch"])
    mock_env.sprint_number = Mock(return_value=user_input_cli_args["sprint_number"])

    with patch("builtins.open", mock_open()) as mocked_file:
        file_writer.write_contents_to_file(user_input_changelog_obj, user_input_dir_list, mock_env, output_path)
        mocked_file.assert_called_once_with(output_path, "w", newline="\n")
