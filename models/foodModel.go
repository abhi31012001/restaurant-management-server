package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

type FOD struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" validate:"required,min=1,max=100"`
	Price      float64            `json:"price" validate:"required"`
	Food_image string             `json:"food_image" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"created_at"`
	Food_id    string
	Menu_id    string
}
