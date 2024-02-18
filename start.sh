#!/bin/sh

set -e
source /app/app.env
echo "DB SOURCE"
echo $DBSOURCE
echo "Run db migrations"
/app/migrate -path /app/migration -database "$DBSOURCE" -verbose up

echo "Start the app"
exec "$@"