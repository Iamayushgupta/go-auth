package router

import (
	"github.com/ayush/go-auth/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/signup", controller.SignUp)
	r.POST("/login", controller.Login)
	return r
}
