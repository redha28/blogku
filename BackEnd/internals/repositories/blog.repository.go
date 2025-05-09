package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"os"
	"strings"
	"time"

	fp "path/filepath"

	"github.com/redha28/blogku/internals/models"
	"github.com/redha28/blogku/internals/utils"
	"github.com/redha28/blogku/pkg"
	"github.com/redis/go-redis/v9"
)

// BlogRepository handles database operations for blogs
type BlogRepository interface {
	Create(blog models.BlogRequest, file *multipart.FileHeader) (int64, string, error)
	GetAll(page, limit int) (models.BlogListResponse, error)
	GetBySlug(slug string) (models.BlogResponse, error)
	Update(id int, blog models.BlogRequestUpdate) (string, error)
	Delete(id int) (string, error)
}

// SQLBlogRepository implements BlogRepository with MySQL
type SQLBlogRepository struct {
	DB  *sql.DB
	RDB *redis.Client
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *sql.DB, rdb *redis.Client) BlogRepository {
	return &SQLBlogRepository{
		DB:  db,
		RDB: rdb,
	}
}

// Create adds a new blog post to the database
func (r *SQLBlogRepository) Create(blog models.BlogRequest, file *multipart.FileHeader) (int64, string, error) {
	// Generate slug from title
	ext := fp.Ext(file.Filename)
	allowedExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	if !allowedExt[ext] {
		return 0, "", fmt.Errorf("file extension not allowed")
	}
	slug := utils.GenerateSlug(blog.Title)

	pkg.Debug("Generated initial slug: " + slug)

	// Ensure slug is unique
	uniqueSlug, err := utils.EnsureUniqueSlug(r.DB, slug, 0)
	if err != nil {
		pkg.Error("Failed to ensure unique slug", err)
		return 0, "", err
	}

	pkg.Debug("Using unique slug: " + uniqueSlug)

	// Get current time
	now := time.Now()
	imagePath := fmt.Sprintf("%s_image%s", uniqueSlug, ext)
	// Insert blog post
	query := "INSERT INTO blogs (title, content, slug, image_path, published_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := r.DB.Exec(query, blog.Title, blog.Content, uniqueSlug, imagePath, now, now, now)
	if err != nil {
		pkg.Error("Failed to insert blog post", err)
		return 0, "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		pkg.Error("Failed to get last insert ID", err)
		return 0, "", err
	}

	// Clear cache for blog list
	ctx := context.Background()
	if err := r.RDB.Del(ctx, "blog:list").Err(); err != nil {
		pkg.Warn("Failed to clear blog list cache: " + err.Error())
	} else {
		pkg.Debug("Blog list cache cleared successfully")
	}

	pkg.GetLogger().InfoWithFields("Blog post created", map[string]interface{}{
		"id":   id,
		"slug": uniqueSlug,
	})

	return id, uniqueSlug, nil
}

