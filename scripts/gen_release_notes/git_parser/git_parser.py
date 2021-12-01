"""
This module contains all of the methods used to interact with git
"""

import os
import sys

import generate
import requests
from helper import helper

# Secrets
# These are set via ENV variables
pat = os.getenv("GITHUB_TOKEN", None)
if pat is None:
    print("FATAL: Env variable GITHUB_TOKEN is not set. Exiting...")
    sys.exit(1)

pr_numb_regex = r"^#[0-9]+$"
merge_commit_regex = r"^Merge"
cherry_commit_regex = r"\(cherry picked from commit \w+\)"

loading_icons = ["|", "/", "-", "\\"]


def get_all_changed_files(repo_owner, repo_name, pr_number):
    """
    This function grabs all of the changed files that are in a given PR, given the following:
    - repo owner
    - repo name
    - pr number

    Returns a list of all changed files. If it cannot find any,
    """

    payload = {"Authorization": "token {secret}".format(secret=pat)}
    iterator = 1
    result = list()

    while True:

        try:
            res = requests.get(
                "https://api.github.com/repos/{owner}/{repo}/pulls/{pr}/files?per_page=100&page={counter}".format(
                    owner=repo_owner, repo=repo_name, pr=pr_number, counter=iterator
                ),
                headers=payload,
                timeout=2.0,
            )
        except requests.exceptions.RequestException:
            print("Error getting files from PR. Did you pass in the right values? Or are you authenticated?")
            return []

        result.extend(res.json())
        if len(res.json()) == 100:
            iterator += 1
        else:
            break

    return result


def get_repo_information(repo_obj):
    """
    Using the current repository, we obtain information on the repository, which includes
    - owner_name
    - repo_name

    Requires for there to be a remote named "origin" in the current local repo clone

    Returns a dict with said values. If it cannot find both, returns an empty dict
    """

    url = repo_obj.remotes["origin"].url
    colon_find = url.find(":")
    dot_r_find = url.rfind(".")

    if colon_find != -1 and dot_r_find != -1:
        url = url[colon_find + 1 : dot_r_find]
        results = url.split("/")
        return {"owner_name": results[0], "repo_name": results[1]}

    return dict()


def get_pr_information(commit_message):
    """
    Taking in a commit message, extracts out the PR number and PR title it is associated with.

    Returns
        - pr_number
        - pr_title
    in a dict.
    """

    data = dict()

    pr_title_parse = commit_message.split(generate.pr_body_prefix)
    data["pr_title"] = pr_title_parse[len(pr_title_parse) - 1].replace("#", "").strip()

    pr_number_parse = commit_message.split(" ")
    for curr_value in pr_number_parse:
        if helper.is_regex_match(pr_numb_regex, curr_value) is True:
            data["pr_number"] = curr_value.replace("#", "").strip()
            break

    return data


def get_pr_change_type(repo_owner, repo_name, pr_number, pr_commit_message):
    """
    Taking in a commit message, we determine what kind of PR change this message is associated with.

    This attempts to use a weight based system to prioritize specific changes over others, which is determined by the weight_change dict

    Returns the following dict
        {
            pr_envs: [
                list_of_envs
            ],
            pr_type: "string"
        }
    """

    if helper.is_regex_search(cherry_commit_regex, pr_commit_message.lower()):
        return {"pr_envs": [generate.misc_env], "pr_type": "cherry"}

    for curr_value in generate.bug_fix_regex:
        if helper.is_regex_search(curr_value, pr_commit_message.lower()):
            return {"pr_envs": [generate.misc_env], "pr_type": "bug"}

    curr_weight_score = dict()
    curr_parse = {"pr_envs": set()}

    for curr_file_obj in get_all_changed_files(repo_owner, repo_name, pr_number):
        file_name = curr_file_obj["filename"]
        helper.add_valid_env_to_set(curr_parse["pr_envs"], file_name)

        # If a file name has the root env name, we hardcode the index location
        # which should contain a specific root pertaining to a project
        if helper.is_regex_search(generate.root_env_dir_name, file_name):
            root_name = file_name.split("/")[2]
        else:
            root_name = file_name.split("/")[0]

        if curr_weight_score.get(root_name) is None:
            curr_weight_score[root_name] = 1
        else:
            curr_weight_score[root_name] = helper.get_weight_change(root_name) + curr_weight_score[root_name]

    total = 0
    result_type = ""
    for curr_key, curr_value in curr_weight_score.items():
        if curr_value > total:
            total = curr_value
            result_type = curr_key

    return {"pr_envs": list(curr_parse["pr_envs"]), "pr_type": result_type}


