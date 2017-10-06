package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func Webhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

	}
	decoded, err := base64.StdEncoding.DecodeString(c.Request.Header.Get("X-Line-Signature"))
	if err != nil {

	}
	hash := hmac.New(sha256.New, []byte("<channel secret>"))
	hash.Write(body)

	hmac.Equal(decoded, hash.Sum(nil))

	c.JSON(200, gin.H{
		"message": "pong",
	})
	fmt.Println("===> A")
	log.Println("===> B")
}
