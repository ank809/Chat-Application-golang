package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/ank809/Chat-Application-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RoomManager struct {
	Clients map[string]*websocket.Conn // Map to store WebSocket connections
	Channel chan models.RoomMessages
}

var RoomMap = make(map[string]*RoomManager) // Map to manage rooms by room ID

func JoinRoom(c *gin.Context) {
	roomId := c.Query("roomId")
	userId := c.Query("userId")

	_, err := ValidateRoomAndUser(roomId, userId)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Initialize the RoomManager if it doesn't exist
	if _, exists := RoomMap[roomId]; !exists {
		RoomMap[roomId] = &RoomManager{
			Clients: make(map[string]*websocket.Conn),
			Channel: make(chan models.RoomMessages),
		}
		go broadCastMessage(roomId)
	}
	// Check if the user already have a connection or not

	existingConn, exists := RoomMap[roomId].Clients[userId]
	if exists {
		// 1. Prevent from making new connection
		// log.Println("Connection already exists")

		// 2. Delete existing connection
		log.Println("Existing connection removed")
		existingConn.Close()
		delete(RoomMap[roomId].Clients, userId)
	}

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	RoomMap[roomId].Clients[userId] = conn

	newMessage := models.ChatMessage{
		SenderId: "system",
		Content:  "New user has entered in the chat",
		SendAt:   time.Now(),
	}
	msg := models.RoomMessages{
		ID:      primitive.NewObjectID(),
		RoomId:  roomId,
		Message: []models.ChatMessage{newMessage},
	}
	// Notify the room that a new user has entered
	RoomMap[roomId].Channel <- msg

	handleWebsocketConn(conn, roomId, userId)
}

func handleWebsocketConn(conn *websocket.Conn, roomID, userID string) {
	defer func() {
		conn.Close()
		delete(RoomMap[roomID].Clients, userID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		newMessage := models.ChatMessage{
			Content:  string(message),
			SendAt:   time.Now(),
			SenderId: userID,
		}
		// // Create a message object
		msg := models.RoomMessages{
			ID:      primitive.NewObjectID(),
			RoomId:  roomID,
			Message: []models.ChatMessage{newMessage},
		}

		SaveMsgToDb(newMessage, roomID)
		RoomMap[roomID].Channel <- msg
	}

}
