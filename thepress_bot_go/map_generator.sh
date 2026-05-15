#!/bin/bash

MAP_FILE="thepress_bot_go/complete_map.md"

echo "# ThePress Bot - Ամբողջական Քարտեզ (Complete Map)" > $MAP_FILE
echo "Այս ֆայլը պարունակում է բոտի բոլոր ֆայլերի, ֆունկցիաների և տվյալների բազայի (DB) կառուցվածքի մասին տեղեկատվություն։" >> $MAP_FILE
echo "" >> $MAP_FILE

echo "## Տվյալների բազա (Database Schema)" >> $MAP_FILE
echo "Բազան SQLite է, որը գտնվում է \`bot_ultimate.db\` ֆայլում։ Ունի հետևյալ աղյուսակները՝" >> $MAP_FILE
echo '```sql' >> $MAP_FILE
grep -A 20 -E "CREATE TABLE" thepress_bot_go/internal/infra/database/sqlite.go >> $MAP_FILE || true
echo '```' >> $MAP_FILE
echo "" >> $MAP_FILE

echo "## Կոդի կառուցվածք (Code Structure & Functions)" >> $MAP_FILE
for file in $(find thepress_bot_go -type f -name "*.go" | sort); do
    echo "### $file" >> $MAP_FILE
    echo '```go' >> $MAP_FILE
    grep -E "^func " $file >> $MAP_FILE || true
    grep -E "^type .* struct" $file >> $MAP_FILE || true
    grep -E "^type .* interface" $file >> $MAP_FILE || true
    echo '```' >> $MAP_FILE
    echo "" >> $MAP_FILE
done
