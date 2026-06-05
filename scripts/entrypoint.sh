#!/bin/sh
set -e

echo "Running migrations..."
./migrate up

echo "Starting application..."
exec ./app -storage=postgres