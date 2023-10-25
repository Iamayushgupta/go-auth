package main

import (
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/router"
	"log"
)

func main() {
	config.ConnectToDB()
	defer config.DB.Close()

	r := router.SetupRouter()
	log.Println("Server is starting and listening on port 8080")
	r.Run(":8080")
}
