#!/bin/bash

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
while ! nc -z postgres 5432; do
  sleep 0.1
done
echo "PostgreSQL is ready!"

# Run migrations
echo "Running migrations..."

goose -dir sql/schema postgres "$DB_URL" up


echo "Migrations completed!" 