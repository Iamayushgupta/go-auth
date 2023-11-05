package service

import (
	"fmt"
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/model"
	"github.com/ayush/go-auth/util"
	"log"
	"net/http"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) SignUp(u *model.User) (int, error) {
	// Implement the SignUp logic here
	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err) // Log the error
		return http.StatusInternalServerError, fmt.Errorf("internal Server Error")
	}

	if util.UserExists(config.DB, u.Username) {
		log.Printf("Username already exists")
		return http.StatusConflict, fmt.Errorf("username already exists")
	}

	if err := util.InsertUser(config.DB, u.Username, hashedPassword); err != nil {
		log.Printf("Failed to sign up :  %v", err) // Log the error
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (us *UserService) Login(u *model.User) (int, error) {
	if !util.UserExists(config.DB, u.Username) {
		log.Printf("Username does not exists")
		return http.StatusNotFound, fmt.Errorf("username does not exists")
	}

	storedPassword, err := util.GetUserPassword(config.DB, u.Username)
	if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError, fmt.Errorf("internal Server Error")
	}

	if err := util.ComparePasswords(storedPassword, u.Password); err != nil {
		log.Printf("Wrong password error : %v", err)
		return http.StatusUnauthorized, fmt.Errorf("password is incorrect")
	}

	return http.StatusAccepted, nil
}
