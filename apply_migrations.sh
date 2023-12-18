#!/bin/bash

# Check if the DATABASE_FILE argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <DATABASE_FILE>"
    exit 1
fi

DATABASE_FILE="$1"
MIGRATIONS_FOLDER="sql/migrations"

# Check if the migrations folder exists
if [ ! -d "$MIGRATIONS_FOLDER" ]; then
    echo "Error: Migrations folder '$MIGRATIONS_FOLDER' not found."
    exit 1
fi

# Find all migration files in the migrations folder
MIGRATION_FILES=$(find "$MIGRATIONS_FOLDER" -type f -name "*.sql" | sort -n)

# Iterate over migration files and apply them in order
for FILE in $MIGRATION_FILES; do
    echo "Applying migration: $FILE"
    sqlite3 "$DATABASE_FILE" < "$FILE"

    # Check the exit status of the previous command
    if [ $? -ne 0 ]; then
        echo "Error applying migration: $FILE"
        exit 1
    fi
done

echo "All migrations applied successfully."
