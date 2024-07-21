package main

import (
	"log"
	"net/http"

	"github.com/ank809/Chat-Application-golang/authentication"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup", authentication.Signup)
	r.GET("/login", authentication.Login)

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println(err)
		return
	}
}
