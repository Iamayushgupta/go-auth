package controller

import (
	"github.com/ayush/go-auth/model"
	"github.com/ayush/go-auth/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	invalidInputError      = "Invalid input"
	internalServerError    = "Internal Server Error"
	usernameExistsError    = "Username already exists"
	usernameNotFoundError  = "Username does not exist"
	incorrectPasswordError = "Password is incorrect"
	tokenGenerationError   = "Could not generate a Token"
)

func SignUp(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, invalidInputError)
		return
	}

	statusCode, err := userService.SignUp(&user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "signup successful"})
}

func Login(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, invalidInputError)
		return
	}

	statusCode, err := userService.Login(&user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	log.Printf("User %s logged in successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "Login successful"})
}

func JwtLogin(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, invalidInputError)
		return
	}

	statusCode, token, err := userService.JwtLogin(&user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	log.Printf("Token Generated Successfully for user %s", user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SecureEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
}

func handleError(c *gin.Context, statusCode int, message string) {
	log.Printf("Error: %s", message)
	c.JSON(statusCode, gin.H{"error": message})
}