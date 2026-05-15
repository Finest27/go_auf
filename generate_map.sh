#!/bin/bash
OUTPUT_FILE="thepress_bot_go/complete_map.md"

echo "# Complete Codebase Map" > $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

find thepress_bot_go -type f -not -path "*/\.git/*" -not -name "*.exe" -not -name "*.exe~" -not -name "*.db" -not -name "*.db-*" -not -name "generate_map.sh" -not -name "complete_map.md" | sort | while read -r file; do
    echo "## $file" >> $OUTPUT_FILE
    echo "\`\`\`go" >> $OUTPUT_FILE
    if [[ $file == *.md || $file == *.html || $file == *.css || $file == *.js || $file == *.json || $file == *.bat || $file == *.sh || $file == *.txt ]]; then
        echo "/* Content of $file */" >> $OUTPUT_FILE
        cat "$file" >> $OUTPUT_FILE
    else
        cat "$file" >> $OUTPUT_FILE
    fi
    echo "\`\`\`" >> $OUTPUT_FILE
    echo "" >> $OUTPUT_FILE
done
