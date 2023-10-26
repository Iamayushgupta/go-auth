package model

import (
	"database/sql"
	"errors"
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
		return err
	}

	_, err = DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to sign up :  %v", err) // Log the error
		return err
	}

	return nil
}

func (u *User) Login(DB *sql.DB) error {
	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE username=?", u.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		log.Printf(u.Username + " username does not exist")
		return errors.New("user not found")
	} else if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return errors.New("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password)); err != nil {
		log.Printf("Wrong password error : %v",err)
		return errors.New("invalid password")
	}
	
	return nil
}