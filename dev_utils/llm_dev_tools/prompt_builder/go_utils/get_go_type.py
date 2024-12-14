import pathlib
import subprocess
import os
import typing


def extract_go_types(project_path: pathlib.Path, requested_types: typing.List[tuple]) -> str:
    res = []
    for file_path, type_name in requested_types:
        type_declaration = _extract_go_type(project_path / file_path, type_name)
        if not type_declaration:
            raise RuntimeError(f"Could not extract type declaration: {file_path}:{type_name}")
        res.append(f"//{file_path}:\n{type_declaration}")
    return "\n\n".join(res)


def _extract_go_type(file_path: pathlib.Path, type_name: str) -> str:
    """
    Extract a type declaration from a Go file using the Go type extractor.

    Args:
        file_path (Path): Path to the Go file
        type_name (str): Name of the type to extract

    Returns:
        str: The extracted type declaration

    Raises:
        FileNotFoundError: If the Go file or extractor directory doesn't exist
        subprocess.CalledProcessError: If the Go script fails to execute
    """
    extractor_dir = pathlib.Path(__file__).parent

    go_script_path = extractor_dir / "get_go_type.go"
    binary_path = extractor_dir / "bin" / "get_go_type"

    if _should_compile(binary_path, go_script_path):
        try:
            subprocess.run(
                ["go", "build", "-o", str(binary_path), str(go_script_path)],
                check=True,
                capture_output=True,
                text=True,
                cwd=str(extractor_dir)
            )
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"Failed to compile Go script: {e.stderr}")

    try:
        result = subprocess.run(
            [str(binary_path), f"{file_path}:{type_name}"],
            check=True,
            capture_output=True,
            text=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to extract type: {e.stderr}")

def _should_compile(binary_path, go_script_path):
    return (not os.path.exists(binary_path) or
            os.path.getmtime(go_script_path) > os.path.getmtime(binary_path))


def usage():
    print(extract_go_types(pathlib.Path("/home/galt/code/ai_chef"), [
        ("internal/pkg/llm/claude.go", "ClaudeClient"),
        ("internal/pkg/llm/claude.go", "IClaudeClient"),
    ]))

if __name__ == '__main__':
    usage()