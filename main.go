package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pikomonde/fam100bot/src/line"
)

func main() {
	if os.Getenv("ENV") == "dev" {
		fmt.Println("==== FAM100BOT ====")
		for {
			var uID, msg string
			fmt.Scan(&uID, &msg)
			line.SimulateWebhook(uID, msg)
		}
	}

	r := gin.New()
	r.GET("/ping", ping)

	// Webhook Endpoint
	r.POST("/line/webhook", line.Webhook)

	r.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
