package ws

import (
	"log"
	"time"

	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JoinGroupChat(c *gin.Context) {
	groupId := c.Query("groupId")
	userId := c.Query("userId")

	_, err := ValidGroupAndUser(groupId, userId)
	if err != nil {
		log.Println(err)
		return
	}
	if _, exists := RoomMap[groupId]; !exists {
		RoomMap[groupId] = &RoomManager{
			Clients: make(map[string]*websocket.Conn),
			Channel: make(chan models.RoomMessages),
		}
		go broadCastMessage(groupId)
	}
	existingConn, exists := RoomMap[groupId].Clients[userId]
	if exists {
		// 1. Prevent from making new connection
		// log.Println("Connection already exists")

		// 2. Delete existing connection
		log.Println("Existing connection removed")
		defer existingConn.Close()
		delete(RoomMap[groupId].Clients, userId)

		return

	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	RoomMap[groupId].Clients[userId] = conn
	newUserMessage := models.ChatMessage{
		SenderId: "system",
		Content:  "New user has entered the room",
		SendAt:   time.Now(),
	}
	// Notify the room that a new user has entered
	RoomMap[groupId].Channel <- models.RoomMessages{
		Message: []models.ChatMessage{newUserMessage},
	}

	handleWebsocketConnGroup(conn, groupId, userId)

}
func handleWebsocketConnGroup(conn *websocket.Conn, groupId, userID string) {
	defer func() {
		conn.Close()
		delete(RoomMap[groupId].Clients, userID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		newMessage := models.ChatMessage{
			SenderId: userID,
			Content:  string(message),
			SendAt:   time.Now(),
		}

		msg := models.RoomMessages{
			ID:      primitive.NewObjectID(),
			RoomId:  groupId,
			Message: []models.ChatMessage{newMessage},
		}
		SaveMsgToDb(newMessage, groupId)
		RoomMap[groupId].Channel <- msg
	}

}
