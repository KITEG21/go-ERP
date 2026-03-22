package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	// Get the directory of the running executable
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	// Look for migrations in different locations
	migrationPaths := []string{
		filepath.Join(execDir, "db", "migrations"),
		filepath.Join(execDir, "..", "db", "migrations"),
		"db/migrations",
		"./db/migrations",
	}

	var m *migrate.Migrate
	for _, migrationPath := range migrationPaths {
		m, err = migrate.NewWithDatabaseInstance(
			"file://"+migrationPath,
			dsn,
			driver,
		)
		if err == nil {
			break
		}
	}

	if m == nil {
		// No migrations folder found, skip migrations
		return nil
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
