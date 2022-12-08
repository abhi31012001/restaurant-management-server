package routers

import (
	"res/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:user_id", controller.GetUser())
	router.POST("/user/signup", controller.SignUp())
	router.POST("/user/login", controller.Login())
}
