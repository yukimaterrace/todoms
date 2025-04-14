package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// SetupTestDBForSuite sets up a test database for the entire test suite
func SetupTestDBForSuite() (*sqlx.DB, func()) {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("todoms"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	if err != nil {
		log.Fatalf("Failed to start test container: %v", err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	// Wait for the database to be ready
	db, err := waitForDB(dsn, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Apply migrations
	err = applyMigrations(dsn)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	teardown := func() {
		pgContainer.Terminate(ctx)
	}

	return db, teardown
}

func waitForDB(dsn string, timeout time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	start := time.Now()

	for time.Since(start) < timeout {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			return db, nil
		}
		time.Sleep(500 * time.Millisecond) // Wait before retrying
	}

	return nil, fmt.Errorf("database connection timeout: %w", err)
}

func applyMigrations(dsn string) error {
	// Get the absolute path to the migration directory
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filename)) // Move up two levels
	migrationPath := filepath.Join(basePath, "migration")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
