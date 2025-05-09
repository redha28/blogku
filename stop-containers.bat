@echo off
echo Stopping containers...
docker stop blogku-mysql blogku-redis blogku-backend blogku-frontend

echo Containers stopped.
