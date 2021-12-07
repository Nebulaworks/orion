"""
A helper function used to write to a file.
"""

import os

import generate
import jinja2


def write_contents_to_file(changelog_obj, dir_changes_set, cli_args, output_file_path):
    """
    Taking in a changelog object (a dict that contains a list of objects) and a set, we format and parse out the contents
    to a written file.

    Does not return anything.
    """

    jinja_env = jinja2.Environment(
        loader=jinja2.FileSystemLoader(os.path.dirname(os.path.realpath(__file__)) + "/templates")
    )

    top_template = jinja_env.get_template("beginning.j2")
    change_dirs = jinja_env.get_template("changed_roots.j2")
    changes_start_template = jinja_env.get_template("changes_start.j2")
    changes_segments_template = jinja_env.get_template("changes_list.j2")
    bottom_template = jinja_env.get_template("end.j2")

    # if the file already exists, we remove it
    if os.path.isfile(output_file_path):
        os.remove(output_file_path)

    # Writing to file in specific order
    with open(output_file_path, "w", newline="\n") as created_file:
        # Top Section
        created_file.write(
            top_template.render(
                next_branch=cli_args.next_release_tag_branch,
                sprint_number=cli_args.sprint_number,
            )
            + "\n"
            + "\n"
        )

        # Middle Sections
        created_file.write(change_dirs.render(env_name="X") + "\n")
        for curr_dir in dir_changes_set:
            created_file.write("- `{dir}`".format(dir=curr_dir) + "\n")

        created_file.write("\n" + changes_start_template.render() + "\n")

        for curr_env in generate.valid_envs:
            created_file.write("\n" + changes_segments_template.render(env_name=curr_env) + "\n")
            for curr_env_changes in changelog_obj.get(curr_env, []):
                created_file.write("- {change_type}".format(change_type=curr_env_changes["type_change"]) + "\n")
                for curr_entry in curr_env_changes["change_entries"]:
                    created_file.write("    - {entry}".format(entry=curr_entry) + "\n")

        # End Section
        created_file.write(
            "\n"
            + bottom_template.render(branch_name=cli_args.next_release_tag_branch.replace("/", "").replace(".", ""))
            + "\n"
        )

    print("Successfully generated changelog file at {loc}!".format(loc=output_file_path))
