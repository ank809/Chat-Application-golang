package main

import (
	"log"
	"net/http"

	"github.com/ank809/Chat-Application-golang/authentication"
	"github.com/ank809/Chat-Application-golang/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup", authentication.Signup)
	r.GET("/login", authentication.Login)
	r.POST("/createRoom", ws.CreateRoom)

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println(err)
		return
	}
}
