package routers

import (
	"res/controller"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(router *gin.Engine) {
	router.GET("/orderItems", controller.GetOrderItems())
	router.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	router.GET("/orderItems-order/:order_id", controller.GetOrderItemByOrder())
	router.POST("/orderItems", controller.CreateOrderItem())
	router.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}
