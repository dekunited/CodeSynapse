#!/bin/bash

set -euo pipefail

# Cleanup when script exits (normal or Ctrl+C)
cleanup() {
  echo "🧹 Stopping and cleaning up containers..."
  docker compose down
}
trap cleanup EXIT

echo "🔧 Building containers..."
docker compose build frontend
docker compose build backend

echo "🚀 Starting services in dev mode..."
docker compose up

