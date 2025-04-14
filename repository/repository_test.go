package repository_test

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// Global variables for test database
var (
	testDB *sqlx.DB
)

// TestMain runs once before all tests in the package
func TestMain(m *testing.M) {
	// Setup test database
	db, teardown := SetupTestDBForSuite()
	testDB = db

	// Run tests
	code := m.Run()

	// Cleanup
	if teardown != nil {
		teardown()
	}

	os.Exit(code)
}
