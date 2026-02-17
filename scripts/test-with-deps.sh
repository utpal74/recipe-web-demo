#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Starting test dependencies via docker compose..."
docker compose up -d

echo "Waiting for services to be healthy..."
tries=0
until docker compose exec -T mongo mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1 && docker compose exec -T redis redis-cli ping >/dev/null 2>&1; do
  tries=$((tries+1))
  if [ "$tries" -ge 30 ]; then
    echo "Timed out waiting for services"
    docker compose logs --no-color
    docker compose down --volumes --remove-orphans
    exit 1
  fi
  sleep 1
done

echo "Services are ready — running tests..."
set +e
go test ./...
TEST_EXIT=$?
set -e

echo "Tearing down test dependencies..."
docker compose down --volumes --remove-orphans

exit $TEST_EXIT
