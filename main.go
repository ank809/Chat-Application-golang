package main

import (
	"log"
	"net/http"

	"github.com/ank809/Chat-Application-golang/authentication"
	"github.com/ank809/Chat-Application-golang/middlewares"
	"github.com/ank809/Chat-Application-golang/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup", authentication.Signup)
	r.GET("/login", authentication.Login)
	authRoutes := r.Group("/", middlewares.AuthMiddleware())
	{
		authRoutes.POST("/createRoom", ws.CreateRoom)
		authRoutes.POST("/createGroup", ws.CreateGroup)
		authRoutes.POST("/adduser", ws.AddUserToGroup)
		authRoutes.GET("/deleteuser", ws.DeleteUserFromGroup)

	}
	r.GET("/joingroupchat", ws.JoinGroupChat)
	r.GET("/joinRoom", ws.JoinRoom)
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println(err)
		return
	}
}
