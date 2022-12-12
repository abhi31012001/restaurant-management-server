package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID            primitive.ObjectID `bson:"_id"`
	Quantity      *string            `json:"quantity"`
	Unit_price    *float64           `json:"unit_price"`
	Food_id       *string            `json:"food_id"`
	Order_id      string             `json:"order_id"`
	Order_item_id string             `json:"order_item_id"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
