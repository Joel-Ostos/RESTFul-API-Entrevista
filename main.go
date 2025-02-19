package main

import (
	"log"
	_ "example.com/models"
	_ "example.com/helpers"
	"example.com/functions"
	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/users", functions.GetUsersHandler) 
  if err := r.Run(":8081"); err != nil {
    log.Fatal("Server failed to start: ", err)
  }
}
