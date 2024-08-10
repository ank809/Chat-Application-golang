package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ValidateRoomAndUser(roomId, userId string) (models.Room, error) {
	var room models.Room
	coll := database.OpenCollection(database.Client, "Rooms")
	// room exists or not
	err := coll.FindOne(context.Background(), bson.M{"roomid": roomId}).Decode(&room)
	if err != nil {
		return models.Room{}, fmt.Errorf("room does not exists")
	}

	// user is  a participant of that room

	filter := bson.M{"roomid": roomId, "participants": bson.M{"$in": []string{userId}}}
	err = coll.FindOne(context.Background(), filter).Decode(&room)
	if err != nil {
		return models.Room{}, fmt.Errorf("you are not allowed to enter in the room")
	}
	return room, nil

}

func broadCastMessage(roomId string) {
	for msg := range RoomMap[roomId].Channel {
		// Broadcast the message to all participants in the room
		for _, message := range msg.Message {
			for participant, conn := range RoomMap[roomId].Clients {
				err := conn.WriteJSON(message)
				if err != nil {
					log.Println("Error broadcasting message to", participant, ":", err)
					conn.Close()
					delete(RoomMap[roomId].Clients, participant)
				}
			}
		}
	}
}

func SaveMsgToDb(msg models.ChatMessage, roomId string) {
	var msgc models.RoomMessages
	msgCollection := database.OpenCollection(database.Client, "Messages")

	// Find the document by roomId
	err := msgCollection.FindOne(context.Background(), bson.M{"roomid": roomId}).Decode(&msgc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no document exists, create a new document
			newMsg := models.RoomMessages{
				ID:      primitive.NewObjectID(),
				RoomId:  roomId,
				Message: []models.ChatMessage{msg},
			}

			_, insertErr := msgCollection.InsertOne(context.Background(), newMsg)
			if insertErr != nil {
				log.Println("Error creating new message document:", insertErr)
				return
			}
			log.Println("New message document created successfully.")
		} else {
			log.Println("Error finding message document:", err)
		}
		return
	}

	// If the document exists, append the new message to the Content slice
	update := bson.M{
		"$push": bson.M{"message": msg},
	}

	_, updateErr := msgCollection.UpdateOne(context.Background(), bson.M{"roomid": roomId}, update, options.Update().SetUpsert(true))
	if updateErr != nil {
		log.Println("Error updating message document:", updateErr)
		return
	}

	log.Println("Message appended to the existing document successfully.")
}

func ValidGroupAndUser(groupid, userId string) (models.Group, error) {
	var group models.Group
	collection := database.OpenCollection(database.Client, "Groups")
	err := collection.FindOne(context.Background(), bson.M{"groupid": groupid}).Decode(&group)
	if err != nil {
		return group, fmt.Errorf("group does not exists")
	}

	filter := bson.M{"groupid": groupid, "participants": bson.M{"$in": []string{userId}}}
	err = collection.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		return models.Group{}, fmt.Errorf("you are not allowed to enter in the room")
	}
	return group, nil
}
