# BlogKu CMS Platform

A full-stack content management system built with Next.js for the frontend and Go for the backend API.

## Project Overview

BlogKu is a modern blog content management system that consists of:
- **Frontend**: Next.js 15 application with responsive design and admin dashboard
- **Backend**: Golang API with MySQL database and Redis caching

## Features

- Public blog listing and detailed blog post views
- Admin authentication with JWT
- Admin dashboard for content management
- Create, read, update, and delete blog posts
- Image upload for blog posts
- Automatic slug generation
- Redis caching for improved performance
- Responsive design

## Tech Stack

### Frontend
- Next.js 15
- React 19
- TypeScript
- Tailwind CSS
- Shadcn UI components

### Backend
- Go (Golang)
- Gin web framework
- MySQL database
- Redis for caching
- JWT for authentication
- Swagger for API documentation

## Setup Instructions

### Option 1: Using Docker Compose (Recommended)

The easiest way to run the entire application stack is using Docker Compose.

#### Prerequisites
- Docker and Docker Compose installed on your system
- Command line/terminal access

#### Installation
1. Clone the repository
2. Navigate to the project root directory
3. Run the following command:

```bash
docker-compose up -d
```

This will start the following services:
- MySQL database on port 3307
- Redis on port 6380
- Backend API on port 8080
- Frontend application on port 3000

You can access the application at http://localhost:3000

### Option 2: Manual Setup

#### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn
   # or
   pnpm install
   ```

3. Run the development server:
   ```bash
   npm run dev
   # or
   yarn dev
   # or
   pnpm dev
   ```

4. Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

#### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd BackEnd
   ```

2. Configure environment variables in `.env` file
3. Run the database migration script:
   ```bash
   mysql -u root -p < migrations/schema.sql
   ```

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Run the application:
   ```bash
   go run cmd/main.go
   ```

## API Documentation

### Authentication
- **POST /api/v1/auth/login** - Admin login
- **POST /auth/logout** - Admin logout

### Public Endpoints
- **GET /blogs** - Get all published blogs with pagination
- **GET /blogs/{slug}** - Get a specific blog by slug

### Admin Endpoints (require authentication)
- **POST /admin/blogs** - Create a new blog post
- **PUT /admin/blogs/{id}** - Update a blog post
- **DELETE /admin/blogs/{id}** - Delete a blog post

For full API documentation, you can access:
- **Postman Collection**: [BlogKu API Collection](https://www.postman.co/workspace/My-Workspace~25b50f28-dc7e-4fd2-88a5-21905a54fec1/collection/22450553-d11d4131-db6e-4335-aca3-1b6c97e01892?action=share&creator=22450553)

## Admin Credentials

For admin access, use the following credentials:
- **Email**: admin@blog.com
- **Password**: admin123

## Docker Compose Guide

### Introduction
Docker Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application's services, and then with a single command, you create and start all the services from your configuration.

### Prerequisites
- Docker installed on your system
- Command line/terminal access

### Installation Instructions

#### Windows
1. Docker Desktop for Windows includes Docker Compose
   - Download Docker Desktop from [Docker Hub](https://hub.docker.com/editions/community/docker-ce-desktop-windows/)
   - Follow the installation wizard
   - Docker Compose is automatically installed with Docker Desktop

#### Mac
1. Docker Desktop for Mac includes Docker Compose
   - Download Docker Desktop from [Docker Hub](https://hub.docker.com/editions/community/docker-ce-desktop-mac/)
   - Follow the installation instructions
   - Docker Compose is automatically installed with Docker Desktop

#### Linux
1. Download the Docker Compose binary:
   ```
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   ```

2. Apply executable permissions:
   ```
   sudo chmod +x /usr/local/bin/docker-compose
   ```

3. Verify the installation:
   ```
   docker-compose --version
   ```

### Common Commands

- Start containers:
  ```
  docker-compose up
  ```

- Start containers in detached mode:
  ```
  docker-compose up -d
  ```

- Stop containers:
  ```
  docker-compose down
  ```

- View logs:
  ```
  docker-compose logs
  ```

- Check status:
  ```
  docker-compose ps
  ```

## Troubleshooting

- If you encounter permission issues on Linux, make sure you've added your user to the docker group:
  ```
  sudo usermod -aG docker $USER
  ```

- For network issues, check if the ports are already in use:
  ```
  netstat -tuln
  ```

- To rebuild images after making changes:
  ```
  docker-compose up --build
  ```

## License

MIT
