ROLE = "You are an expert Golang developer."
GO_STACK = "You are contributing to a project with the following tech stack: Golang."
PROJECT_DESCRIPTION = "The project is a system for generating and managing cooking recipes."
AVOID_COMMENTS = "Do not include any comments in the generated code."
AVOID_EXPLANATIONS = "Avoid any explanations."

COMMON_PREFIX = f"{ROLE} {GO_STACK} {PROJECT_DESCRIPTION} {AVOID_COMMENTS}"

COMMON_POSTFIX = f"{AVOID_COMMENTS} {AVOID_EXPLANATIONS}"