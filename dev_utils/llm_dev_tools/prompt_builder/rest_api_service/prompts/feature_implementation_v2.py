from common import COMMON_POSTFIX, AVOID_EXPLANATIONS, ROLE, NO_TESTS_IMPLEMENTATION, NO_LOGIC_IMPLEMENTATION, \
    ACKNOWLEDGE, COMMON_PREFIX
from rest_api_service.prompts.predefined_prompt import PredefinedPrompt
from rest_api_service.prompts.validators import get_named_args_validator
"""
Prompts for generating code that implements a general feature described - second attempt

The whole flow should include design for the feature, break down implementation into tasks,
implementation of the tasks including tests implementation.
"""

INITIAL_DESIGN_ANALYSIS = PredefinedPrompt(
    name="INITIAL_DESIGN_ANALYSIS",
    prefix=[COMMON_PREFIX],
    postfix=[AVOID_EXPLANATIONS, NO_LOGIC_IMPLEMENTATION],
    prompt_template="Suggest a high level design for {feature_description}\n"
                    "Base your high level design on these types and interfaces:\n{types}\n"
                    "Provide Short description of the main flows to be implemented as part of the feature\n",
    args_validator=get_named_args_validator(["feature_description", "types"]))
BREAK_DOWN_TASKS = PredefinedPrompt(
    name="BREAK_DOWN_TASKS",
    prompt_template="Given the following project structure:\n"
                    "{project_structure}\n"
                    "The Go module name is {go_module_name}\n"
                    "Break that feature implementation into high level essential development tasks.\n"
                    "Avoid any explanations. Generate the list as a yaml code artifact in the following format:\n"
                    "```\n"
                    "tasks:\n"
                    "  - implementation_priority: 1\n"
                    "    name: \"Add user authentication\"\n"
                    "    description: \"Implement JWT-based user authentication\"\n"
                    "    new_files: [\"internal/pkg/llm/claude.go\", \"internal/pkg/llm/types.go\"]\n"
                    "    modified_files: [\"internal/pkg/models/recipe.go\"]\n"
                    "    test_cases: [\"Should validate JWT token\", \"Should reject invalid credentials\"]\n"
                    "```",
    args_validator=get_named_args_validator(["project_structure", "go_module_name"]))
IMPLEMENT_SUBTASK = PredefinedPrompt(
    name="IMPLEMENT_SUBTASK",
    postfix=[COMMON_POSTFIX, NO_TESTS_IMPLEMENTATION],
    prompt_template="The project structure is as follows:\n{project_structure}\n"
                    "The Go module name is {go_module_name}\n"
                    "Based on the above feature description and design implement this subtask:\n{task}",
    args_validator=get_named_args_validator(["project_structure", "go_module_name", "task"]))
GATHER_SUBTASK_TEST_CASES = PredefinedPrompt(
    name="IMPLEMENT_SUBTASK_TEST_CASES",
    prefix=[ROLE],
    postfix=[AVOID_EXPLANATIONS],
    prompt_template="The Go module name is {go_module_name}\n"
                    "Based on the above feature description and design and subtask implementation list required test cases:\n"
                    "{task}\n"
                    "Do not include non-functional edge cases\n"
                    "Each test case should consist of:\n"
                    "- Test case name\n"
                    "- Description of the initial state to be mocked and set up\n"
                    "- What actions should be taken as part of the test case\n"
                    "- What kind of assertions should be made to validate that the tested code works properly",
    args_validator=get_named_args_validator(["go_module_name", "task"]))
TASK_TESTS_IMPLEMENTATION = PredefinedPrompt(
    name="TASK_TESTS_IMPLEMENTATION",
    postfix=[AVOID_EXPLANATIONS],
    prompt_template="In a file named located in a test package implement the above test cases using testify github.com/stretchr/testify\n"
                    "Each test function should be named as follows Test<tested component name>_<tested function>__<use case name>(t *testing.T)\n"
                    "Each test function should have the following structure:\n"
                    "//Arrange\n"
                    "<set up and declare expected mocked store objects behaviour>\n"
                    "//Act\n"
                    "<Execute the tested flow>\n"
                    "//Assert\n"
                    "<Assertion statements using testify>\n")
