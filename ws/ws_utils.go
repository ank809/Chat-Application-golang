package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
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
		for participant, conn := range RoomMap[roomId].Clients {
			err := conn.WriteJSON(msg.Content)
			if err != nil {
				log.Println("Error broadcasting message to", participant, ":", err)
				conn.Close()
				delete(RoomMap[roomId].Clients, participant)
			}
		}
	}
}

func SaveMsgToDb(msg models.Messages, roomID string) {
	msgcoll := database.OpenCollection(database.Client, "Messages")
	_, err := msgcoll.InsertOne(context.Background(), msg)
	if err != nil {
		log.Println(err)
	}
	// Send the message to the room's broadcast channel
	roomCollection := database.OpenCollection(database.Client, "Rooms")
	update := bson.M{"$push": bson.M{"messages": msg.Message_id}}
	_, err = roomCollection.UpdateOne(context.TODO(), bson.M{"roomid": roomID}, update)
	if err != nil {
		log.Println("Error updating room's message list:", err)
	}
}
