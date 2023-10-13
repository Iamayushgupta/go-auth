package model

import (
	"database/sql"
	"fmt"
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
		return fmt.Errorf("failed to hash password: %v", err)
	}

	_, err = DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to sign up, possibly duplicate username: %v", err)
	}

	return nil
}

func (u *User) Login(DB *sql.DB) error {
	var storedPassword string
	row := DB.QueryRow("SELECT password FROM users WHERE username=?", u.Username)
	err := row.Scan(&storedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password))
	if err != nil {
		return err
	}

	return nil
}
