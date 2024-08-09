package ws

import (
	"context"
	"strings"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func AddUserToGroup(c *gin.Context) {
	var group models.Group
	groupId := c.Query("groupid")
	users := c.Query("users")
	participant := strings.Split(users, ",")

	coll := database.OpenCollection(database.Client, "Groups")
	err := coll.FindOne(context.Background(), bson.M{"groupid": groupId}).Decode(&group)
	if err != nil {
		c.JSON(400, "Group does not exists")
		return
	}
	update := bson.M{"$push": bson.M{"participants": bson.M{"$each": participant}}}
	_, err = coll.UpdateOne(context.Background(), bson.M{"groupid": groupId}, update)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, "User successfully added to group")
}
