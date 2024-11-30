#!/bin/bash

DIRECTORY="/home/galt/code/ai_chef"
TEMP_FILE=$(mktemp)

TREE_EXCLUDE_DIRS=(".git" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github" "dev_utils")

# Define files to copy
FILES=("dev_utils/local_dev_env/compose.yaml" "Dockerfile" "makefile")

# Define directories to copy
DIRECTORIES=()
#DIRECTORIES=("deploy/fake_client/fake_data"
#"internal/app/webserver/controllers")

cd "$DIRECTORY" || exit 1

# Generate the exclude pattern for tree command
EXCLUDE_PATTERN=""
for EXCLUDE_DIR in "${TREE_EXCLUDE_DIRS[@]}"; do
    EXCLUDE_PATTERN+="-I $EXCLUDE_DIR "
done

tree "$DIRECTORY" $EXCLUDE_PATTERN > "$TEMP_FILE"

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "File: $file" >> "$TEMP_FILE"
        cat "$file" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    else
        echo "File not found: $file"
    fi
done

for directory in "${DIRECTORIES[@]}"; do
    if [ -d "$directory" ]; then
        # Recursively copy files in the directory
        while IFS= read -r -d '' file; do
            echo "File: $file" >> "$TEMP_FILE"
            cat "$file" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
        done < <(find "$directory" -type f -print0)

        echo "" >> "$TEMP_FILE"
    else
        echo "Directory not found: $directory"
    fi
done

# Copy the contents to the clipboard
xclip -selection clipboard < "$TEMP_FILE"

echo "Copied $(wc -c < "$TEMP_FILE") characters"

rm "$TEMP_FILE"

cd - > /dev/null || exit 1