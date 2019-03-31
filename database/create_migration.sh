#!/usr/bin/env bash
timestamp=`date +%s`

if [ "$#" -gt 1 ]; then
    migration_dir=$1
    filename=$2
else
    migration_dir=sql/migrations
    filename=$1
fi

if [ "$migration_dir" = "" ]; then
    echo "No migration directory selected"
    exit 1
fi

if [ "$filename" = "" ]; then
    echo "No filename specified"
    exit 1
fi

echo "$migration_dir/${timestamp}_${filename}_up.sql"
touch "$migration_dir/${timestamp}_${filename}_up.sql"

echo "$migration_dir/${timestamp}_${filename}_down.sql"
touch "$migration_dir/${timestamp}_${filename}_down.sql"