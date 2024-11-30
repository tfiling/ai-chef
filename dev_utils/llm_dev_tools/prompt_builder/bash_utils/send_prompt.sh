#!/bin/bash

#Debug flag:
#set -x

# Example usage:
# send_prompt "Hello, this is my prompt"
send_prompt() {
    local text_to_send="$1"

    # Hardcoded URL - modify this to your desired website
    TARGET_URL="https://claude.ai/"

    # Function to check if Chrome is ready
    wait_for_chrome() {
        local window_id
        for i in {1..3000}; do  # Wait up to 30 seconds
            window_id=$(xdotool search --name "Google Chrome" | tail -1)
            if [ ! -z "$window_id" ]; then
                return 0
            fi
            sleep 0.01
        done
        return 1
    }

    # If Chrome isn't running, start it using setsid to detach it
    if ! pgrep -x "chrome" > /dev/null; then
        setsid google-chrome "$TARGET_URL" </dev/null &>/dev/null &
        echo "Starting Chrome..."
    fi

    # Wait for Chrome to open
    if ! wait_for_chrome; then
        echo "Timeout waiting for Chrome to open. Aborting."
        return 1
    fi
    echo "claude tab is ready"

    # Get the window ID
    WINDOW_ID=$(xdotool search --name "Google Chrome" | tail -1)

    # Activate the window
    xdotool windowactivate --sync $WINDOW_ID

    # Paste the text (simulates Ctrl+V)
    echo -n "$text_to_send" | xclip -selection clipboard
    xdotool key ctrl+v
    sleep 0.5
    xdotool key Return

    # Cleanup:
    # Ensure xclip process is terminated
    pkill xclip 2>/dev/null || true

    return 0
}

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 'your prompt text'"
    exit 1
fi

# Execute the function and capture its return value
send_prompt "$1"
exit_code=$?

# Clean exit with proper status code
exit $exit_code