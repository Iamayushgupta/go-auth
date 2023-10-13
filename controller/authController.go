package controller

import (
	"github.com/ayush/go-auth/model"
	"net/http"
	"github.com/ayush/go-auth/config"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind JSON input"})
		return
	}

	err := user.SignUp(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign up"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind JSON input"})
		return
	}

	err := user.Login(config.DB)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
