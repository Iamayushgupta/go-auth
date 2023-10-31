package controller

import (
	"log"
	"net/http"
	"os"
	"github.com/ayush/go-auth/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func SignUp(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode,err := user.SignUp()
	if err != "" {
		c.JSON(statusCode, gin.H{"error": err})
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "signup successful"})
}


//Basic-Auth Login
func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode,err := user.Login()

	if err !="" {
		c.JSON(statusCode, gin.H{"error":err})
		return
	}

	log.Printf("User %s logged up successfully", user.Username)
	c.JSON(statusCode, gin.H{"message": "Login successful"})
}


var store = sessions.NewCookieStore([]byte("secret"))

//Session Login
func SessionLogin(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for login: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot bind JSON input"})
		return
	}

	statusCode,err := user.Login()

	if err !="" {
		c.JSON(statusCode, gin.H{"error":err})
		return
	}
		
	// store - An instance of a session store created using the Gorilla Sessions package. 
	//This store is responsible for managing session data and cookies.
	//Get(c.Request, "mysession"): This method tries to retrieve an existing
	//session by the name "mysession" from the user's browser cookies. 
	//If a session with this name does not exist, a new session is created.
	session, _ := store.Get(c.Request, os.Getenv("SESSION_NAME"))
	session.Values["username"] = user.Username

	// session.Save (c.Request, c.Writer): This method writes the session cookie to the user's browser. 
	//It takes the current HTTP request and response writer as arguments.
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		// log.Printf("Failed to save the session : %v", err)
		return
	}

	log.Printf("User %s Logged In successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

//Session Dashboard
func Dashboard(c *gin.Context) {
	session, _ := store.Get(c.Request, os.Getenv("SESSION_NAME"))
	username := session.Values["username"]
	log.Printf("User %s accessed the dashboard",username)
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the dashboard", "user": username})
}

//Session Logout
func SessionLogout(c *gin.Context) {
	session, _ := store.Get(c.Request, os.Getenv("SESSION_NAME"))
	log.Printf("User %s logged out successfully",session.Values["username"])
	session.Values["username"] = nil
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

//Session Auth
func AuthRequired(c *gin.Context) {
	session, _ := store.Get(c.Request, os.Getenv("SESSION_NAME"))
	if session.Values["username"] == nil {
		log.Printf("Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}


//Login using JWT
func JwtLogin(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	statusCode,err, token := user.JwtLogin()
	if err != "" {
		c.JSON(statusCode, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SecureEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
}