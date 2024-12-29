"""
Prompts for generating code that implements a general feature described - first attempt

The whole flow should include design for the feature, break down implementation into tasks,
implementation of the tasks including tests implementation.
"""
from common import COMMON_POSTFIX, AVOID_EXPLANATIONS, ROLE, NO_TESTS_IMPLEMENTATION, NO_LOGIC_IMPLEMENTATION, \
    ACKNOWLEDGE
from rest_api_service.prompts.predefined_prompt import PredefinedPrompt
from rest_api_service.prompts.validators import get_named_args_validator

DESCRIBE_FEATURE_STEP_1 = PredefinedPrompt(
    name="DESCRIBE_FEATURE_STEP_1",
    postfix=[COMMON_POSTFIX, NO_LOGIC_IMPLEMENTATION],
    prompt_template="Suggest a high level design for {feature_description}\n"
                    "Your high level design should include the following sections:\n"
                    "1. API exposed by new components\n"
                    "2. Short description of the main flows to be implemented as part of the feature\n",
    args_validator=get_named_args_validator(["feature_description"]))
DESIGN_STEP_1_ADJUSTMENTS = PredefinedPrompt(
    name="DESIGN_STEP_2_WITH_ADJUSTMENTS",
    postfix=[COMMON_POSTFIX, NO_LOGIC_IMPLEMENTATION, ACKNOWLEDGE],
    prompt_template="I went over the suggested API and flows and changed the scope and implementation details a little bit:\n"
                    "1. API exposed by new component:\n"
                    "{modified_suggested_api}\n"
                    "2. Short description of the main flows:\n"
                    "{modified_suggested_main_flows}\n",
    args_validator=get_named_args_validator(["modified_suggested_main_flows",
                                             "modified_suggested_api"]))
DESIGN_STEP_2 = PredefinedPrompt(
    name="DESIGN_STEP_2_ACCEPTING_STEP_1",
    postfix=[COMMON_POSTFIX, NO_LOGIC_IMPLEMENTATION],
    prompt_template="Great job!\n"
                    "Provide these additional design details based on the new scope and implementation details:\n"
                    "3. List files to be created in the format: [\"internal/pkg/llm/claude.go\", \"internal/pkg/llm/types.go\"]\n"
                    "4. List files to be modified in the format: [\"internal/pkg/llm/claude.go\", \"internal/pkg/llm/types.go\"]\n"
                    "5. Implementation of new types that are mentioned in the chat but not yet implemented")
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
INCLUDE_UNATTENDED_FILES = PredefinedPrompt(
    # There might be a situation where not all the new files or modified files mentioned in the design are included in the tasks breakdown
    name="INCLUDE_UNATTENDED_FILES",
    prompt_template="Given the following project structure:\n"
                    "{project_structure}\n"
                    "The Go module name is {go_module_name}\n"
                    "Update the tasks or create new tasks if necessary attending these files mentioned in the above "
                    "design but were not included in the subtasks breakdown:\n"
                    "Files that should be modified according to the design but are not included in the above subtasks "
                    "breakdown: {unattended_modified_files}\n"
                    "new files not included in the above subtasks breakdown: {unattended_new_files}\n"
                    "Avoid any explanations. Generate the complete tasks list as a yaml artifact in the following format:\n"
                    "```\n"
                    "tasks:\n"
                    "  - implementation_priority: 1\n"
                    "    name: \"Add user authentication\"\n"
                    "    description: \"Implement JWT-based user authentication\"\n"
                    "    new_files: [\"internal/pkg/llm/claude.go\", \"internal/pkg/llm/types.go\"]\n"
                    "    modified_files: [\"internal/pkg/models/recipe.go\"]\n"
                    "    test_cases: [\"Should validate JWT token\", \"Should reject invalid credentials\"]\n"
                    "```",
    args_validator=get_named_args_validator(
        ["project_structure", "go_module_name", "unattended_new_files", "unattended_modified_files"]))
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
