#!/bin/bash
echo "# ThePressUSA Auto-Journalist Bot - Complete Architecture Map" > complete_map.md
echo "" >> complete_map.md
echo "## Codebase Structure & Components" >> complete_map.md
echo "" >> complete_map.md

for file in $(find . -name "*.go" | sort); do
    echo "### File: \`$file\`" >> complete_map.md
    echo "\`\`\`go" >> complete_map.md

    # Extract package name
    grep -E "^package " "$file" | head -n 1 >> complete_map.md
    echo "" >> complete_map.md

    # Extract interfaces
    grep -E "^type [A-Za-z0-9_]+ interface" "$file" >> complete_map.md

    # Extract structs
    grep -E "^type [A-Za-z0-9_]+ struct" "$file" >> complete_map.md

    # Extract functions
    grep -E "^func " "$file" | sed 's/ {.*//' >> complete_map.md

    echo "\`\`\`" >> complete_map.md
    echo "" >> complete_map.md
done
