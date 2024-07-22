package ws

import (
	"context"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/helpers"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func CreateRoom(c *gin.Context) {

	var roomparticipant models.RoomParticipant
	var existingRoom models.Room

	err := c.BindJSON(&roomparticipant)
	if err != nil {
		c.JSON(400, "Error in binding json")
		return
	}
	coll := database.OpenCollection(database.Client, "Rooms")
	filter := bson.M{
		"participants": bson.M{"$all": []string{roomparticipant.User1, roomparticipant.User2}},
	}
	err = coll.FindOne(context.TODO(), filter).Decode(&existingRoom)
	if err == nil {
		c.JSON(400, "Room with these participant already exists")
		return
	}

	var room models.Room = models.Room{
		ID:           primitive.NewObjectID(),
		RoomID:       helpers.GenerateRoomId(),
		Participants: []string{roomparticipant.User1, roomparticipant.User2},
		Messages:     []string{},
	}

	_, err = coll.InsertOne(context.TODO(), room)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"success": "Room created successuflly",
		"User1":   roomparticipant.User1,
		"User2":   roomparticipant.User2,
	})
}
