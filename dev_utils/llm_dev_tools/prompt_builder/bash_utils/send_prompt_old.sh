#!/bin/bash

set -x

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

    # If Chrome isn't running, start it
    if ! pgrep -x "chrome" > /dev/null; then
        google-chrome &
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

#    # Delete previous prompt - unnecessary ATM
#    sleep 0.5
#    xdotool key --clearmodifiers --delay 100 ctrl+a
#    sleep 0.5
#    xdotool key --clearmodifiers --delay 100 Delete
#    sleep 0.5

    # Paste the text (simulates Ctrl+V)
    echo -n "$text_to_send" | xclip -selection clipboard
    xdotool key ctrl+v
    sleep 0.5
    xdotool key Return
    return

    # Clear potential stuck keys - from https://github.com/jordansissel/xdotool/issues/43
#    Commented bc it might get the script stuck
#    sleep 0.5
#    xdotool keyup Meta_L Meta_R Alt_L Alt_R Super_L Super_R
}

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 'your prompt text'"
    exit 1
fi
send_prompt "$1"