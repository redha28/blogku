package models

import (
	"mime/multipart"
	"time"
)

// Blog represents the blog post model
type Blog struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	Slug        string    `json:"slug"`
	ImagePath   string    `json:"image_path"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BlogRequest is used for creating/updating blog posts
type BlogRequest struct {
	Title   string                `form:"title" binding:"required"`
	Content string                `form:"content" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type BlogRequestUpdate struct {
	Title   string                `form:"title" binding:"omitempty"`
	Content string                `form:"content" binding:"omitempty"`
	Image   *multipart.FileHeader `form:"image" binding:"omitempty"`
}

// BlogResponse is used for API responses
type BlogResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Slug        string    `json:"slug"`
	ImagePath   string    `json:"image_path"`
	PublishedAt time.Time `json:"published_at"`
}

// BlogListResponse is used for paginated list responses
type BlogListResponse struct {
	Total int            `json:"total"`
	Blogs []BlogResponse `json:"blogs"`
	Meta  MetaPagination `json:"meta"`
}

type MetaPagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPage  int `json:"totalPage"`
	TotalItems int `json:"totalItems"`
}
