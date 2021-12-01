"""
Test function to check if the methods declared in helper work as intended

"""
from unittest.mock import patch

import pytest
from gen_release_notes import generate
from gen_release_notes.helper import helper

sample_type_changes = ["env", "packer", "somemysteryvalue"]

sample_string_parsers = [
    {
        "key": "RC/0.1",
        "expected_value": "stg",
    },
    {"key": "https://github.com/some/repo", "expected_value": "url"},
    {"key": "0.1.0", "expected_value": "prod"},
    {"key": "someunknownvalue", "expected_value": None},
]

sample_regex_match = [
    {"regex": r"^helloworld$", "context": "helloworld", "expected_result": True},
    {"regex": r"^RC\/0\.1$", "context": "RC/0.1", "expected_result": True},
    {"regex": r"^0\.1\.2$", "context": "0.1.2", "expected_result": True},
    {"regex": r"^some\ faulty\ check$", "context": "this is a faulty check", "expected_result": False},
]

sample_regex_search = [
    {"regex": r"string", "context": "string another", "expected_result": True},
    {"regex": r"this", "context": "checkthisrunoutoandonandonandone", "expected_result": True},
    {"regex": r"0\.10\.10", "context": "[WIP] Releasing tag 0.10.10", "expected_result": True},
    {"regex": r"RC\/0\.12", "context": "Nothing to see here", "expected_result": False},
]

sample_check_key_dict = [
    {
        "dicts": [{"something": "value", "aha": "nice"}, {"another": "value", "confusing": "right"}],
        "search_key": "another",
        "search_value": "value",
        "expected_result": 1,
    },
    {
        "dicts": [{"something": "value", "aha": "nice"}, {"another": "value", "confusing": "right"}],
        "search_key": "nada",
        "search_value": "value",
        "expected_result": -1,
    },
    {
        "dicts": [{"something": "value", "aha": "nice"}, {"another": "value", "confusing": "right"}],
        "search_key": "value",
        "search_value": "nice",
        "expected_result": -1,
    },
]

sample_check_valid_env = [
    {"mod_set": set(), "env_check": "dev", "parse_string": "env/dev/someroot/project"},
    {"mod_set": {"general"}, "env_check": "prod", "parse_string": "env/prod/anotherroot/project"},
    {"mod_set": {"dev"}, "env_check": "dev", "parse_string": "module/someenv"},
    {"mod_set": set(), "env_check": "dev", "parse_string": "module/someenv"},
]


@pytest.mark.parametrize("user_type_changes", sample_type_changes)
@patch.object(generate, "weight_scale", return_value={"env": 2, "packer": 20})
def test_get_weight_change(mocked_weight_scale, user_type_changes):
    """
    Test to confirm that the weights are expected
    """

    results = helper.get_weight_change(user_type_changes)
    if mocked_weight_scale().get(user_type_changes) is not None:
        assert results == mocked_weight_scale()[user_type_changes]
    else:
        assert results == 1


@pytest.mark.parametrize("user_string_parser", sample_string_parsers)
@patch("generate.tag_regex", r"^[0-9]{1,}.[0-9]{1,}.[0-9]{1,}$")
@patch("generate.rc_regex", r"^RC/[0-9]{1,}.[0-9]{1,}$")
@patch("generate.url_regex", r"^https://github.com/some/repo")
def test_get_env(user_string_parser):
    """
    Test to confirm that the right env are interpolated
    """

    results = helper.get_env(user_string_parser["key"])
    assert results == user_string_parser["expected_value"]


@pytest.mark.parametrize("user_regex_match", sample_regex_match)
def test_is_regex_match(user_regex_match):
    """
    Test to confirm the regex match works as intended
    """

    results = helper.is_regex_match(user_regex_match["regex"], user_regex_match["context"])
    assert results is user_regex_match["expected_result"]


@pytest.mark.parametrize("user_regex_search", sample_regex_search)
def test_is_regex_search(user_regex_search):
    """
    Test to confirm the regex search works as expected
    """

    results = helper.is_regex_search(user_regex_search["regex"], user_regex_search["context"])
    assert results is user_regex_search["expected_result"]


@pytest.mark.parametrize("user_check_key_dict", sample_check_key_dict)
def test_check_if_key_in_list_dict(user_check_key_dict):
    """
    Test to confirm that the search for a specific dict value works as intended
    """

    results = helper.check_if_key_in_list_dict(
        user_check_key_dict["dicts"], user_check_key_dict["search_key"], user_check_key_dict["search_value"]
    )
    assert results == user_check_key_dict["expected_result"]


@pytest.mark.parametrize("user_check_valid_env", sample_check_valid_env)
@patch.object(generate, "misc_env", return_value="general")
@patch.object(generate, "valid_envs", return_value=["dev", "tst", "prod"])
def test_add_valid_env_to_set(mocked_valid_envs, mocked_misc_env, user_check_valid_env):
    """
    Test ti conform if a valid env is in a set works as intended
    """

    results = helper.add_valid_env_to_set(user_check_valid_env["mod_set"], user_check_valid_env["parse_string"])
    assert user_check_valid_env["env_check"] in mocked_valid_envs()

    # These are the following tests:
    # 1. If the environment to parse is in the modded set
    # 2. If the modded set does NOT have a general environment value and instead has the expected env set
    # 3. If the modded set HAS the general environment value and only said value is in the set
    assert (
        (results is True and user_check_valid_env["env_check"] in user_check_valid_env["mod_set"])
        or (
            results is False
            and mocked_misc_env() not in user_check_valid_env["mod_set"]
            and user_check_valid_env["env_check"] in user_check_valid_env["mod_set"]
        )
        or (
            results is True
            and mocked_misc_env() in user_check_valid_env["mod_set"]
            and len(user_check_valid_env["mod_set"]) == 1
        )
    )
