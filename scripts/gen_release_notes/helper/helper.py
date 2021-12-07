"""
This python script contains helper methods that are used in the generate
script, gen_release_notes.py

Note that these are generic helper methods
"""

import re

import generate


def get_weight_change(type_change):
    """
    Depending on the type of change that is passed into this method, returns a corresponding value

    If it cannot find said type in weight_scale, returns a default value of 1
    """

    if generate.weight_scale.get(type_change, None) is not None:
        return generate.weight_scale[type_change]
    return 1


def get_env(string_parser):
    """
    Given a string, checks if that string matches the pattern for a specific
    environment. Generally:
        tags = prod
        RC branch = test
        URL = URL

    Returns said env/is URL. If it cannot find a match, returns None
    """

    if is_regex_match(generate.tag_regex, string_parser) is True:
        return generate.prd_env_name
    if is_regex_match(generate.rc_regex, string_parser) is True:
        return generate.stg_env_name
    if is_regex_match(generate.url_regex, string_parser) is True:
        return "url"

    return None


def is_regex_match(regex_pattern, search_context):
    """
    Wrapper method to perform regex searches on a passed in regex pattern and a string

    If the search context matches the regex_pattern, returns True. Else false.
    """

    regex_search = re.compile(regex_pattern)
    if regex_search.match(search_context) is not None:
        return True
    return False


def is_regex_search(regex_pattern, search_context):
    """
    Wrapper method to check if the regex pattern is within the search context.

    If it is, returns True. False if not.
    """

    regex_search = re.compile(regex_pattern)
    if regex_search.search(search_context) is not None:
        return True
    return False


def check_if_key_in_list_dict(list_obj, dict_key, dict_value):
    """
    Helper method to help check if a specific dictionary is in a list by checking if a specific dict key is
    in said dictionary.

    Returns the index point where this dict was found. If it cannot find it, returns -1
    """

    for index, curr_value in enumerate(list_obj):
        if curr_value.get(dict_key, None) == dict_value:
            return index
    return -1


def add_valid_env_to_set(set_obj, parse_string):
    """
    Taking in a set and a file string, we check to see
    if the parsed string contains the requested env string in it.

    If so, we also check to make sure that the passed in set does NOT contain the
    misc_env value in it. This is done to prioritize environmental changes over specific env changes.
    """

    # Checks if the parsed string has an env path in it
    for curr_env in generate.valid_envs:
        if is_regex_search(curr_env, parse_string):
            if generate.misc_env in set_obj:
                set_obj.remove(generate.misc_env)
            set_obj.add(curr_env)
            return True

    # Checks if the doesn't exist any other higher envs in the set of envs
    found = False
    for curr_set_item in set_obj:
        if curr_set_item in generate.valid_envs and curr_set_item != generate.misc_env:
            found = True
    if found is False:
        set_obj.add(generate.misc_env)
        return True

    return False
