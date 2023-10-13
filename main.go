package main

import (
	"github.com/ayush/go-auth/config"
	"github.com/ayush/go-auth/router"
)

func main() {
	config.ConnectToDB()
	defer config.DB.Close()

	r := router.SetupRouter()
	r.Run(":8080")
}