// GetAll retrieves all blog posts with pagination
func (r *SQLBlogRepository) GetAll(page, limit int) (models.BlogListResponse, error) {
	ctx := context.Background()
	var response models.BlogListResponse

	cacheKey := fmt.Sprintf("blog:list:page:%d:limit:%d", page, limit)

	// Try to get from cache first
	cachedBlogs, err := r.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedBlogs), &response); err == nil {
			pkg.Debug("Returning blogs from cache")
			return response, nil
		}
	} else {
		pkg.Debug("Cache miss for blog list, fetching from database")
	}

	offset := (page - 1) * limit

	// Count total blogs
	var total int
	err = r.DB.QueryRow("SELECT COUNT(*) FROM blogs").Scan(&total)
	if err != nil {
		pkg.Error("Failed to count blogs", err)
		return response, err
	}

	// Get blogs with pagination (sorted by published_at DESC)
	rows, err := r.DB.Query(`
		SELECT id, title, content, slug, image_path, published_at
		FROM blogs
		ORDER BY published_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		pkg.Error("Failed to query blogs", err)
		return response, err
	}
	defer rows.Close()

	blogs := []models.BlogResponse{}
	for rows.Next() {
		var blog models.BlogResponse
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Slug, &blog.ImagePath, &blog.PublishedAt); err != nil {
			pkg.Error("Failed to scan blog row", err)
			return response, err
		}
		blogs = append(blogs, blog)
	}

	// Hitung total halaman
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response = models.BlogListResponse{
		Total: total,
		Blogs: blogs,
		Meta: models.MetaPagination{
			Page:       page,
			Limit:      limit,
			TotalPage:  totalPage,
			TotalItems: total,
		},
	}

	// Cache the result
	cacheData, _ := json.Marshal(response)
	if err := r.RDB.Set(ctx, cacheKey, cacheData, 10*time.Minute).Err(); err != nil {
		pkg.Warn("Failed to cache blog list: " + err.Error())
	} else {
		pkg.Debug("Blog list cached successfully")
	}

	pkg.GetLogger().InfoWithFields("Retrieved blog list", map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPage,
	})

	return response, nil
}

// GetBySlug retrieves a blog post by slug
func (r *SQLBlogRepository) GetBySlug(slug string) (models.BlogResponse, error) {
	ctx := context.Background()
	var blog models.BlogResponse

	// Try to get from cache first
	cacheKey := "blog:slug:" + slug
	cachedBlog, err := r.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedBlog), &blog); err == nil {
			return blog, nil
		}
	}

	// If not in cache, get from database
	query := "SELECT id, title, content, slug, image_path, published_at FROM blogs WHERE slug = ? LIMIT 1"
	err = r.DB.QueryRow(query, slug).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Slug, &blog.ImagePath, &blog.PublishedAt)
	if err != nil {
		return blog, err
	}

	// Cache the result
	cacheData, _ := json.Marshal(blog)
	r.RDB.Set(ctx, cacheKey, cacheData, 30*time.Minute)

	return blog, nil
}

// Update modifies an existing blog post
func (r *SQLBlogRepository) Update(id int, blog models.BlogRequestUpdate) (string, error) {
	// Ambil slug lama untuk invalidasi cache
	var existingSlug string
	err := r.DB.QueryRow("SELECT slug FROM blogs WHERE id = ?", id).Scan(&existingSlug)
	if err != nil {
		return "", err
	}

	// Buat potongan query dinamis
	fields := []string{}
	values := []any{}

	// Periksa apakah title diisi
	var newSlug string
	if blog.Title != "" {
		newSlug = utils.GenerateSlug(blog.Title)
		uniqueSlug, err := utils.EnsureUniqueSlug(r.DB, newSlug, id)
		if err != nil {
			return "", err
		}
		newSlug = uniqueSlug

		fields = append(fields, "title = ?", "slug = ?")
		values = append(values, blog.Title, newSlug)
	}

	// Periksa apakah content diisi
	if blog.Content != "" {
		fields = append(fields, "content = ?")
		values = append(values, blog.Content)
	}

	// Kalau tidak ada yang berubah
	if len(fields) == 0 {
		return existingSlug, nil
	}

	// Tambahkan updated_at
	fields = append(fields, "updated_at = ?")
	now := time.Now()
	values = append(values, now)

	// Tambahkan id ke parameter
	values = append(values, id)

	// Buat query update
	query := fmt.Sprintf("UPDATE blogs SET %s WHERE id = ?", strings.Join(fields, ", "))
	_, err = r.DB.Exec(query, values...)
	if err != nil {
		return "", err
	}

	// Clear cache lama
	ctx := context.Background()
	r.RDB.Del(ctx, "blog:list")
	r.RDB.Del(ctx, "blog:slug:"+existingSlug)

	return existingSlug, nil
}

// Delete removes a blog post
func (r *SQLBlogRepository) Delete(id int) (string, error) {
	// Get slug before deletion for cache invalidation
	var slug string
	var imagePath string
	err := r.DB.QueryRow("SELECT slug, image_path FROM blogs WHERE id = ?", id).Scan(&slug, &imagePath)
	if err != nil {
		return "", err
	}

	oldPath := fp.Join("public", "uploads", imagePath)
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		log.Println("[WARNING] Failed to delete old file:", err)
	}

	// Delete blog post
	_, err = r.DB.Exec("DELETE FROM blogs WHERE id = ?", id)
	if err != nil {
		return "", err
	}

	// Clear caches
	ctx := context.Background()
	r.RDB.Del(ctx, "blog:list")
	r.RDB.Del(ctx, "blog:slug:"+slug)

	return slug, nil
}
