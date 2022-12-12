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

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Intenal error in Order item "})
			return
		}
		var allMenu []bson.M

		if err = result.All(ctx, &allMenu); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allMenu)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No order item has been found"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")

		allOrdersItem, err := orderItemCollection.Find(ctx, bson.M{"order_id": orderId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Intenal error in Order item "})
			return
		}
		var allMenu []bson.M

		if err = allOrdersItem.All(ctx, &allMenu); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allMenu)
	}
}
func IteamByOrder(id string) (orderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookUpStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullandEmptyArray", true}}}}
	lookUpDoorStage := bson.D{{"$lookup", bson.D{{"food", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}

	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullandEmptyArray", true}}}}
	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "tab;e_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullandEmptyArray", true}}}}
	projectStage := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"amount", "$food_price"},
				{"total_count", 1},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			},
		},
	}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"order_item", bson.D{{"$push", "$$ROOT"}}}}}}
	projectStage2 := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"payment_due", 1},
				{"total_count", 1},
				{"table_number", "$_id.table_number"},
				{"order_items", 1},
			},
		},
	}
	result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookUpStage,
		unwindStage,
		lookUpDoorStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		panic(err)
	}
	if err = result.All(ctx, &orderItems); err != nil {
		panic(err)
	}
	defer cancel()
	return orderItems, err

}
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem OrderItemPack
		var order models.Order
		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not agble to bind Json in orderItem"})
			return
		}
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsTobeInserted := []interface{}{}
		order.Table_id = orderItem.Table_id
		order_id := OrderItemOrderCreator(order)
		for _, Orderitem := range orderItem.Order_items {
			Orderitem.Order_id = order_id
			validateErr := validate.Struct(orderItem)
			if validateErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validateErr})
				return
			}
			Orderitem.ID = primitive.NewObjectID()
			Orderitem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			Orderitem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			Orderitem.Order_item_id = Orderitem.ID.Hex()
			var num = toFixed(*Orderitem.Unit_price, 2)
			Orderitem.Unit_price = &num
			orderItemsTobeInserted = append(orderItemsTobeInserted, Orderitem)

		}
		inserted, err := orderItemCollection.InsertMany(ctx, orderItemsTobeInserted)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, inserted)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")
		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"order_item_id": orderItemId}
		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.Unit_price})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})

		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItem.Food_id})

		}
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OrderaItems is not Updated"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}

}
