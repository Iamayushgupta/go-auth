// util/db.go
package util

import (
	"database/sql"
	"fmt"
)

// InsertUser inserts a new user into the database.
func InsertUser(db *sql.DB, username, password string) error {
	_, err := db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

// GetUserPassword retrieves the password for a user from the database.
func GetUserPassword(db *sql.DB, username string) (string, error) {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&storedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to get user password: %v", err)
	}
	return storedPassword, nil
}

func UserExists(db *sql.DB, username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username).Scan(&count)
	if err != nil {
		fmt.Printf("Error checking user existence: %v\n", err)
		return false
	}
	return count > 0
}
