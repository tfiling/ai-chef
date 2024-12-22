import pathlib
import subprocess
import os
import typing


def run_go_script(script_name: str, args: typing.List[str], base_dir=pathlib.Path(__file__).parent) -> str:
    go_script_path = base_dir / script_name
    bin_file_name = script_name.split(".go")[0]
    binary_path = base_dir / "bin" / bin_file_name

    if _should_compile(binary_path, go_script_path):
        try:
            subprocess.run(
                ["go", "build", "-o", str(binary_path), str(go_script_path)],
                check=True,
                capture_output=True,
                text=True,
                cwd=str(base_dir)
            )
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"Failed to compile Go script: {e.stderr}")

    try:
        result = subprocess.run(
            [str(binary_path)] + args,
            check=True,
            capture_output=True,
            text=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to extract types: {e.stderr}")


def _should_compile(binary_path, go_script_path):
    return (not os.path.exists(binary_path) or
            os.path.getmtime(go_script_path) > os.path.getmtime(binary_path))
