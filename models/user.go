package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Name        string             `json:"name" binding:"required"`
	About       string             `json:"about"`
	Email       string             `json:"email" binding:"required"`
	PhoneNumber string             `json:"phoneNumber" binding:"required" `
	Profile_Url string             `json:"profile_url"`
	Password    string             `json:"password" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
