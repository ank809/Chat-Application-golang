package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Messages struct {
	ID           primitive.ObjectID  `bson:"_id"`
	Message_id   string              `json:"message_id"`
	Message_from string              `json:"message_from"`
	Message_to   string              `json:"message_to"`
	Content      string              `json:"content"`
	CreatedAt    primitive.Timestamp `json:"created_at"`
}
