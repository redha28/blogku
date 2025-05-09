package utils

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-friendly slug from a title
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-{2,}")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from beginning and end
	slug = strings.Trim(slug, "-")

	return slug
}

// EnsureUniqueSlug makes sure the slug is unique in the database
func EnsureUniqueSlug(db *sql.DB, slug string, excludeID int) (string, error) {
	baseSlug := slug
	counter := 1
	uniqueSlug := slug

	for {
		var exists bool
		var query string

		if excludeID > 0 {
			query = "SELECT EXISTS(SELECT 1 FROM blogs WHERE slug = ? AND id != ?)"
			err := db.QueryRow(query, uniqueSlug, excludeID).Scan(&exists)
			if err != nil {
				return "", err
			}
		} else {
			query = "SELECT EXISTS(SELECT 1 FROM blogs WHERE slug = ?)"
			err := db.QueryRow(query, uniqueSlug).Scan(&exists)
			if err != nil {
				return "", err
			}
		}

		if !exists {
			return uniqueSlug, nil
		}

		// If slug exists, add a counter suffix
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
}
