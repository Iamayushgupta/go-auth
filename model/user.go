package model

import (
	"database/sql"
	"errors"
	"log"
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/ayush/go-auth/config"
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

func (u *User) JwtSignUp() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	_, err = config.DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to sign up: %v", err)
	}

	return nil
}

func (u *User) JwtLogin() (string, error) {
	var user User
	row := config.DB.QueryRow("SELECT username, password FROM users WHERE username=?", u.Username)
	err := row.Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("query error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := createToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("could not create token: %v", err)
	}

	return token, nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %v", err)
	}
	return tokenString, nil
}