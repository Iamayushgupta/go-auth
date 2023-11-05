package router

import (
	"github.com/ayush/go-auth/controller"
	"github.com/ayush/go-auth/middleware"
	"github.com/gin-gonic/gin"
)

// gin.Engine represents Gin Router
func SetupRouter() *gin.Engine {
	r := gin.Default()
	// Making diagram, showing user, your service and database
	r.POST("/v1/signup", controller.SignUp)
	r.POST("/v1/login", controller.Login)

	// Use redis instead of local cache and read why you are using redis here and not mysql
	// Make a signup for this as well. As part of this, we have to create a cookie and return it in our signup function
	//Changing entire session logic to allow multiple users to login simultaneously
	r.POST("/v2/login", controller.SessionLogin)
	r.GET("/v2/dashboard", controller.AuthRequired, controller.Dashboard)
	r.GET("/v2/logout", controller.AuthRequired,controller.SessionLogout)

	r.POST("/v3/login", controller.JwtLogin)
	r.GET("/v3/secure", middleware.Authenticate(), controller.SecureEndpoint)
	return r
}