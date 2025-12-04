#!/bin/bash

# Script to rebuild and restart the backend service

echo "Stopping services..."
docker compose down

echo "Rebuilding server..."
docker compose build server

echo "Starting services..."
docker compose up -d

echo "Waiting for services to be ready..."
sleep 5

echo "Server logs:"
docker logs dating_server --tail 50

echo ""
echo "Done! Server is running on http://localhost:8080"
echo "Health check: http://localhost:8080/health"
