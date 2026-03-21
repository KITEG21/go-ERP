package database

import (
	"fmt"
	"os"
)

func RunMigrations() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// AutoMigrate models using existing DB connection
	// This creates tables if they don't exist
	err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255),
			password VARCHAR(255) NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS departments (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT
		);
		
		CREATE TABLE IF NOT EXISTS workers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			phone VARCHAR(50),
			department_id INTEGER REFERENCES departments(id),
			salary DECIMAL(10,2),
			hire_date DATE
		);
		
		CREATE TABLE IF NOT EXISTS attendances (
			id SERIAL PRIMARY KEY,
			worker_id INTEGER REFERENCES workers(id),
			check_in TIME,
			check_out TIME,
			date DATE,
			status INTEGER DEFAULT 0
		);
		
		CREATE TABLE IF NOT EXISTS payrolls (
			id SERIAL PRIMARY KEY,
			worker_id INTEGER REFERENCES workers(id),
			month VARCHAR(7),
			base_salary DECIMAL(10,2),
			bonus DECIMAL(10,2),
			deductions DECIMAL(10,2),
			net_salary DECIMAL(10,2),
			status INTEGER DEFAULT 0
		);
	`).Error

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
