package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	// Open MySQL connection
	var err error
	db, err = sql.Open("mysql", "root:ayushsql@tcp(127.0.0.1:3306)/authDB")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	r := gin.Default()
	r.POST("/signup", SignUp)
	r.POST("/login", Login)
	r.Run(":8080")
}

func SignUp(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "cannot bind JSON input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	_, err = db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", credentials.Username, string(hashedPassword))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to sign up, possibly duplicate username"})
		return
	}

	c.JSON(200, gin.H{"message": "signup successful"})
}

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "cannot bind JSON input"})
		return
	}

	var storedPassword string
	row := db.QueryRow("SELECT password FROM users WHERE username=?", credentials.Username)
	err := row.Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(401, gin.H{"error": "incorrect username or password"})
		} else {
			c.JSON(500, gin.H{"error": "database error"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "incorrect username or password"})
		return
	}

	c.JSON(200, gin.H{"message": "login successful"})
}
