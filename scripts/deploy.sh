#!/bin/sh
set -e

echo "Deploying..."
docker stop portfolio 2>/dev/null || true
docker rm portfolio 2>/dev/null || true

docker run -itd \
  --name portfolio \
  --restart always \
  -p 8000:8000 \
  --network proxy \
  --env-file /root/docker/go-fiber-vue/.env \
  portfolio

docker image prune -f

printf "\nDeploy success"