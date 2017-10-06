package main

import (
	"github.com/gin-gonic/gin"
	line "github.com/pikomonde/fam100bot/src/line"
)

func main() {
	r := gin.New()
	r.GET("/ping", ping)

	// Line API
	r.POST("/line/webhook", line.Webhook)

	r.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
