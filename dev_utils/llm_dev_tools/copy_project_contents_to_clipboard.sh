#!/bin/bash

source ./utils/project_tree.sh

DIRECTORY="/home/galt/code/ai_chef"
TEMP_FILE=$(mktemp)

# Define exclusions: List directories and files to exclude
# EXCLUDE_DIRS=(".git" "deploy" "cmd" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github")
# EXCLUDE_FILES=("go.mod" "go.sum" "Dockerfile" "makefile" ".gitignore" "README.md" "*.png" "*.ico" "package-lock.json")
EXCLUDE_DIRS=(".git" "deploy" "cmd" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github" "run_outputs" "__pycache__" ".venv" "llm_dev_tools")
EXCLUDE_FILES=("go.mod" "go.sum" "Dockerfile" "makefile" ".gitignore" "README.md" "*.png" "*.ico" "package-lock.json")


if [ ! -d "$DIRECTORY" ]; then
  echo "$DIRECTORY does not exist."
  exit 1
fi

project_tree "$DIRECTORY" "$TEMP_FILE"

cd "$DIRECTORY" || exit 1

FIND_CMD="find ."
for EXCLUDE_DIR in "${EXCLUDE_DIRS[@]}"; do
    FIND_CMD+=" -not -path '*/${EXCLUDE_DIR}/*' -not -path '*/${EXCLUDE_DIR}'"
done
FIND_CMD+=" -type f"
for EXCLUDE_FILE in "${EXCLUDE_FILES[@]}"; do
    FIND_CMD+=" ! -name '$EXCLUDE_FILE'"
done
FIND_CMD+=" -print0"

# Execute find command and write the contents of the files to the temporary file
eval "$FIND_CMD" | while IFS= read -r -d $'\0' file; do
    echo "File: $file"
    echo "$file" >> "$TEMP_FILE"
    cat "$file" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
done

# Copy the contents to the clipboard
xclip -selection clipboard < "$TEMP_FILE"

echo "Copied $(wc -c < "$TEMP_FILE") characters"

rm "$TEMP_FILE"

cd - > /dev/null || exit 1