package controller

import (
	"github.com/ayush/go-auth/model"
	"net/http"
	"github.com/ayush/go-auth/config"
	"log" 
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind JSON input"})
		return
	}

	err := user.SignUp(config.DB)
	if err != nil {
		log.Printf("Failed to sign up user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign up"})
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for login: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind JSON input"})
		return
	}

	err := user.Login(config.DB)
	if err != nil {
		log.Printf("Login failed for user %s: %v", user.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect username or password"})
		return
	}

	log.Printf("User %s logged in successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}