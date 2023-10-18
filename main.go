package main

import (
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/router"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	config.ConnectToDB()
	defer config.DB.Close()

	r := router.SetupRouter()
	log.Println("Server is starting and listening on port 8080")
	r.Run(":8080")
}
