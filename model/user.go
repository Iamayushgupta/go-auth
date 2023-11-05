package model

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/util"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Moving this generate password to another class. User class does not need to know the exact implementation  *****
// That class should handle bcrypt and other relevant methods    *****
// Using zap for logging - Uber zap library

func (u *User) SignUp() (int, error) {
	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err) // Log the error
		return http.StatusInternalServerError, fmt.Errorf("internal Server Error")
	}

	// Moving these db functions to another class    *****
	// It will take username and password as input and execute this.   *****

	if util.UserExists(config.DB, u.Username) {
		log.Printf("Username already exists")
		return http.StatusConflict, fmt.Errorf("username already exists")
	}

	if err := util.InsertUser(config.DB, u.Username, hashedPassword); err != nil {
		log.Printf("Failed to sign up :  %v", err) // Log the error
		return http.StatusInternalServerError, err
	}

	// Return nil instead of empty string  *****
	return http.StatusOK, nil
}

// Making this more readable
func (u *User) Login() (int, error) {
	// Rename to fetched password or something  *****
	// Moving this method to some DB class   *****
	if !util.UserExists(config.DB, u.Username) {
		log.Printf("Username does not exists")
		return http.StatusNotFound, fmt.Errorf("username does not exists")
	}

	storedPassword, err := util.GetUserPassword(config.DB, u.Username)
	// Returning errors instead of strings  *****
	if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError, fmt.Errorf("internal Server Error")
	}

	// This method will go to another class   ******
	if err := util.ComparePasswords(storedPassword, u.Password); err != nil {
		log.Printf("Wrong password error : %v", err)
		return http.StatusUnauthorized, fmt.Errorf("password is incorrect")
	}

	return http.StatusAccepted, nil
}

// Fix the return fields, maybe token, class with fields err code and err
func (u *User) JwtLogin() (int, string, string) {
	storedPassword, err := util.GetUserPassword(config.DB, u.Username)
	// Abstracting this compare password logic
	if err == sql.ErrNoRows {
		log.Printf(u.Username + " username does not exist")
		return http.StatusNotFound, "Username does not exist", ""
	} else if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError, "Internal Server Error", ""
	}

	if err := util.ComparePasswords(storedPassword, u.Password); err != nil {
		log.Printf("Wrong password error : %v", err)
		return http.StatusUnauthorized, "Password is incorrect", ""
	}

	token, err := util.GenerateToken(u.Username)
	if err != nil {
		return http.StatusInternalServerError, "Could not generate a Token", ""
	}

	return http.StatusOK, "", token
}
