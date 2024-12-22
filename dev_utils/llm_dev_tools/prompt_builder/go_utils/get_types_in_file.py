import pathlib
import subprocess
import os
import typing

import consts
from go_utils.go_script_wrapper import run_go_script


def extract_go_types_from_files(file_paths: typing.List[str]) -> str:
    """
    Extract all type declarations from specified Go files.

    Args:
        file_paths (List[str]): List of file paths relative to project_path

    Returns:
        str: All extracted type declarations, separated by newlines

    Raises:
        RuntimeError: If extraction fails for any file
    """
    res = []
    for file_path in file_paths:
        type_declarations = run_go_script("get_types_in_file.go", [str(consts.PROJECT_PATH / file_path)])
        if not type_declarations:
            raise RuntimeError(f"Could not extract type declarations from: {file_path}")
        res.append(f"//{file_path}:\n{type_declarations}")
    return "\n\n".join(res)


def usage_example():
    print(extract_go_types_from_files([
        "internal/pkg/llm/claude.go",
        "internal/pkg/llm/parser.go",
        "internal/pkg/llm/recipe_generator.go",
    ]))


if __name__ == '__main__':
    usage_example()