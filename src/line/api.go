package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type WebhookEvents struct {
	Events []WebhookEvent `json:"events"`
}
type WebhookEvent struct {
	ReplyToken string          `json:"replyToken"`
	Type       string          `json:"type"`
	Timestamp  int64           `json:"timestamp"`
	Source     WebhookSource   `json:"source"`
	Message    WebhookMessage  `json:"message"`
	Postback   WebhookPostback `json:"postback"`
	Beacon     WebhookBeacon   `json:"beacon"`
}
type WebhookSource struct {
	Type    string `json:"type"`
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
	RoomID  string `json:"roomId"`
}
type WebhookMessage struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Text      string  `json:"text"`      // type: text
	Filename  string  `json:"fileName"`  // type: file
	FileSize  string  `json:"fileSize"`  // type: file
	Title     string  `json:"title"`     // type: location
	Address   string  `json:"address"`   // type: location
	Latitude  float64 `json:"latitude"`  // type: location
	Longitude float64 `json:"longitude"` // type: location
	PackageID float64 `json:"packageId"` // type: sticker
	StickerID float64 `json:"stickerId"` // type: sticker
	GroupID   float64 `json:"groupId"`   // type: join/leave
}
type WebhookPostback struct {
	Data   string      `json:"data"`
	Params interface{} `json:"params"`
}
type WebhookBeacon struct {
	HwID string `json:"hwid"`
	Type string `json:"type"`
	DM   string `json:"dm"`
}

func Webhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(c.Request.Header.Get("X-Line-Signature"))
	if err != nil {
		log.Println(err)
		return
	}
	hash := hmac.New(sha256.New, []byte(os.Getenv("CHANNELSECRET")))
	hash.Write(body)

	if !hmac.Equal(decoded, hash.Sum(nil)) {
		log.Println("HMAC not equal")
		return
	}

	fmt.Println("===> A")
	var webhookObj WebhookEvents
	err = json.Unmarshal(body, &webhookObj)
	if err != nil {
		log.Println(err)
		return
	}

	// Reply Message
	bot, err := linebot.New(os.Getenv("CHANNELSECRET"), os.Getenv("CHANNELTOKEN"))
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := bot.ReplyMessage(webhookObj.Events[0].ReplyToken, linebot.NewTextMessage(webhookObj.Events[0].Message.Text)).Do(); err != nil {
		log.Println(err)
		return
	}

	// Push Message
	for i := 0; i < 250; i++ {
		if _, err := bot.PushMessage(webhookObj.Events[0].Source.UserID, linebot.NewTextMessage("Hello >> "+strconv.FormatInt(int64(i), 10))).Do(); err != nil {
			log.Println(err)
			return
		}
	}
}
