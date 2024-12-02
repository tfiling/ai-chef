project_tree() {
  TREE_EXCLUDE_DIRS=(".git" ".idea" "node_modules" ".venv" "test" "e2e" "docs" "ops" ".ci" ".github" "dev_utils")

  local directory="$1"
  local temp_file="$2"

  if [ ! -d "$directory" ]; then
    echo "$directory does not exist."
    exit 1
  fi

  # Build an array of exclude patterns
  local exclude_args=()
  for exclude_dir in "${TREE_EXCLUDE_DIRS[@]}"; do
    exclude_args+=("-I" "$exclude_dir")
  done

  # Use the array expansion to properly pass arguments to tree
  tree "$directory" "${exclude_args[@]}" > "$temp_file"
}