"""
This is a python approach to the problem of doing the following:
- Grab all of the file changes
- Grab the difference between two points of time
- Generate a markdown file that contains this automated change log

The goal of this is to make this as universal as possible, but also keep it Python native
"""

import argparse
import json
import os
import sys

from file_writer import file_writer
from git import Repo
from git_parser import git_parser

# Global Configuration
try:
    custom_path = os.path.expanduser("~") + "/.config/.gen_release.config"
    file_path = os.getenv("CONFIG_LOC", custom_path)
    with open(file_path, "r") as config_file:
        json_obj = json.load(config_file)
except FileNotFoundError:
    print("Error, there doesn't exist a configuration file at {path}".format(path=file_path))
    print("Place a config file in said path or specify a custom path via setting CONFIG_LOC.")
    print("A sample config can be obtained in Nebulaworks/orion/scripts/gen_release_notes/gen_release_sample.config")
    sys.exit(1)

main_branch_name = json_obj["main_branch_name"]
dev_env_name = json_obj["dev_env_name"]
stg_env_name = json_obj["stg_env_name"]
prd_env_name = json_obj["prd_env_name"]
misc_env = json_obj["misc_env_name"]
valid_envs = [misc_env, dev_env_name, stg_env_name, prd_env_name]

root_env_dir_name = json_obj["root_env_dir_name"]

repo_url = json_obj["repo_url"]
url_repo_endpoint = repo_url.replace("https://", "").split("/")[0]

pr_body_prefix = json_obj["pr_body_prefix"]
output_file_name = json_obj["output_file_name"]
weight_scale = json_obj["weight_scale"]

tag_regex = r"{}".format(json_obj["tag_regex"]).replace(".", r"\.")
rc_regex = r"{}".format(json_obj["rc_regex"]).replace(".", r"\.")
url_regex = r"^{url}".format(url=url_repo_endpoint)
bug_fix_regex = json_obj["bug_fix_regex"]


def main():
    """
    Main Invoker of this script.

    Does not return anything.
    """

    parser = argparse.ArgumentParser()
    parser.add_argument("last_release_tag_branch", help="tag or branch name of the starting point to check for changes")
    parser.add_argument(
        "next_release_tag_branch", help="tag or branch name of the ending point to stop checking for changes"
    )
    parser.add_argument(
        "release_changes_branch",
        help="the name of the branch that is used for determing the changed directories for a release",
    )
    parser.add_argument("sprint_number", help="the current sprint number these notes were generated from")
    py_args = parser.parse_args()

    repo = Repo(repo_url)
    dir_changes_list = git_parser.gen_all_modded_dirs(repo, py_args.release_changes_branch)
    change_log_obj = git_parser.gen_changelog_obj(
        repo, py_args.last_release_tag_branch, py_args.next_release_tag_branch
    )
    file_writer.write_contents_to_file(change_log_obj, dir_changes_list, py_args, "gen_output.md")


if __name__ == "__main__":
    main()
