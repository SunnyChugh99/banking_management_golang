#!/bin/sh

set -e
. /app/app.env
echo "DB SOURCE"
echo $DB_SOURCE
echo "Run db migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"