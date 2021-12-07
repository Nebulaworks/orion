# gen_release_notes Python Module

This directory contains the python code to generate a `markdown` friendly ChangeLog for a git repository. This changelog takes into account the following elements:
- Changed directories in a given branch (usually a branch that is used in a release)
- "Smartly" logging all commits and their respective merge PRs based on deployment environment and project
- Outlining all hotfixes and bugs that came up in said release

<details>
<summary>Example Output</summary>

```bash
$ genreleasenotes 0.28.0 0.29.0 ms/DEV-1449-prd 32
Successfully obtained changelog objecty, please wait. -
Successfully generated changelog file at gen_output.md!

$ cat gen_output.md

## 0.29.0

<details>
  <summary>details</summary>

Release for Sprint 32:

The following environments are updated for X environment:
- `env/prod/someProject`
- `env/prod/someProject2`
- `env/prod/someProject3`

The following roots are skipped in this release:
> Remove this section if irrelevant

The following roots are updated, but not `terraform applied`:
> Remove this section if irrelevant

## New Behavior

### general Changes
- cherry
    - HotFix (cherry picked from commit 01855dcbbc2d90d3c0122d8ecc9830174672e748) [#1](https://github.com/SomeOrg/SomeRepo/pull/1)
- ansible
    - Some change [#2](https://github.com/SomeOrg/SomeRepo/pull/2)
- tests
    - ...
- modules
    - ...
- docs
    - ...
- bug
    - ...

### dev Changes
- terraform
    - ...
- networking
    - ...
- modules
    - ...
- devops
    - ...
- ansible
    - ...

### tst Changes
- shared-services
    - ...
- modules
    - ...

### prod Changes
- shared-services
    - ...
- master
    - ...
- terraform
    - ...
- modules
    - ...
- networking
    - ...
- devops
    - ...
- security
    - ...
- ansible
    - ...

[Return to Change Header](#0290)

[Return to Top](#Changelog)

</details>
```
</details>

## Pre-Requirements
- Python Version 3.8.X

## How to Install

### Download Via Clone Repo

1. Clone repository and navigate to this root
2. `pip install .`

### Download Via Repo Link

1. `pip install -e "git+https://github.com/Nebulaworks/orion.git@<commit_hash>#egg=genreleasenotes&subdirectory=scripts/gen_release_notes"`

## How to Use

1. Create a new file, `.gen_release.config`, which contains configuration for this script
    - Details on this config can be seen [here](#configuration)
    - A sample configuration can be found [here](./gen_release_sample.config)

2. Export the following environment variables:

```bash
# By default, the location for this file is in `~/.config/.gen_release.config`
# this can be omitted if one is ok with this default path
$ export CONFIG_LOC="path/to/.gen_release.config"

# This token should have repository read access
$ export GIT_TOKEN="pat_value" 
```

3. Invoke the command

```bash
# Help output
$ genreleasenotes -h
usage: genreleasenotes [-h] last_release_tag_branch next_release_tag_branch release_changes_branch sprint_number

positional arguments:
  last_release_tag_branch
                        tag or branch name of the starting point to check for changes
  next_release_tag_branch
                        tag or branch name of the ending point to stop checking for changes
  release_changes_branch
                        the name of the branch that is used for determing the changed directories for a release
  sprint_number         the current sprint number these notes were generated from

optional arguments:
  -h, --help            show this help message and exit

# Example Usage
$ genreleasenotes RC/0.28 RC/0.29 ms/DEV-1449-tst 32
```

## Configuration

The nature of this script is to be as flexible as however many use cases there are. As such, this script has an external configuration file that can alter the behavior of the changelog generation.

Below is a list of all of the configurable values

| Variable     | Description     | Example Value |
|--------------|-----------------|---------------|
| tag\_regex | String value to check for tags in a repo. Note that special characters (i.e. `.`) do not need escape characters. |  `"^[0-9]{1,}.[0-9]{1,}.[0-9]{1,}$"` |
| rc\_regex | String value to check for release candidate (e.g. specific branches) in a repo. Note that special characters (i.e. `.`) do not need escape characters. | `"^RC/[0-9]{1,}.[0-9]{1,}$"` | =
| bug\_fix\_regex | List of values that correspond to a bug fix. | `["fix","hotfix","bug"]` |
| main\_branch\_name | Name of the main/master branch in a repo. | `"master"` |
| tst\_env\_name | Name of test environment. | `"tst"` |
| prd\_env\_name | Name of production environment. | `"prd"` |
| dev\_env\_name | Name of development environment. | `"dev"` |
| misc\_env\_name | Name of general/non-environment specific changes. | `"general"` |
| root\_env\_dir\_name | Name of the directory that corresponds to the start of the environment roots | `"env"` |
| repo\_url | Repository to analyze with this script. | `"https://github.com/Nebulaworks/orion"` |
| pr\_body\_prefix | String to determine what pattern to look for in pull request titles. | `"] -"` |
| output\_file\_name | Path to output the results. | `"gen_output_results.md"`|
| weight\_scale | Scale to determine change priority. More details [below](#weight-scale). | `["env": "1", "modules": "10", "docs": "20"]` |

### Weight Scale

This script determines change types based on a general weight scale. This weight scale is numeric and determines change types across your entire repository.

Weights are calculated on the summation of change types within a pull request. These change types are determined by the root directory that a change is located in. For example, suppose we have a PR that has the following changes:

```
env/dev/account_name/global/main.tf
docs/document_1/README.md
awslambda/team/lambda/lambda.py
```

The change type for each of the files is determined as such:

```
env -> account_name
docs
awslambda
```

> For roots that pertain to environmental changes (i.e. env/dev), a special parameter is used to specify these type of roots: `root_env_dir_name`. These will take the **second** subdirectory that is in that path which uses that particular specified root name.

Each of these chage types are then weighted based on the weight scale that is passed into the config. A good default set can be seen in the sample config. Nevertheless, these values are highly subjective and can produce different outputs based on your own use case.

## Testing

To test this python code, `tox` can be leveraged:

```
tox
```