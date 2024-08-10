package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	SenderId string    `json:"senderid"`
	Content  string    `json:"content"`
	SendAt   time.Time `json:"send_at"`
}

type RoomMessages struct {
	ID      primitive.ObjectID `bson:"_id"`
	RoomId  string             `json:"roomid"`
	Message []ChatMessage      `json:"message"`
}
