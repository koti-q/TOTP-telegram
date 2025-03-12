package migrations

import (
	"database/sql"
	"log"
)

// Migrate function to create tables
func Migrate(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id BIGINT PRIMARY KEY UNIQUE
	)
	CREATE TABLE IF NOT EXISTS secrets (
		user_id BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		secret_id SERIAL PRIMARY KEY,
		secret_name TEXT UNIQUE,
		secret_value TEXT
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
