#!/bin/bash
#TODO - this version opens a fresh claude tab

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

    # Open Chrome with the specific URL
    google-chrome "$TARGET_URL" &
    echo "opened claude tab"

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

    # Wait a bit for the page to load
    sleep 3

    # Delete previous prompt
    xdotool key --clearmodifiers --delay 100 ctrl+a
    sleep 0.5
    xdotool key --clearmodifiers --delay 100 Delete
    sleep 0.5

    # Paste the text (simulates Ctrl+V)
    echo -n "$text_to_send" | xclip -selection clipboard
    xdotool key --clearmodifiers ctrl+v
    xdotool key --clearmodifiers Return
}

[[ $# -eq 0 ]] && { echo "Usage: $0 'your prompt text'"; exit 1; }
send_prompt "$1"