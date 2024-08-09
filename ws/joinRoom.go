package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/ank809/Chat-Application-golang/helpers"
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
	Channel chan models.Messages
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

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Initialize the RoomManager if it doesn't exist
	if _, exists := RoomMap[roomId]; !exists {
		RoomMap[roomId] = &RoomManager{
			Clients: make(map[string]*websocket.Conn),
			Channel: make(chan models.Messages),
		}
		go broadCastMessage(roomId)
	}

	RoomMap[roomId].Clients[userId] = conn

	// Notify the room that a new user has entered
	RoomMap[roomId].Channel <- models.Messages{
		Content:   "New user has entered the room",
		CreatedAt: time.Now(),
	}

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

		// Create a message object
		msg := models.Messages{
			ID:           primitive.NewObjectID(),
			Message_id:   helpers.GetUniqueKey(),
			Message_from: userID,
			CreatedAt:    time.Now(),
			Content:      string(message),
			RoomId:       roomID,
		}

		SaveMsgToDb(msg, roomID)
		RoomMap[roomID].Channel <- msg
	}

}
