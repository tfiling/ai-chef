import json
import typing
from pathlib import Path

import yaml

from rest_api_service.prompts import feature_implementation_v1 as feature_implementation_prompts
from rest_api_service.prompts.predefined_prompt import format_prompt
from utils.fix_broken_tests import fix_tests
from utils.project_structure import generate_tree_structure
from utils.send_prompt import send_prompt
from utils.user_interactions import wait_for_user, ask_yes_no, get_user_input, get_yaml_from_user


# TODO - run tree command as part of prompts generation
# TODO - infra for a predicate for executing the step

def main(**flow_args):
    send_prompt(format_prompt(feature_implementation_prompts.DESCRIBE_FEATURE_STEP_1, **flow_args))
    run_design_step_2(**flow_args)
    tasks = break_down_to_subtasks(**flow_args)
    iterate_feature_tasks(tasks=tasks, **flow_args)


def run_design_step_2(**flow_args):
    if ask_yes_no("Would you like to change suggested API or main flows?"):
        apply_manual_feedback_loop()
        modified_suggested_api = get_user_input("modified_suggested_api") or "You did a great job in this part"
        modified_suggested_main_flows = get_user_input("modified_suggested_main_flows") or "You did a great job in this part"
        send_prompt(format_prompt(feature_implementation_prompts.DESIGN_STEP_1_ADJUSTMENTS,
                                  modified_suggested_api=modified_suggested_api,
                                  modified_suggested_main_flows=modified_suggested_main_flows,
                                  **flow_args))
    send_prompt(format_prompt(feature_implementation_prompts.DESIGN_STEP_2, **flow_args))
    apply_manual_feedback_loop()

def break_down_to_subtasks(**flow_args) -> list:
    send_prompt(format_prompt(feature_implementation_prompts.BREAK_DOWN_TASKS, **flow_args))
    unattended_modified_files, unattended_new_files = detect_unattended_files()
    if unattended_modified_files or unattended_new_files:
        send_prompt(format_prompt(feature_implementation_prompts.INCLUDE_UNATTENDED_FILES,
                                  unattended_new_files=unattended_new_files,
                                  unattended_modified_files=unattended_modified_files,
                                  **flow_args))
    apply_manual_feedback_loop()
    return get_yaml_from_user("Paste implementation tasks here")["tasks"]


def detect_unattended_files():
    unattended_new_files = set(json.loads(get_user_input("Paste a JSON of new files or nothing if list is empty", "[]")))
    unattended_modified_files = set(json.loads(get_user_input("Paste a JSON of modified files or nothing if list is empty", "[]")))
    tasks = get_yaml_from_user("Paste implementation tasks here")["tasks"]
    for task in tasks:
        unattended_new_files.difference_update(task["new_files"])
        unattended_modified_files.difference_update(task["modified_files"])
    return unattended_modified_files, unattended_new_files


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
            fix_tests(base_dir="/home/galt/code/ai_chef")


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
    # This is a first version of feature implementation flow.
    # Current results are too flaky and the flow feels unnatural(Especially when I disagree about the suggested design).
    project_structure = generate_tree_structure()
    params = read_input_params_yaml()
    print(params)
    main(project_structure=project_structure, **params)
