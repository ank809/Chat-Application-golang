package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Messages struct {
	ID           primitive.ObjectID `bson:"_id"`
	RoomId       string             `json:"roomid"`
	Message_id   string             `json:"message_id"`
	Message_from string             `json:"message_from"`
	Content      string             `json:"content"`
	CreatedAt    time.Time          `json:"created_at"`
}
