#!/bin/bash
OUTPUT="thepress_bot_go/complete_map.md"
echo "# Complete Architecture & Code Map" > $OUTPUT
echo "" >> $OUTPUT
echo "## Directory Structure" >> $OUTPUT
tree thepress_bot_go >> $OUTPUT || find thepress_bot_go >> $OUTPUT
echo "" >> $OUTPUT

echo "## Go Files and Functions" >> $OUTPUT
find thepress_bot_go -name "*.go" | while read -r file; do
    echo "### File: $file" >> $OUTPUT
    echo '```go' >> $OUTPUT
    grep -E "^(package|type|func |func \()" "$file" >> $OUTPUT
    echo '```' >> $OUTPUT
    echo "" >> $OUTPUT
done

echo "## Database Schema (SQLite)" >> $OUTPUT
sqlite3 thepress_bot_go/bot_ultimate.db ".schema" >> $OUTPUT

echo "Map generated at $OUTPUT"
