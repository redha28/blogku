# Blog CMS API

A simple Blog CMS API built with Go, using MySQL for data storage and Redis for caching.

## Features

- Admin authentication using JWT
- Blog post management (CRUD operations)
- Automatic slug generation based on blog titles
- Redis caching for improved performance

## Setup Instructions

### Prerequisites

- Go 1.16+
- MySQL 5.7+
- Redis 6.0+

### Installation

1. Clone the repository
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

Or with hot reload:

```bash
fresh
```

## API Documentation

### Authentication

#### Login

- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
- **Response**:
  ```json
  {
    "id": 1,
    "username": "admin",
    "email": "admin@blog.com",
    "token": "your.jwt.token"
  }
  ```

### Blog Posts

#### Get All Blog Posts

- **URL**: `/api/v1/blogs`
- **Method**: `GET`
- **Query Parameters**:
  - `page` (default: 1)
  - `limit` (default: 10)
- **Response**:
  ```json
  {
    "total": 10,
    "blogs": [
      {
        "id": 1,
        "title": "Blog Title",
        "content": "Blog content...",
        "slug": "blog-title",
        "published_at": "2023-01-01T12:00:00Z"
      }
    ]
  }
  ```

#### Get Blog Post by Slug

- **URL**: `/api/v1/blogs/{slug}`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "id": 1,
    "title": "Blog Title",
    "content": "Blog content...",
    "slug": "blog-title",
    "published_at": "2023-01-01T12:00:00Z"
  }
  ```

#### Create Blog Post (Admin only)

- **URL**: `/api/v1/admin/blogs`
- **Method**: `POST`
- **Headers**:
  - `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "title": "New Blog Post",
    "content": "Content of the blog post..."
  }
  ```
- **Response**:
  ```json
  {
    "message": "Blog post created successfully",
    "blog": {
      "id": 2,
      "title": "New Blog Post",
      "content": "Content of the blog post...",
      "slug": "new-blog-post",
      "published_at": "2023-01-02T12:00:00Z"
    }
  }
  ```

#### Update Blog Post (Admin only)

- **URL**: `/api/v1/admin/blogs/{id}`
- **Method**: `PUT`
- **Headers**:
  - `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "title": "Updated Blog Post",
    "content": "Updated content..."
  }
  ```
- **Response**:
  ```json
  {
    "message": "Blog post updated successfully"
  }
  ```

#### Delete Blog Post (Admin only)

- **URL**: `/api/v1/admin/blogs/{id}`
- **Method**: `DELETE`
- **Headers**:
  - `Authorization: Bearer {token}`
- **Response**:
  ```json
  {
    "message": "Blog post deleted successfully"
  }
  ```

## Admin Credentials

- **Username**: admin
- **Password**: admin123

## License

MIT
