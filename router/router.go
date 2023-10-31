package router

import (
	"github.com/ayush/go-auth/controller"
	"github.com/ayush/go-auth/middleware"
	"github.com/gin-gonic/gin"
)

// gin.Engine represents Gin Router
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/basic-auth/signup", controller.SignUp)
	r.POST("/basic-auth/login", controller.Login)

	r.POST("/session/login", controller.SessionLogin)
	r.GET("/session/dashboard", controller.AuthRequired, controller.Dashboard)
	r.GET("/session/logout", controller.AuthRequired,controller.SessionLogout)

	r.POST("/jwt/login", controller.JwtLogin)
	r.GET("/jwt/secure", middleware.Authenticate(), controller.SecureEndpoint)
	return r
}