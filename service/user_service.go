package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/model"
	"github.com/ayush/go-auth/util"
	"github.com/gin-gonic/gin"
)

const (
	invalidInputError      = "invalid input"
	internalServerError    = "internal Server Error"
	usernameExistsError    = "username already exists"
	usernameNotFoundError  = "username does not exist"
	incorrectPasswordError = "password is incorrect"
	tokenGenerationError   = "could not generate a Token"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) BindUserFromJSON(c *gin.Context) (*model.User, error) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to read input: %v", err)
		return nil, fmt.Errorf(invalidInputError)
	}
	return &user, nil
}

func (us *UserService) SignUp(u *model.User) (int, error) {
	// Implement the SignUp logic here
	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return http.StatusInternalServerError, fmt.Errorf(internalServerError)
	}

	if util.UserExists(config.DB, u.Username) {
		log.Printf("Username already exists")
		return http.StatusConflict, fmt.Errorf(usernameExistsError)
	}

	if err := util.InsertUser(config.DB, u.Username, hashedPassword); err != nil {
		log.Printf("Failed to sign up :  %v", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (us *UserService) Login(u *model.User) (int, error) {
	if !util.UserExists(config.DB, u.Username) {
		log.Printf("Username does not exists")
		return http.StatusNotFound, fmt.Errorf(usernameNotFoundError)
	}

	storedPassword, err := util.GetUserPassword(config.DB, u.Username)
	if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError, fmt.Errorf(internalServerError)
	}

	if err := util.ComparePasswords(storedPassword, u.Password); err != nil {
		log.Printf("Wrong password error : %v", err)
		return http.StatusUnauthorized, fmt.Errorf(incorrectPasswordError)
	}

	return http.StatusAccepted, nil
}

func (us *UserService) JwtLogin(u *model.User) (int, *model.Token, error) {
	if !util.UserExists(config.DB, u.Username) {
		log.Printf("Username does not exists")
		return http.StatusNotFound, &model.Token{TokenString: ""}, fmt.Errorf(usernameNotFoundError)
	}

	storedPassword, err := util.GetUserPassword(config.DB, u.Username)
	if err != nil {
		log.Printf("Database error during login for user %s: %v", u.Username, err)
		return http.StatusInternalServerError, &model.Token{TokenString: ""}, fmt.Errorf(internalServerError)
	}

	if err := util.ComparePasswords(storedPassword, u.Password); err != nil {
		log.Printf("Wrong password error : %v", err)
		return http.StatusUnauthorized, &model.Token{TokenString: ""}, fmt.Errorf(incorrectPasswordError)
	}

	token, err := util.GenerateToken(u.Username)
	if err != nil {
		return http.StatusInternalServerError, &model.Token{TokenString: ""}, fmt.Errorf(tokenGenerationError)
	}

	return http.StatusAccepted, token, nil
}
