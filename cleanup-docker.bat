@echo off
echo Stopping containers...
docker stop blogku-mysql blogku-redis blogku-backend blogku-frontend

echo Removing containers...
docker rm blogku-mysql blogku-redis blogku-backend blogku-frontend

echo Removing images...
docker rmi blogku-backend-image blogku-frontend-image

echo Removing network...
docker network rm blogku-network

echo Docker cleanup complete.
