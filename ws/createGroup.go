package ws

import (
	"context"
	"strings"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func CreateGroup(c *gin.Context) {
	var group models.Group
	users := c.Query("users")
	participants := strings.Split(users, ",")
	groupId := c.Query("groupid")
	if len(groupId) < 4 {
		c.JSON(400, "Length of group id should be greater than 4")
	}
	if len(participants) < 2 {
		c.JSON(400, gin.H{"error": "Atleast 2 users are required"})
		return
	}
	// check if group exists or not

	collection := database.OpenCollection(database.Client, "Groups")
	err := collection.FindOne(context.Background(), bson.M{"groupid": groupId}).Decode(&group)
	if err == nil {
		c.JSON(400, "Group already exists")
		return
	}
	group = models.Group{
		ID:           primitive.NewObjectID(),
		GroupId:      groupId,
		Participants: participants,
		Messages:     []string{},
	}

	_, err = collection.InsertOne(context.Background(), group)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, "Group successfully created")

}
