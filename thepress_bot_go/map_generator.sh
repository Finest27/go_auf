#!/bin/bash
OUTPUT="complete_map.md"
echo "# Complete Architecture Map" > $OUTPUT
echo "## Database Schema" >> $OUTPUT
sqlite3 bot_ultimate.db .schema >> $OUTPUT

echo "" >> $OUTPUT
echo "## Directory Structure" >> $OUTPUT
tree -L 4 >> $OUTPUT

echo "" >> $OUTPUT
echo "## Go Files, Structs, and Functions" >> $OUTPUT
for file in $(find . -type f -name "*.go" | grep -v "vendor"); do
    echo "### File: $file" >> $OUTPUT
    echo '```go' >> $OUTPUT
    grep -E '^(type|func) ' "$file" >> $OUTPUT
    echo '```' >> $OUTPUT
    echo "" >> $OUTPUT
done

echo "Map generation complete."
