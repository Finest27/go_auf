#!/bin/bash
OUTPUT="complete_map.md"
echo "# Բոտի Ամբողջական Քարտեզ (Complete Codebase Map)" > $OUTPUT
echo "" >> $OUTPUT

echo "## Պանակների Կառուցվածք (Directory Structure)" >> $OUTPUT
echo "\`\`\`" >> $OUTPUT
find . -type d -not -path "*/\.*" | sort >> $OUTPUT
echo "\`\`\`" >> $OUTPUT
echo "" >> $OUTPUT

echo "## Գո Ֆայլեր և Ֆունկցիաներ (Go Files and Functions)" >> $OUTPUT
for f in $(find . -name "*.go" | sort); do
    echo "### $f" >> $OUTPUT
    echo "\`\`\`go" >> $OUTPUT
    grep -E "^func " "$f" >> $OUTPUT
    echo "\`\`\`" >> $OUTPUT
    echo "" >> $OUTPUT
done

echo "## Տվյալների Բազայի Կառուցվածք (Database Information)" >> $OUTPUT
echo "Բազան SQLite է։ Ստորև ներկայացված են SQL հրամանները՝ աղյուսակներ ստեղծելու համար." >> $OUTPUT
echo "\`\`\`sql" >> $OUTPUT
grep -A 20 -B 2 "CREATE TABLE" $(find . -name "*.go") | grep -v "\-\-" >> $OUTPUT
echo "\`\`\`" >> $OUTPUT

echo "Map generated."