def gen_changelog_obj(repo_obj, start_point, end_point):
    """
    Parses the changes that are in between two points of time.
    The resulting object will look like this:

    {
        env_name: [
            {
                type_change: "name",
                change_entries: [
                    "string_1",
                    "string_2"
                ]
            },
            {
                type_change: "name",
                change_entries: [
                    "string_1",
                    "string_2"
                ]
            },
        ]
    }
    """

    results = {generate.dev_env_name: [], generate.stg_env_name: [], generate.prd_env_name: [], generate.misc_env: []}
    repo_data = get_repo_information(repo_obj)
    commit_list = (
        repo_obj.git.log(
            "{branch_1}..{branch_2}".format(branch_1=start_point, branch_2=end_point),
            "--cherry-pick",
            "--first-parent",
            "--format='%s %b'",
        )
        .replace("'", "")
        .split("\n")
    )
    commit_list = [item for item in commit_list if item != ""]

    load_index = 0
    for curr_index, curr_commit in enumerate(commit_list):
        print(
            "Currently parsing changes in repository, please wait. {icon}".format(icon=loading_icons[load_index]),
            end="\r",
            flush=True,
        )

        if helper.is_regex_search(merge_commit_regex, curr_commit) is True:
            if curr_index + 1 < len(commit_list) and helper.is_regex_match(
                cherry_commit_regex, commit_list[curr_index + 1]
            ):
                full_commit = "{orig_mess} {cherry}".format(orig_mess=curr_commit, cherry=commit_list[curr_index + 1])
                pr_info = get_pr_information(full_commit)
            else:
                pr_info = get_pr_information(curr_commit)

            pr_result = get_pr_change_type(
                repo_data["owner_name"], repo_data["repo_name"], pr_info["pr_number"], pr_info["pr_title"]
            )
            change_entry = "{pr_title} [#{pr_number}]({host_url}/{owner}/{repo}/pull/{pr_number})".format(
                pr_title=pr_info["pr_title"],
                pr_number=pr_info["pr_number"],
                host_url=generate.url_repo_endpoint,
                owner=repo_data["owner_name"],
                repo=repo_data["repo_name"],
            )

            for curr_env in pr_result["pr_envs"]:
                index_value = helper.check_if_key_in_list_dict(
                    results.get(curr_env, []), "type_change", pr_result["pr_type"]
                )
                if index_value != -1:
                    results[curr_env][index_value]["change_entries"].append(change_entry)
                else:
                    results[curr_env].append({"type_change": pr_result["pr_type"], "change_entries": [change_entry]})

        load_index += 1
        if load_index == len(loading_icons):
            load_index = 0

    print("")
    print("Successfully obtained changelog object!")
    return results


def gen_all_modded_dirs(repo_obj, end_branch_point):
    """
    Using a repo reference and a branch ref, obtain a list of all of the directories that are
    getting changed, comparing it to the master/generate branch.

    Returns a sorted set of all directories that are being modified
    """

    result = set()
    branch_commit_numb = len(
        repo_obj.git.log(
            "{generate}..{branch}".format(generate=generate.main_branch_name, branch=end_branch_point),
            "--format='%s %b'",
        ).split("\n")
    )
    file_list = repo_obj.git.diff(
        "--dirstat=files,0", "HEAD~{numb}".format(numb=branch_commit_numb), "--name-only"
    ).split("\n")

    for curr_file_path in file_list:
        index = curr_file_path.rfind("/")
        if index != -1:
            modded = curr_file_path[0 : index + 1]
            result.add(modded)
        else:
            result.add(curr_file_path)

    result = sorted(result)
    return result
