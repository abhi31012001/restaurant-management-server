package main

import (
	"log"
	"os"
	"res/database"
	"res/middleware"
	"res/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())
	routers.UserRoutes(router)
	router.Use(middleware.Authentication())
	routers.FoodRoutes(router)
	routers.MenuRoutes(router)
	routers.TableRoutes(router)
	routers.OrderRoutes(router)
	routers.OrderItemRoutes(router)
	routers.InvoiceRoutes(router)
	router.Run(":" + port)

}
