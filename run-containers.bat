@echo off
echo Stopping any existing containers...
docker stop blogku-mysql blogku-redis blogku-backend blogku-frontend 2>NUL
docker rm blogku-mysql blogku-redis blogku-backend blogku-frontend 2>NUL

echo Creating network...
docker network create blogku-network 2>NUL

echo Building backend image...
cd BackEnd
docker build -t blogku-backend-image .
cd ..

echo Building frontend image...
cd frontend
docker build -t blogku-frontend-image .
cd ..

echo Starting MySQL container...
docker run -d --name blogku-mysql ^
  --network blogku-network ^
  -e MYSQL_ROOT_PASSWORD=password ^
  -e MYSQL_DATABASE=blogku ^
  -e MYSQL_USER=blogku ^
  -e MYSQL_PASSWORD=blogku ^
  -p 3307:3306 ^
  mysql:8.0 ^
  --default-authentication-plugin=mysql_native_password

echo Starting Redis container...
docker run -d --name blogku-redis ^
  --network blogku-network ^
  -p 6380:6379 ^
  redis:alpine

echo Waiting for databases to initialize...
timeout /t 5

echo Starting backend container...
docker run -d --name blogku-backend ^
  --network blogku-network ^
  -e DBUSER=root ^
  -e DBPASS=password ^
  -e DBHOST=blogku-mysql ^
  -e DBPORT=3306 ^
  -e DBNAME=blogku ^
  -e JWT_ISSUER=BLOGKU_JWT_ISSUER ^
  -e JWT_SECRET=1234_BLOG_KU_PALING_POPULER_1234 ^
  -e RDSHOST=blogku-redis ^
  -e RDSPORT=6379 ^
  -e ADMIN_API_KEY=super-secret-admin-api-key-change-me ^
  -p 8080:8080 ^
  blogku-backend-image

echo Starting frontend container...
docker run -d --name blogku-frontend ^
  --network blogku-network ^
  -e NEXT_PUBLIC_API_URL=http://blogku-backend:8080/api/v1 ^
  -p 3000:3000 ^
  blogku-frontend-image

echo All containers started. Application should be available at http://localhost:3000
echo Backend API at http://localhost:8080
