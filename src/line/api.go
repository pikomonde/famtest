package line

import "github.com/gin-gonic/gin"

func Webhook(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
