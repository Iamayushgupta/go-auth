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
	userGroup1 := r.Group("/v1")
	{
		userGroup1.POST("/signup", func(c *gin.Context) {
			controller.SignUp(c, userService)
		})
		userGroup1.POST("/login", func(c *gin.Context) {
			controller.Login(c, userService)
		})
	}
	userGroup2 := r.Group("/v2")
	{
		userGroup2.POST("/login", func(c *gin.Context) {
			controller.SessionLogin(c, userService)
		})
		userGroup2.GET("/dashboard", controller.Dashboard)
		userGroup2.GET("/logout", controller.SessionLogout)
	}

	userGroup3 := r.Group("/v3")
	{
		userGroup3.POST("/login", func(c *gin.Context) {
			controller.JwtLogin(c, userService)
		})
		userGroup3.GET("/dashboard", middleware.Authenticate(), controller.SecureEndpoint)
	}

	return r
}
