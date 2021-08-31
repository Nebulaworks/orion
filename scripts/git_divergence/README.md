# Git Divergence

Python script to outline *stale/outdated unmerged* remote branches in a given Git Repository. Provides an output of said branches, rating them on a scale of how diverged a branch is (sorted by greatest to least).

## Usage

```
$ GIT_REPO_PATH=/path/to/git/repository python divergence.py
```

### Required Parameters

- `GIT_REPO_PATH`: Path to the repository that will be used by the script. Specified as an env variable.
    - ex: `GIT_REPO_PATH=/path/to/git/repository python divergence.py`

### Optional Parameters

- `GIT_AUTH_BRANCH`: The name of the authoritative branch in a repository. Specified as an env variable. If not set, defaults to `master`.
    - ex: `GIT_REPO_PATH=/path/to/git/repository GIT_AUTH_BRANCH=main python divergence.py`

> All parameters can be specified either in-line or using `export`

## Example Output

<details>
<summary>details</summary>

```
Repo at /samplerepo successfully loaded.
Repo description: Unnamed repository; edit this file 'description' to name the repository.
Remote named "origin" with URL "git@github.com:SampleRepo/Repo.git"
Branch                                         Behind   Ahead   DSB  DSLC DIVERGENCE
refs/stash                                      126       2      10    10        0
0.10.1                                          2879      1     207   207        0
0.11.1                                          2774      2     192   192        0
0.12.1                                          2640      2     175   175        0
0.14.1                                          2337      1     150   150        0
0.15.0                                          2170      2     140   140        0
0.16.0                                          1966      1     123   123        0
0.17.0                                          1845      1     110   110        0
0.17.3                                          1845      4     110   110        0
0.18.1                                          1555      1      95    95        0
0.19.1                                          1362      1      82    82        0
0.20.1                                          1166      1      67    67        0
0.21.1                                          942       1      40    40        0
0.21.2                                          942       2      40    40        0
0.22.1                                          763       1      40    40        0
0.23.1                                          419       1      24    24        0
0.24.0                                          200       3      12    12        0
0.7.1                                           3259      2     262   262        0
0.7.2                                           3259      4     262   262        0
```
</details>

## Testing

```
tox
```