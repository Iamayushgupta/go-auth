package controller

import (
	"github.com/ayush/go-auth/model"
	"net/http"
	"github.com/ayush/go-auth/config"
	"log" 
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

//Basic-Auth Signup
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


//Basic-Auth Login
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


var store = sessions.NewCookieStore([]byte("secret"))

//Session Login
func SessionLogin(c *gin.Context) {
	var credentials model.User
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot bind JSON input"})
		return
	}

	if err := credentials.Login(config.DB); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	session, _ := store.Get(c.Request, "mysession")
	session.Values["username"] = credentials.Username
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

//Session Dashboard
func Dashboard(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	username := session.Values["username"]
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the dashboard", "user": username})
}

//Session Logout
func SessionLogout(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	session.Values["username"] = nil
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

//Session Auth
func AuthRequired(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	if session.Values["username"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}