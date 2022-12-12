package controller

import (
	"context"
	"log"
	"net/http"
	"res/database"
	"res/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := tableCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allTable []bson.M
		if err = result.All(ctx, &allTable); err != nil {
			log.Fatal(err.Error())
			return
		}
		c.JSON(http.StatusOK, allTable)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var tableId = c.Param("table_id")
		if err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, internalErr := tableCollection.InsertOne(ctx, table)
		if internalErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": internalErr})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var tableId = c.Param("table_id")
		var table models.Table
		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateObj primitive.D
		filter := bson.E{"table_id", tableId}
		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
		}
		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", table.Updated_at})

		upsert := true
		obj := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &obj)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
