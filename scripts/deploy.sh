#!/bin/sh
set -e

echo "Deploying..."

go mod download
CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

go run ./cmd/console openapi

podman build -t main-app -f Dockerfile .

podman stop main-app 2>/dev/null || true
podman rm main-app 2>/dev/null || true

podman run -itd \
  --name main-app \
  --restart always \
  -p 8080:8080 \
  # --env-file .env \
  main-app

podman image prune -f

printf "\nDeploy success"
