from common import COMMON_POSTFIX, AVOID_EXPLANATIONS, ACKNOWLEDGE
from rest_api_service.prompts.predefined_prompt import PredefinedPrompt
from rest_api_service.prompts.validators import get_named_args_validator


def __get_entity_model():
    """Expects the user to past the type declaration of the model."""
    print("Paste new entity model here(avoid empty lines):")
    lines = []
    while True:
        line = input()
        if line == '':  # Empty line to stop
            break
        lines.append(line)
    return {"model": "\n".join(lines)}


# TODO - extract tech stack: Mongo, testing 3rd-party pkg
DESCRIBE_PROJECT = PredefinedPrompt(
    name="DESCRIBE_PROJECT",
    postfix=f"{ACKNOWLEDGE}",
    prompt_template="You are an expert Golang developer. "
                    "You are contributing to a project with the following tech stack: Golang, MongoDB, Docker. "
                    "The project is a system for generating and managing cooking recipes. "
                    "The project structure is as follows:\n{project_structure}\n"
                    "The Go module name is {go_module_name}\n"
                    "My next message will include further instructions.",
    args_validator=get_named_args_validator(["project_structure", "go_module_name"]))
NEW_ENTITY_CONTROLLER_INTERFACE = PredefinedPrompt(
    name="NEW_ENTITY_CONTROLLER_INTERFACE",
    postfix=COMMON_POSTFIX,
    prompt_template="I would like to implement a new controller Go type that would expose CRUD API "
                    "for an entity named {entity_name}. "
                    "The controller Go type is declared as follows:\n```"
                    "type {entity_name}Controller struct {{\n"
                    "{entity_name}Store store.I{entity_name}Store\n"
                    "}}```\n"
                    "Name the new file <entity name>_controller.go\n"
                    "Implement a method with the following signature:\n```"
                    "func (c *{entity_name}Controller) RegisterRoutes(router fiber.Router) error {{\n"
                    "router.<Router method>(<const declaring the API endpoint>, c.<handler method>)\n"
                    "<additional declarations>\n"
                    "}}```.\n"
                    "{entity_name}'s backend model is as follows: ```{model}```\n"
                    "Do not implement the methods declared in RegisterRoutes method.\n",
    args_validator=get_named_args_validator(["entity_name", "model"]))
CONTROLLER_IMPLEMENTATION = PredefinedPrompt(
    name="CONTROLLER_IMPLEMENTATION",
    postfix=COMMON_POSTFIX,
    prompt_template="Implement the methods used in func (c *{entity_name}Controller) RegisterRoutes(router fiber.Router) error",
    args_validator=get_named_args_validator(["entity_name"]))
STORE_INTERFACE = PredefinedPrompt(
    name="STORE_INTERFACE",
    postfix=COMMON_POSTFIX,
    prompt_template="In a file named <entity name>_store.go located in store package declare an interface for store.I{entity_name}Store",
    args_validator=get_named_args_validator(["entity_name"]))
STORE_IMPLEMENTATION = PredefinedPrompt(
    name="STORE_IMPLEMENTATION",
    postfix=COMMON_POSTFIX,
    prompt_template="In a file named <entity name>_store.go located in store package implement new Go type named Mongo{entity_name}Store implementing the above interface store.I{entity_name}Store",
    args_validator=get_named_args_validator(["entity_name"]))
STORE_MOCK = PredefinedPrompt(
    name="STORE_MOCK",
    postfix=COMMON_POSTFIX,
    prompt_template="Implement new Go type named MockI{entity_name}Store implementing the above interface store.I{entity_name}Store\n"
                    "MockI{entity_name}Store should embed mock.Mock from \"github.com/stretchr/testify/mock\"\n",
    args_validator=get_named_args_validator(["entity_name"]))
CONTROLLER_TEST_CASES = PredefinedPrompt(
    name="CONTROLLER_TEST_CASES",
    postfix=AVOID_EXPLANATIONS,
    prompt_template="List test cases for {entity_name}Controller based on the routes declared in RegisterRoutes\n"
                    "For each tested method exposing an API route list its specific test cases including edge cases\n"
                    "Do not include non-functional edge cases\n"
                    "Do not include integration tests since the implementation will be using a mock\n"
                    "Each test case should consist of:\n"
                    "- Test case name\n"
                    "- Description of the initial state to be mocked and set up\n"
                    "- What actions should be taken as part of the test case\n"
                    "- What kind of assertions should be made to validate that the tested code works properly",
    args_validator=get_named_args_validator(["entity_name"]))
CONTROLLER_TESTS_IMPLEMENTATION = PredefinedPrompt(
    name="CONTROLLER_TESTS_IMPLEMENTATION",
    postfix=AVOID_EXPLANATIONS,
    prompt_template="In a file named <entity name>_controller_test.go located in controllers_test package implement the above test cases for {entity_name}Controller\n using testify github.com/stretchr/testify\n"
                    "Each test function should be named as follows Test{entity_name}Controller_<method name>__<use case name>(t *testing.T)\n"
                    "Each test function should setup a fiber test app and register the {entity_name}Controller's API endpoints as part of Arrange part\n"
                    "Each test function should have the following structure:\n"
                    "//Arrange\n"
                    "<set up and declare expected mocked store objects behaviour>\n"
                    "//Act\n"
                    "<Execute the tested flow>\n"
                    "//Assert\n"
                    "<Assertion statements using testify>\n",
    args_validator=get_named_args_validator(["entity_name"]))
