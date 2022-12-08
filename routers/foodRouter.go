package routers

import (
	"res/controller"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.Engine) {
	router.GET("/foods", controller.GetFoods())
	router.GET("/foods/:food_id", controller.GetFood())
	router.POST("/foods", controller.CreateFood())
	router.PATCH("/foods/:food_id", controller.UpdateFood())
}
