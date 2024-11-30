import typing
from pathlib import Path

import yaml
from rest_api_service.prompts.predefined_prompt import format_prompt
from utils.send_prompt import send_prompt
from rest_api_service.prompts import crud_api_controller as crud_api_controller_prompts
from utils.project_structure import generate_tree_structure


def run_new_crude_api_controller_flow(**prompt_args):
    # TODO - add step for submitting test report
    # TODO - incorporate linting feedback(some error return values are not handled in the tests)
    # TODO - improve code quality:
    #  use consts when possible
    prompts = [
        crud_api_controller_prompts.DESCRIBE_PROJECT,
        crud_api_controller_prompts.NEW_ENTITY_CONTROLLER_INTERFACE,
        crud_api_controller_prompts.CONTROLLER_IMPLEMENTATION,
        crud_api_controller_prompts.STORE_INTERFACE,
        crud_api_controller_prompts.STORE_IMPLEMENTATION,
        crud_api_controller_prompts.STORE_MOCK,
        crud_api_controller_prompts.CONTROLLER_TEST_CASES,
        crud_api_controller_prompts.CONTROLLER_TESTS_IMPLEMENTATION
    ]
    for predefined_prompt in prompts:
        send_prompt(format_prompt(predefined_prompt, **prompt_args))






def read_input_params_yaml(file_path: str = "./params.yaml") -> typing.Dict[str, typing.Any]:
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
    run_new_crude_api_controller_flow(project_structure=project_structure, **params)
