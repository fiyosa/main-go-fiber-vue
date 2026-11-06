#!/bin/sh
set -e

echo "Deploying..."
docker stop main-app 2>/dev/null || true
docker rm main-app 2>/dev/null || true

docker run -itd \
  --name main-app \
  --restart always \
  -p 8000:8000 \
  --env-file .env \
  main-app

docker image prune -f

printf "\nDeploy success"
