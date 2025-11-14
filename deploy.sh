#!/bin/bash
echo "Deploying OmeTV Pro..."
docker-compose down
docker-compose up --build -d
sleep 30
docker-compose up --force-recreate certbot
echo "Deployed! Visit https://client.yourdomain.com"