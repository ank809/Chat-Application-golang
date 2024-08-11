package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID           primitive.ObjectID `bson:"_id"`
	RoomID       string             `json:"roomID"  binding:"required"`
	Participants []string           `json:"participants"`
}

type RoomParticipant struct {
	User1 string `json:"user1" binding:"required"`
	User2 string `json:"user2" binding:"required"`
}

type Group struct {
	ID           primitive.ObjectID `bson:"_id"`
	GroupId      string             `json:"groupid"`
	Participants []string           `json:"participants"`
	Messages     []string           `json:"messages"`
}
