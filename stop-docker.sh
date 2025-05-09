#!/bin/bash

# Stop containers
echo "Stopping containers..."
docker stop blogku-mysql blogku-redis blogku-backend blogku-frontend

# Optional: Remove containers
# echo "Removing containers..."
# docker rm blogku-mysql blogku-redis blogku-backend blogku-frontend

echo "Containers stopped."
