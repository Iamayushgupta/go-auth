package main

import (
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/router"
	"log"
)

func main() {
	config.ConnectToDB()
	config.InitializeRedisClient()
	//Defer ensures that the DB connection is closed when main.go exits
	defer config.DB.Close()

	//Setting up the router
	r := router.SetupRouter()
	log.Println("Server is starting and listening on port 8080")
	r.Run(":8080")
}
