#!/bin/bash

set -e

echo "Starting database container..."
cd .. && docker compose -f docker-compose.yaml up -d db

echo "Building frontend..."
cd frontend && npm install && npm run build

echo "Starting backend..."
cd .. && go run .