#!/bin/sh
set -e

echo "Deploying..."

go mod download
CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

go run ./cmd/console openapi

podman build -t main-app -f Dockerfile .

if [ ! -f .env ]; then
  echo "Missing .env in $(pwd); refusing to run because podman would create a directory instead." >&2
  exit 1
fi

podman stop main-app 2>/dev/null || true
podman rm main-app 2>/dev/null || true

podman run -itd \
  --name main-app \
  --restart always \
  -p 8080:8080 \
  -v "$PWD/.env:/app/.env:ro,Z" \
  main-app

podman image prune -f

printf "\nDeploy success"
