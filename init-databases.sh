#!/bin/bash
set -e
set -o pipefail

# Debug information
echo "Using POSTGRES_USER: $POSTGRES_USER"
echo "Using POSTGRES_HOST: $POSTGRES_HOST"

# Wait for postgres to be ready
until PGPASSWORD=$POSTGRES_PASSWORD pg_isready -h $POSTGRES_HOST -U "$POSTGRES_USER"; do
  echo "Waiting for postgres..."
  sleep 1
done

echo "Postgres is ready, creating databases..."

# Create databases
for i in {1..5}; do
  echo "Creating database chirpy$i..."
  # Check if database exists
  if PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "chirpy$i"; then
    echo "Database chirpy$i already exists"
  else
    PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U "$POSTGRES_USER" -c "CREATE DATABASE chirpy$i;"
    if [ $? -eq 0 ]; then
      echo "Database chirpy$i created successfully"
    else
      echo "Failed to create database chirpy$i"
    fi
  fi
done 