#!/bin/bash

set -e

echo "Starting database container..."
cd .. && docker compose -f docker-compose.yaml up -d db

echo "Building frontend..."
cd frontend && npm run build

echo "Building backend..."
cd ..  && go run .