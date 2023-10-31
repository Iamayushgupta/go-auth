package model

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/util"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) SignUp() (int,string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err) // Log the error
		return http.StatusInternalServerError, "Internal Server Error"
	}

	_, err = config.DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to sign up :  %v", err) // Log the error
		return http.StatusConflict,"Username already exist"
	}

	return http.StatusOK,""
}

func (u *User) Login() (int,string) {
	var storedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username=?", u.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		log.Printf(u.Username + " username does not exist")
		return http.StatusNotFound,"Username does not exist"
	} else if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError,"Internal Server Error"
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password)); err != nil {
		log.Printf("Wrong password error : %v",err)
		return http.StatusUnauthorized,"Password is incorrect"
	}
	
	return http.StatusAccepted,""
}

func (u *User) JwtLogin() (int,string, string) {
	var storedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username=?", u.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		log.Printf(u.Username + " username does not exist")
		return http.StatusNotFound,"Username does not exist",""
	} else if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError,"Internal Server Error",""
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password)); err != nil {
		log.Printf("Wrong password error : %v",err)
		return http.StatusUnauthorized,"Password is incorrect",""
	}
	
	token, err := util.GenerateToken(u.Username)
	if err != nil {
		return http.StatusInternalServerError, "Could not generate a Token",""
	}

	return http.StatusOK,"",token
}
