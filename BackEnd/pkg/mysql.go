package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

var DB *sql.DB

func Connect() (*sql.DB, error) {
	// Get the environment variables for the database connection
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")

	// Validate and use defaults if necessary
	if dbHost == "" {
		dbHost = "localhost"
		log.Println("DBHOST not set, using default:", dbHost)
	}

	if dbPort == "" {
		dbPort = "3306" // Default MySQL port
		log.Println("DBPORT not set, using default:", dbPort)
	}

	if dbName == "" {
		return nil, fmt.Errorf("DBNAME environment variable is required")
	}

	// Create the connection string for MySQL
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Print connection info for debugging (hide password)
	connectionInfo := fmt.Sprintf("Connecting to MySQL: User: %s, Host: %s, Port: %s, Database: %s",
		dbUser, dbHost, dbPort, dbName)
	log.Println(connectionInfo)

	// Open the connection to the MySQL database
	var err error
	DB, err = sql.Open("mysql", dbString)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Check if the database is reachable
	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	log.Println("MySQL DB connected successfully")
	return DB, nil
}
