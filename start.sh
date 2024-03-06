#!/bin/sh

set -e
. /app/app.env
echo "DB SOURCE"
echo $DB_SOURCE

echo "Start the app"
exec "$@"