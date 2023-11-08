package controller

import (
	"log"
	"net/http"

	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/service"
	"github.com/ayush/go-auth/util"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, statusCode int, message string) {
	log.Printf("Error: %s", message)
	c.JSON(statusCode, gin.H{"error": message})
}

func SignUp(c *gin.Context, userService *service.UserService) {
	user, err := userService.BindUserFromJSON(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := userService.SignUp(user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "signup successful"})
}

func Login(c *gin.Context, userService *service.UserService) {
	user, err := userService.BindUserFromJSON(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := userService.Login(user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	log.Printf("User %s logged in successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "Login successful"})
}

//Session-Login

func SessionLogin(c *gin.Context, userService *service.UserService) {
	user, err := userService.BindUserFromJSON(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := userService.Login(user)
	if err != nil {
		handleError(c, statusCode, err.Error())
		return
	}

	// First you check if there is a session existing for this user or not
	sessionToken, err := util.GenerateRandomString(32)
	if err != nil {
		handleError(c, 500, err.Error())
	}

	err = util.SetSession(config.RedisClient, user.Username, sessionToken)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create a session")
		return
	}

	// Set the session token in a cookie
	// You can let the user be logged in for 1 hour
	// Fix local
	c.SetCookie("session", sessionToken, 3600, "/", "localhost", false, true)

	c.JSON(statusCode, gin.H{"message": "Login successful"})
}

func SessionLogout(c *gin.Context) {
	// Retrieve the session token from the cookie
	sessionToken, err := c.Cookie("session")
	if err != nil {
		handleError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Delete the session from Redis
	err = util.DeleteSession(config.RedisClient, sessionToken)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete session")
		return
	}

	// Clear the session cookie
	c.SetCookie("session", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Session Dashboard
func Dashboard(c *gin.Context) {
	// Get the session token from the cookie
	sessionToken, err := c.Cookie("session")
	if err != nil {
		// No session token found, consider the user not logged in
		handleError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if the session token is valid
	username, err := util.GetUsernameFromSession(config.RedisClient, sessionToken)
	if err != nil {
		// Handle error, e.g., session not found or invalid session
		handleError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// User is authenticated; you can provide the protected content here
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the dashboard, " + username})
}

func JwtLogin(c *gin.Context, userService *service.UserService) {
	user, err := userService.BindUserFromJSON(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, token, err := userService.JwtLogin(user)
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
