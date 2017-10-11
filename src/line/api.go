package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var LINE LineSetting

func Init() {
	LINE.ChannelSecret = os.Getenv("CHANNELSECRET")
	LINE.ChannelToken = os.Getenv("CHANNELTOKEN")
}

func Webhook(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(c.Request.Header.Get("X-Line-Signature"))
	if err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}
	hash := hmac.New(sha256.New, []byte(LINE.ChannelSecret))
	hash.Write(body)

	// Authentication
	if !hmac.Equal(decoded, hash.Sum(nil)) {
		log.Println("[LineWebhook] " + "HMAC not equal")
		return
	}

	var webhookObj WebhookEvents
	err = json.Unmarshal(body, &webhookObj)
	if err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}

	// Set Bot
	bot, err := linebot.New(LINE.ChannelSecret, LINE.ChannelToken)
	if err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}

	// Reply Message
	if _, err := bot.ReplyMessage(webhookObj.Events[0].ReplyToken, linebot.NewTextMessage(webhookObj.Events[0].Message.Text)).Do(); err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}

	// Push Message
	for i := 0; i < 250; i++ {
		if _, err := bot.PushMessage(webhookObj.Events[0].Source.UserID, linebot.NewTextMessage("Hello >> "+strconv.FormatInt(int64(i), 10))).Do(); err != nil {
			log.Println("[LineWebhook] " + err.Error())
			return
		}
	}
}
