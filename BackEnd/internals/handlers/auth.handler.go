package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/internals/models"
	"github.com/redha28/blogku/internals/repositories"
	"github.com/redha28/blogku/pkg"
)

// AuthController handles authentication operations
type AuthController struct {
	repository repositories.AuthRepository
}

// NewAuthController creates a new auth controller
func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		repository: repositories.NewAuthRepository(db),
	}
}

// Login authenticates an admin user
// @Summary Admin login
// @Description Authenticate an admin user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.AdminLogin true "Login Request"
// @Success 200 {object} models.AdminResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var loginRequest models.AdminLogin

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add debugging information
	log.Printf("Login attempt with email: %s", loginRequest.Email)

	// Get admin using the versatile authentication method
	admin, hashedPassword, err := a.repository.GetAdminForAuth(loginRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with identifier: %s", loginRequest.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Printf("Database error during auth: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Safety check for nil admin
	if admin == nil {
		log.Println("Admin object is nil despite no database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "System error"})
		return
	}

	// Use Argon2 to compare passwords
	hasher := pkg.InitHashConfig()
	hasher.UseDefaultConfig()

	log.Printf("Comparing password for user: %s (ID: %d)", admin.Username, admin.ID)

	match, err := hasher.CompareHashAndPassword(hashedPassword, loginRequest.Password)
	// Fix: Check if err is nil before calling err.Error()
	if err != nil {
		log.Printf("Password comparison error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !match {
		log.Println("Password does not match")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	payload := pkg.NewPayload(strconv.Itoa(admin.ID), "admin")
	token, err := payload.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set cookie with JWT
	c.SetCookie(
		"authToken",
		token,
		int(time.Hour.Seconds()*24), // 24 hours
		"/",
		"",
		false, // set to true in production with HTTPS
		true,  // HTTP only
	)

	c.JSON(http.StatusOK, models.AdminResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
		Token:    token, // still include in response for API clients
	})
}

// Logout clears the auth cookie
// @Summary Admin logout
// @Description Clear the authentication cookie for the admin user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func (a *AuthController) Logout(c *gin.Context) {
	c.SetCookie(
		"authToken",
		"",
		-1, // expire immediately
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// CreateAdmin creates a new admin user (private endpoint)
// @Summary Create a new admin
// @Description Create a new admin user (requires API key)
// @Tags auth
// @Accept json
// @Produce json
// @Param adminRequest body models.AdminCreate true "Admin Create Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/admin/create [post]
func (a *AuthController) CreateAdmin(c *gin.Context) {
	var adminRequest models.AdminCreate

	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	exists, err := a.repository.CheckIfAdminExists(adminRequest.Username, adminRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash password with Argon2
	hasher := pkg.InitHashConfig()
	hasher.UseDefaultConfig()
	hashedPass, err := hasher.GenHashedPassword(adminRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new admin using repository
	id, err := a.repository.CreateAdmin(adminRequest, hashedPass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin created successfully",
		"admin": gin.H{
			"id":       id,
			"username": adminRequest.Username,
			"email":    adminRequest.Email,
		},
	})
}
