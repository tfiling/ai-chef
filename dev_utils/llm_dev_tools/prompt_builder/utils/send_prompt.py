import pathlib
import subprocess
import sys
import traceback

from utils.user_interactions import wait_for_user


def send_prompt(text: str):
    """
    Sends a prompt to Claude via the send_prompt bash script

    Args:
        text (str): The prompt text to send
    """
    send_prompt_path = pathlib.Path(__file__).parent.parent.absolute() / "bash_utils" / "send_prompt.sh"

    try:
        wait_for_user("Focus on Claude's input\nClick OK to continue")
        result = subprocess.run([send_prompt_path, text],
                                capture_output=True,
                                text=True,
                                check=True,
                                timeout=10)
        print(result.stdout)
    except subprocess.TimeoutExpired:
        print("Bash script got stuck. Continuing.")
        print(traceback.format_exc())
    except subprocess.CalledProcessError as e:
        print(f"Error executing script: {e.stderr}", file=sys.stderr)
        sys.exit(1)
    except FileNotFoundError:
        print("Error: send_prompt.sh script not found", file=sys.stderr)
        sys.exit(1)
