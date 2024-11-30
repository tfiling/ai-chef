import subprocess
import sys
import traceback
import typing

from utils.send_prompt import send_prompt


def fix_tests(base_dir):
    try:
        subprocess.run(["/usr/local/go/bin/go", "test", "-v", "./..."],
                                cwd=base_dir,
                                capture_output=True,
                                text=True,
                                check=True,
                                timeout=120)
        print("All tests passed")
    except subprocess.CalledProcessError as ex:
        _fix_tests_using_claude(ex.output)
        print(f"Error executing command: {ex.stderr}", file=sys.stderr)
        sys.exit(1)
    except subprocess.TimeoutExpired:
        print("Bash command got stuck. Continuing.")
        print(traceback.format_exc())
    except FileNotFoundError as ex:
        print(ex)
        sys.exit(1)

def _fix_tests_using_claude(tests_output):
    failed_tests = _extract_failed_tests_outputs(tests_output)
    for test in failed_tests:
        send_prompt(f"Fix this test failure:\n{test}")

def _extract_failed_tests_outputs(tests_output: str) -> typing.List[str]:
    return ["=== RUN" + t for t in tests_output.split("=== RUN") if "--- FAIL" in t]

if __name__ == '__main__':
    fix_tests(base_dir="/home/galt/code/ai_chef")