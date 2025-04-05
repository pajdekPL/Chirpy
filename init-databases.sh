#!/bin/bash

# Wait for postgres to be ready
until PGPASSWORD=$POSTGRES_PASSWORD pg_isready -h postgres -U "$POSTGRES_USER"; do
  echo "Waiting for postgres..."
  sleep 1
done

# Create databases
for i in {1..5}; do
  PGPASSWORD=$POSTGRES_PASSWORD psql -h postgres -U "$POSTGRES_USER" -c "CREATE DATABASE IF NOT EXISTS chirpy$i;"
done 