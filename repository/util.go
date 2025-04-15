package repository

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to the database
func ConnectDB() (*sqlx.DB, error) {
	dbUser := GetEnvOrDefault("DB_USER", "admin")
	dbPassword := GetEnvOrDefault("DB_PASSWORD", "admin")
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/todoms?sslmode=disable", dbUser, dbPassword)
	return sqlx.Connect("postgres", dbURL)
}

// GetEnvOrDefault returns the value of an environment variable or the default if not set
func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
