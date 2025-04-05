#!/bin/bash

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
while ! nc -z postgres 5432; do
  sleep 0.1
done
echo "PostgreSQL is ready!"

# Run migrations
echo "Running migrations..."
# Construct the database URL with the database name
if [ -n "$DB_NAME" ]; then
  # Extract the base URL without any database name or parameters
  BASE_URL=$(echo "$DB_URL" | sed -E 's/\/[^\/]+(\?.*)?$//')
  
  # Construct the full URL with the database name and parameters
  export DB_URL="${BASE_URL}/${DB_NAME}?sslmode=disable"
  echo "Using database URL: $DB_URL"
else
  echo "Warning: DB_NAME is not set. Using default DB_URL."
fi

goose -dir /migrations postgres "$DB_URL" up


echo "Migrations completed!" 