# BlogKu - Simple Blog CMS

This is a full-stack blog CMS application built with Next.js for the frontend and a separate backend API.

## Features

- Public blog listing and detail pages
- Admin authentication
- Admin dashboard for blog management
- CRUD operations for blog posts
- Responsive design

## Getting Started

### Frontend Setup (Next.js)

1. Install dependencies:

```bash
npm install
# or
yarn
# or
pnpm install
# or
bun install
```

2. Run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

### Backend API

The frontend interacts with a backend API at `http://localhost:8080` which provides the following endpoints:

- **Authentication**
  - `POST /api/v1/auth/login` - Login for admin user

- **Public Endpoints**
  - `GET /api/v1/blogs` - Get all published blogs
  - `GET /api/v1/blogs/:slug` - Get a specific blog by slug

- **Admin Endpoints** (require authentication)
  - `GET /api/v1/admin/blogs` - Get all blogs (admin view)
  - `POST /api/v1/admin/blogs` - Create a new blog
  - `GET /api/v1/admin/blogs/:id` - Get a specific blog by ID
  - `PUT /api/v1/admin/blogs/:id` - Update a specific blog
  - `DELETE /api/v1/admin/blogs/:id` - Delete a specific blog

## Admin Credentials

For admin access, use the credentials provided by the backend service.

## Tech Stack

- **Frontend**:
  - Next.js 15
  - React 19
  - TypeScript
  - Tailwind CSS for styling

- **Backend**:
  - REST API provided separately
