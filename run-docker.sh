#!/bin/bash

# Stop and remove existing containers
docker stop blogku-mysql blogku-redis blogku-backend blogku-frontend || true
docker rm blogku-mysql blogku-redis blogku-backend blogku-frontend || true

# Create network if it doesn't exist
docker network create blogku-network || true

# Build images
echo "Building backend image..."
docker build -t blogku-backend-image ./BackEnd

echo "Building frontend image..."
docker build -t blogku-frontend-image ./frontend

# Run MySQL
echo "Starting MySQL container..."
docker run -d \
  --name blogku-mysql \
  --network blogku-network \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=blogku \
  -e MYSQL_USER=blogku \
  -e MYSQL_PASSWORD=blogku \
  -p 3307:3306 \
  --restart unless-stopped \
  mysql:8.0 \
  --default-authentication-plugin=mysql_native_password

# Run Redis
echo "Starting Redis container..."
docker run -d \
  --name blogku-redis \
  --network blogku-network \
  -p 6380:6379 \
  --restart unless-stopped \
  redis:alpine

# Wait for databases to be ready
echo "Waiting for databases to initialize..."
sleep 10

# Run Backend
echo "Starting backend container..."
docker run -d \
  --name blogku-backend \
  --network blogku-network \
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
  -p 8080:8080 \
  --restart unless-stopped \
  blogku-backend-image

# Run Frontend
echo "Starting frontend container..."
docker run -d \
  --name blogku-frontend \
  --network blogku-network \
  -e NEXT_PUBLIC_API_URL=http://blogku-backend:8080/api/v1 \
  -p 3000:3000 \
  --restart unless-stopped \
  blogku-frontend-image

echo "All containers started. Application should be available at http://localhost:3000"
echo "Backend API at http://localhost:8080"
