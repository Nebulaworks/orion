import datetime
import os
import re
import sys
import argparse

from git import Repo

MASTER = os.getenv("GIT_AUTH_BRANCH", default="master")

DATETIME = "%Y-%m-%d %H:%M:%S%z"

ap = argparse.ArgumentParser()
ap.add_argument("-p", "--pattern", required=False, help="regex pattern for BRANCH_FILTER")

args = vars(ap.parse_args())

if args["pattern"] is not None:
    # A filter patern was passed in
    pattern = re.compile(args["pattern"])
    BRANCH_FILTER = pattern
else:
    BRANCH_FILTER = r"(.*RC\/.*)|([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)"


def get_first_commit(repo, branch):

    other_shas = set()
    first_commit = None

    for parent_commit in repo.iter_commits(MASTER):
        other_shas.add(parent_commit.hexsha)

    for commit in repo.iter_commits(branch):
        if commit.hexsha not in other_shas:
            first_commit = commit

    return first_commit


def get_divergence_report(repo):

    report = []

    for head in repo.refs:

        br = str(head)

        if re.match(BRANCH_FILTER, br):
            continue

        a, b, dslc, dsb, d = divergence_by_branch(repo, br)

        if a:
            row = {
                "Branch": br,
                "Behind": b,
                "Ahead": a,
                "DSB": dsb,
                "DSLC": dslc,
                "Divergence": d,
            }
            report.append(row)

    return sorted(report, key=lambda k: k["Divergence"], reverse=True)


def get_max_len_of_col(repo_list, col_name):

    top_length = 0
    for _, curr_element in enumerate(repo_list, 0):
        curr_len = len(curr_element.get(col_name, ""))
        if curr_len > top_length:
            top_length = curr_len

    return top_length


def divergence_by_branch(repo, branch):

    latest_commit = repo.commit(branch)
    first_commit = get_first_commit(repo, branch)

    a = list(repo.iter_commits(MASTER + "@{u}.." + branch))
    b = list(repo.iter_commits(branch + ".." + MASTER + "@{u}"))

    c = latest_commit.authored_datetime
    d = datetime.datetime.now(tz=c.tzinfo)

    DSLC = d - c

    e = repo.commit(first_commit)

    DSB = d - e.authored_datetime

    D = DSB.days * (DSB.days - DSLC.days)

    return len(a), len(b), DSLC.days, DSB.days, D


def print_commit(commit):
    print(str(commit.hexsha))
    print(
        '"{}" by {} ({})'.format(
            commit.summary, commit.author.name, commit.author.email
        )
    )
    print(str(commit.authored_datetime))
    print(
        str("Commit count: {} and Commit size: {}".format(commit.count(), commit.size))
    )


def print_heads(repo):

    for head in repo.heads:
        print('Repo Head named "{}"'.format(head))


def print_repository(repo):
    print("Repo description: {}".format(repo.description))

    for remote in repo.remotes:
        print('Remote named "{}" with URL "{}"'.format(remote, remote.url))


def main():

    repo_path = os.getenv("GIT_REPO_PATH")

    if MASTER != "master":
        print(
            "GIT_AUTH_BRANCH has been set, using {branch} instead of master...".format(
                branch=MASTER
            )
        )
    try:
        repo = Repo(repo_path)
    except:
        return []

    if repo.bare:
        print("Could not load repository at {} :(".format(repo_path))

    repo.config_reader()
    print("Repo at {} successfully loaded.".format(repo_path))

    print_repository(repo)

    report = get_divergence_report(repo)

    branch_name_len = str(get_max_len_of_col(report, "Branch"))
    table = "{:" + branch_name_len + "} {:^8} {:^8} {:^4} {:4} {:8}"

    print(table.format("Branch", "Behind", "Ahead", "DSB", "DSLC", "DIVERGENCE"))

    for b in report:
        print(
            table.format(
                b["Branch"],
                b["Behind"],
                b["Ahead"],
                b["DSB"],
                b["DSLC"],
                b["Divergence"],
            )
        )


if __name__ == "__main__":

    main()
