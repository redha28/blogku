package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/internals/models"
	"github.com/redha28/blogku/internals/repositories"
	"github.com/redha28/blogku/internals/utils"
	"github.com/redha28/blogku/pkg"
	"github.com/redis/go-redis/v9"
)

// BlogController handles blog-related operations
type BlogController struct {
	repository repositories.BlogRepository
}

// NewBlogController creates a new blog controller
func NewBlogController(db *sql.DB, rdb *redis.Client) *BlogController {
	return &BlogController{
		repository: repositories.NewBlogRepository(db, rdb),
	}
}

// CreateBlog creates a new blog post with an image
// @Summary Create a new blog post
// @Description Create a new blog post with title, content, and image
// @Tags blogs
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Blog Title"
// @Param content formData string true "Blog Content"
// @Param image formData file true "Blog Image"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/blogs [post]
func (b *BlogController) CreateBlog(c *gin.Context) {
	// Parse multipart form with 10 MB max memory
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		pkg.Error("Failed to parse multipart form", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	var blogRequest models.BlogRequest
	if err := c.ShouldBind(&blogRequest); err != nil {
		pkg.Error("Invalid blog request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		pkg.Error("Missing image file", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	// Create blog post with image
	id, slug, err := b.repository.Create(blogRequest, file)
	if err != nil {
		pkg.Error("Failed to create blog post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post: " + err.Error()})
		return
	}
	pkg.GetLogger().InfoWithFields("Uploading image", map[string]any{
		"fileName": file.Filename,
		"fileSize": file.Size,
	})
	fileName, _, err := utils.NewUtils().FileHandling(c, file, slug, "")
	if err != nil {
		pkg.Error("Failed to upload image", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}
	log.Println("[INFO] Image uploaded successfully:", fileName)
	pkg.GetLogger().InfoWithFields("Blog post created", map[string]any{
		"id":    id,
		"slug":  slug,
		"title": blogRequest.Title,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Blog post created successfully",
		"blog": gin.H{
			"id":       id,
			"title":    blogRequest.Title,
			"content":  blogRequest.Content,
			"slug":     slug,
			"imageUrl": fileName,
		},
	})
}

// GetAllBlogs retrieves all blog posts with pagination
// @Summary Get all blog posts
// @Description Retrieve all blog posts with pagination
// @Tags blogs
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.BlogListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /blogs [get]
func (b *BlogController) GetAllBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := b.repository.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBlogBySlug retrieves a blog post by slug
// @Summary Get a blog post by slug
// @Description Retrieve a single blog post by its slug
// @Tags blogs
// @Produce json
// @Param slug path string true "Blog Slug"
// @Success 200 {object} models.BlogResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /blogs/{slug} [get]
func (b *BlogController) GetBlogBySlug(c *gin.Context) {
	slug := c.Param("slug")

	blog, err := b.repository.GetBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog post"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// UpdateBlog updates a blog post
func (b *BlogController) UpdateBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blogRequest models.BlogRequestUpdate
	if err := c.ShouldBindJSON(&blogRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if blogRequest.Title == "" && blogRequest.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (title or content) must be provided"})
		return
	}
	_, err = b.repository.Update(id, blogRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog post"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post updated successfully"})
}

// DeleteBlog deletes a blog post
func (b *BlogController) DeleteBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	_, err = b.repository.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog post"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}
