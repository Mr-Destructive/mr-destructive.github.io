#!/bin/bash

# Directory containing the files
DIRECTORY="posts/"

find "$DIRECTORY" -type f -name "*.md" | while read -r file; do
  content=$(<"$file")

  # Check if the file contains YAML frontmatter (starting with ---)
  if [[ "$content" =~ ^--- ]]; then
    # Extract the first instance of YAML frontmatter (removes the --- lines)
    yaml_frontmatter=$(echo "$content" | sed -n '/^---$/,/^---$/p' | sed '1d;$d')

    # Convert YAML to JSON using yq
    json_frontmatter=$(echo "$yaml_frontmatter" | yq eval -o=json -)

    # Remove only the first YAML frontmatter block and the surrounding ---
    content=$(echo "$content" | sed '1,/^---$/d')

    # Add the JSON frontmatter and the rest of the content
    new_content="$json_frontmatter\n$content"

    # Write the new content back into the file
    echo -e "$new_content" > "$file"
  fi
done
