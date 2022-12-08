package routers

import (
	"res/controller"

	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.Engine) {
	router.GET("/tables", controller.GetTables())
	router.GET("/tables/:table_id", controller.GetTable())
	router.POST("/tables", controller.CreateTable())
	router.PATCH("/tables/:table_id", controller.UpdateTable())
}
