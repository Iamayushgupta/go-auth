package model

import (
	"database/sql"
	"fmt"
	"log" 
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) SignUp(DB *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err) // Log the error
		return fmt.Errorf("failed to hash password: %v", err)
	}

	_, err = DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to sign up, possibly duplicate username: %v", err) // Log the error
		return fmt.Errorf("failed to sign up, possibly duplicate username: %v", err)
	}

	log.Printf("User %s signed up successfully", u.Username) // Log success
	return nil
}

func (u *User) Login(DB *sql.DB) error {
	var storedPassword string
	row := DB.QueryRow("SELECT password FROM users WHERE username=?", u.Username)
	err := row.Scan(&storedPassword)
	if err != nil {
		log.Printf("Failed to retrieve user: %v", err) // Log the error
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password))
	if err != nil {
		log.Printf("Login failed for user %s: %v", u.Username, err) // Log the error
		return err
	}

	log.Printf("User %s logged in successfully", u.Username) // Log success
	return nil
}
