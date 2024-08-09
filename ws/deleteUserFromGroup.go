package ws

import (
	"context"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func DeleteUserFromGroup(c *gin.Context) {
	var room models.Room
	groupId := c.Query("groupId")
	userId := c.Query("userId")

	coll := database.OpenCollection(database.Client, "Groups")
	err := coll.FindOne(context.TODO(), bson.M{"groupid": groupId}).Decode(&room)

	if err != nil {
		c.JSON(400, "Group does not found")
		return
	}

	// check if user exists in group
	update := bson.M{"$pull": bson.M{"participants": userId}}
	res, err := coll.UpdateOne(context.Background(), bson.M{"groupid": groupId}, update)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	if res.ModifiedCount == 0 {
		c.JSON(400, "User not found")
		return
	}
	c.JSON(200, "User deleted successfully")

}
