package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

func SaveSecret(userID int64, name, secret string) error {
	query := "INSERT INTO secrets (user_id, secret_name, secret_key) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, userID, name, string(secret))
	if err != nil {
		return err
	}
	return nil
}

func GetSecret(userID int64, name string) (string, error) {
	var secret string
	query := "SELECT secret_key FROM secrets WHERE user_id = $1 AND secret_name = $2"
	err := db.QueryRow(query, userID, name).Scan(&secret)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func GetUser(userID int64) (bool, error) {
	query := "SELECT user_id FROM USERS WHERE user_id = $1"
	err := db.QueryRow(query, userID).Scan(&userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddUser(userID int64) error {
	query := "INSERT INTO USERS (user_id) VALUES ($1)"
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSecret(userID int64, name string, secret string) error {
	query := "DELETE FROM secrets WHERE user_id = $1 AND secret_name = $2 AND secret_key = $3"
	_, err := db.Exec(query, userID, name, secret)
	if err != nil {
		return err
	}
	return nil
}
