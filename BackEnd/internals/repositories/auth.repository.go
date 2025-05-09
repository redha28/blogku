package repositories

import (
	"database/sql"

	"github.com/redha28/blogku/internals/models"
)

// AuthRepository handles database operations for authentication
type AuthRepository interface {
	GetAdminForAuth(identifier string) (*models.Admin, string, error)
	CheckIfAdminExists(username, email string) (bool, error)
	CreateAdmin(admin models.AdminCreate, hashedPassword string) (int64, error)
}

// SQLAuthRepository implements AuthRepository with MySQL
type SQLAuthRepository struct {
	DB *sql.DB
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &SQLAuthRepository{
		DB: db,
	}
}

// GetAdminForAuth retrieves an admin by email or username for authentication
func (r *SQLAuthRepository) GetAdminForAuth(identifier string) (*models.Admin, string, error) {
	var admin models.Admin
	var hashedPassword string

	// This query will match either email or username
	query := "SELECT id, username, password, email FROM admins WHERE email = ? OR username = ? LIMIT 1"
	err := r.DB.QueryRow(query, identifier, identifier).Scan(&admin.ID, &admin.Username, &hashedPassword, &admin.Email)
	if err != nil {
		return nil, "", err
	}

	return &admin, hashedPassword, nil
}

// CheckIfAdminExists checks if an admin exists by username or email
func (r *SQLAuthRepository) CheckIfAdminExists(username, email string) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM admins WHERE username = ? OR email = ?"
	err := r.DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateAdmin creates a new admin in the database
func (r *SQLAuthRepository) CreateAdmin(admin models.AdminCreate, hashedPassword string) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT INTO admins (username, password, email, created_at) VALUES (?, ?, ?, NOW())",
		admin.Username,
		hashedPassword,
		admin.Email,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
