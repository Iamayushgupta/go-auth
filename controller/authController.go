package controller

import (
	"github.com/ayush/go-auth/model"
	"github.com/ayush/go-auth/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Making constants for error message
func SignUp(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode, err := userService.SignUp(&user)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "signup successful"})
}

// Basic-Auth Login
func Login(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode, err := userService.Login(&user)

	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User %s logged in successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "Login successful"})
}

// Login using JWT
func JwtLogin(c *gin.Context, userService *service.UserService) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode, token, err := userService.JwtLogin(&user)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Token Generated Successfully for user %s", user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SecureEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
}
