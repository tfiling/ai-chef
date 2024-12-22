import pathlib
import typing

import consts
from go_utils.go_script_wrapper import run_go_script


def extract_specific_go_types(requested_types: typing.List[tuple]) -> str:
    res = []
    for file_path, type_name in requested_types:
        type_declaration = run_go_script("get_specific_type_in_file.go", [f"{consts.PROJECT_PATH / file_path}:{type_name}"])
        if not type_declaration:
            raise RuntimeError(f"Could not extract type declaration: {file_path}:{type_name}")
        res.append(f"//{file_path}:\n{type_declaration}")
    return "\n\n".join(res)


def usage_example():
    print([
        ("internal/pkg/llm/claude.go", "ClaudeClient"),
        ("internal/pkg/llm/claude.go", "IClaudeClient"),
    ])

if __name__ == '__main__':
    usage_example()