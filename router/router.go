package router

import (
	"github.com/ayush/go-auth/controller"
	"github.com/ayush/go-auth/middleware"
	"github.com/ayush/go-auth/service"
	"github.com/gin-gonic/gin"
)

// gin.Engine represents Gin Router
func SetupRouter() *gin.Engine {
	r := gin.Default()
	userService := service.NewUserService()
	// Making diagram, showing user, your service and database
	r.POST("/v1/signup", func(c *gin.Context) {
		controller.SignUp(c, userService) // Pass the userService to the SignUp handler
	})
	r.POST("/v1/login", func(c *gin.Context) {
		controller.Login(c, userService) // Pass the userService to the SignUp handler
	})

	// Use redis instead of local cache and read why you are using redis here and not mysql
	// Make a signup for this as well. As part of this, we have to create a cookie and return it in our signup function
	//Changing entire session logic to allow multiple users to login simultaneously
	// r.POST("/v2/login", controller.SessionLogin)
	r.GET("/v2/dashboard", controller.AuthRequired, controller.Dashboard)
	r.GET("/v2/logout", controller.AuthRequired, controller.SessionLogout)

	r.POST("/v3/login", controller.JwtLogin)
	r.GET("/v3/secure", middleware.Authenticate(), controller.SecureEndpoint)
	return r
}
