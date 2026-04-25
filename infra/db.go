package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// InitDB creates an object of our DB
func InitDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// Schema creation
	schema := `
	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT
	);

	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		group_id INT REFERENCES groups(id),
		paid_by INT REFERENCES users(id),
		amount NUMERIC
	);

	CREATE TABLE IF NOT EXISTS participants (
		expense_id INT REFERENCES expenses(id),
		user_id INT REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS group_members (
		group_id INT REFERENCES groups(id),
		user_id INT REFERENCES users(id),
		PRIMARY KEY (group_id, user_id)
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return db
}
