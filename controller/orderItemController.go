package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func IteamByOrder(id string) (orderItems []primitive.M, err error) {

}
func CreateOrderItem() gin.HandlerFunc {
	return func(*gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(*gin.Context) {

	}
}
