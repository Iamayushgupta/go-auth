package controller

import (
	"github.com/ayush/go-auth/model"
	"net/http"
	"github.com/ayush/go-auth/config"
	"log" 
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

//The context "c" provides access to information about the current HTTP request.
func SignUp(c *gin.Context) {
	var user model.User

	//The c.BindJSON -> Attempts to bind the JSON data from the request body to the user variable. 
	//The c.BindJSON(&user) function automatically parses the JSON from the request body and populates the user struct.
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON input for signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := user.SignUp(config.DB)
	if err != nil {
		// Sends back a JSON response
		c.JSON(500, gin.H{"error": err})
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}


//Basic-Auth Login
func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := user.Login(config.DB)

	if err != nil {
		if err.Error() == "user not found"{
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "invalid password"  {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	log.Printf("User %s signed up successfully", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}


var store = sessions.NewCookieStore([]byte("secret"))

//Session Login
func SessionLogin(c *gin.Context) {
	var credentials model.User
	if err := c.BindJSON(&credentials); err != nil {
		log.Printf("Failed to bind JSON input for login: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot bind JSON input"})
		return
	}

	err := credentials.Login(config.DB)

	if err != nil {
		if err.Error() == "user not found"{
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "invalid password"  {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	// store - An instance of a session store created using the Gorilla Sessions package. 
	//This store is responsible for managing session data and cookies.
	//Get(c.Request, "mysession"): This method tries to retrieve an existing s
	//ession by the name "mysession" from the user's browser cookies. 
	//If a session with this name does not exist, a new session is created.
	session, _ := store.Get(c.Request, "mysession")
	session.Values["username"] = credentials.Username

	// session.Save (c.Request, c.Writer): This method writes the session cookie to the user's browser. 
	//It takes the current HTTP request and response writer as arguments.
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	log.Printf("User %s signed up successfully", credentials.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

//Session Dashboard
func Dashboard(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	username := session.Values["username"]
	log.Printf("Dashboard Accessed")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the dashboard", "user": username})
}

//Session Logout
func SessionLogout(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	log.Printf("User %s loggest out successfully",session.Values["username"])
	session.Values["username"] = nil
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

//Session Auth
func AuthRequired(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")
	if session.Values["username"] == nil {
		log.Printf("Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}


func JwtSignUp(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := user.JwtSignUp(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sign up successful"})
}

func JwtLogin(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := user.JwtLogin()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SecureEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
}