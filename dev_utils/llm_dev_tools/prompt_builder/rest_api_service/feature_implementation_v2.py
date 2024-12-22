import typing
from pathlib import Path

import yaml

import consts
from go_utils.get_specific_type_in_file import extract_specific_go_types
from go_utils.get_types_in_file import extract_go_types_from_files
from rest_api_service.prompts import feature_implementation_v2 as feature_implementation_prompts
from rest_api_service.prompts.predefined_prompt import format_prompt
from utils.fix_broken_tests import fix_tests
from utils.project_structure import generate_tree_structure
from utils.send_prompt import send_prompt
from utils.user_interactions import wait_for_user, ask_yes_no, get_yaml_from_user


# TODO - run tree command as part of prompts generation
# TODO - infra for a predicate for executing the step

def main(**flow_args):
    declared_types = _extract_type_declarations(**flow_args)
    if not declared_types:
        raise RuntimeError("Could not extract relevant types")
    send_prompt(format_prompt(feature_implementation_prompts.INITIAL_DESIGN_ANALYSIS, **flow_args, types=declared_types))
    tasks = break_down_to_subtasks(**flow_args)
    iterate_feature_tasks(tasks=tasks, **flow_args)

def _extract_type_declarations(**flow_args) -> str:
    res = ""
    if flow_args.get("relevant_specific_types", []):
        res += extract_specific_go_types(flow_args.get("relevant_specific_types", []))
        res += "\n"
    if flow_args.get("relevant_types", []):
        res += extract_go_types_from_files(flow_args.get("relevant_types", []))
        res += "\n"
    return res


def break_down_to_subtasks(**flow_args) -> list:
    send_prompt(format_prompt(feature_implementation_prompts.BREAK_DOWN_TASKS, **flow_args))
    apply_manual_feedback_loop()
    return get_yaml_from_user("Paste implementation tasks here")["tasks"]


def iterate_feature_tasks(tasks: typing.List[dict], **flow_args):
    tasks = sorted(tasks, key=lambda t: t["implementation_priority"])
    print("Generated tasks:")
    print("\n".join((f"{task["implementation_priority"]}){task["name"]}" for task in tasks)))
    for task in tasks:
        print(f"implementing {task["name"]}")
        if not ask_yes_no(f"Would you like to implement this feature:\n"
                   f"{yaml.dump(task, sort_keys=False, default_flow_style=False)}"):
            # USer decided not to implement task
            continue
        wait_for_user("Focus on Claude's input. Override any previous task implementation using message edit")
        send_prompt(format_prompt(feature_implementation_prompts.IMPLEMENT_SUBTASK,
                                  task=str(task), **flow_args), should_wait_user=False)
        apply_manual_feedback_loop()
        if task["test_cases"]:
            send_prompt(format_prompt(feature_implementation_prompts.GATHER_SUBTASK_TEST_CASES,
                                      task=str(task), **flow_args))
            send_prompt(format_prompt(feature_implementation_prompts.TASK_TESTS_IMPLEMENTATION, **flow_args))
            apply_manual_feedback_loop()
            fix_tests(base_dir=consts.PROJECT_PATH)


def apply_manual_feedback_loop():
    # The user should use the chat for small changes and click OK when completed
    wait_for_user("Apply feedback manually if needed. Click OK when completed.")


def read_input_params_yaml(file_path: str = "./params.yaml") -> typing.Dict[str, typing.Any]:
    # TODO - support default params values(required for modified_suggested_api and modified_suggested_main_flows)
    path = Path(file_path).resolve()
    try:
        with open(path, 'r', encoding='utf-8') as yaml_file:
            return yaml.safe_load(yaml_file)
    except FileNotFoundError:
        raise FileNotFoundError(f"YAML file not found at: {path}")
    except yaml.YAMLError as ex:
        raise yaml.YAMLError(f"Error parsing YAML file: {ex}")
    except PermissionError:
        raise PermissionError(f"Permission denied accessing file: {path}")


if __name__ == '__main__':
    project_structure = generate_tree_structure()
    params = read_input_params_yaml()
    print(params)
    main(project_structure=project_structure, **params)
