import os
import subprocess

DIRECTORY = "/home/galt/code/ai_chef"
EXCLUDE_DIRS = [
    # Generic
    ".git",
    ".idea",
    "node_modules",
    ".venv",
    "test",
    "e2e",
    "docs",
    "ops",
    ".ci",
    ".github",
    # AI Chef specific:
    "dev_utils",
]

# TODO - supporting hiding test files
def generate_tree_structure() -> str:
    if not os.path.isdir(DIRECTORY):
        raise FileNotFoundError(f"{DIRECTORY} does not exist.")

    exclude_pattern = " ".join([f"-I {dir}" for dir in EXCLUDE_DIRS])
    tree_command = f"tree {DIRECTORY} {exclude_pattern}"
    print(tree_command)

    try:
        result = subprocess.run(
            tree_command,
            shell=True,
            check=True,
            capture_output=True,
            text=True
        )
        return result.stdout

    except subprocess.CalledProcessError as ex:
        print(f"Error running tree command: {ex}")
        print(f"Failed to generate tree structure: {ex}")
        exit(1)
    except Exception as ex:
        print(f"Unexpected error: {ex}")
        print(f"Failed to generate tree structure: {ex}")
        exit(1)
