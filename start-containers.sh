#!/bin/bash
# Script to manually manage containers without docker-compose

echo "Removing any previous containers..."
docker rm -f blogku-frontend blogku-backend blogku-mysql blogku-redis 2>/dev/null || true

echo "Creating docker network..."
docker network create blogku-network 2>/dev/null || true

echo "Starting MySQL..."
docker run -d --name blogku-mysql \
  --network blogku-network \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=blogku \
  -e MYSQL_USER=blogku \
  -e MYSQL_PASSWORD=blogku \
  -p 3307:3306 \
  mysql:8.0 \
  --default-authentication-plugin=mysql_native_password

echo "Starting Redis..."
docker run -d --name blogku-redis \
  --network blogku-network \
  -p 6380:6379 \
  redis:alpine

echo "Building and starting Backend..."
docker build -t blogku-backend ./BackEnd
docker run -d --name blogku-backend \
  --network blogku-network \
  -p 8080:8080 \
  -e DBUSER=root \
  -e DBPASS=password \
  -e DBHOST=blogku-mysql \
  -e DBPORT=3306 \
  -e DBNAME=blogku \
  -e JWT_ISSUER=BLOGKU_JWT_ISSUER \
  -e JWT_SECRET=1234_BLOG_KU_PALING_POPULER_1234 \
  -e RDSHOST=blogku-redis \
  -e RDSPORT=6379 \
  -e ADMIN_API_KEY=super-secret-admin-api-key-change-me \
  blogku-backend

echo "Building and starting Frontend..."
docker build -t blogku-frontend ./frontend
docker run -d --name blogku-frontend \
  --network blogku-network \
  -p 3000:3000 \
  -e NEXT_PUBLIC_API_URL=http://blogku-backend:8080/api/v1 \
  -e PORT=3000 \
  blogku-frontend

echo "All containers are started!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"
echo ""
echo "To view logs: docker logs blogku-frontend -f"
