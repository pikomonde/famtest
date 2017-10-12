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

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pikomonde/fam100bot/src/fambot"
)

var LINE LineSetting
var GAMES []fambot.GameInfo

func Init() {
	LINE.ChannelSecret = os.Getenv("CHANNELSECRET")
	LINE.ChannelToken = os.Getenv("CHANNELTOKEN")

	// Set Bot
	var err error
	LINE.Bot, _ = linebot.New(LINE.ChannelSecret, LINE.ChannelToken)
	if err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}
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

	if webhookObj.Events[0].Message.Type == "text" {
		userMsg := webhookObj.Events[0].Message.Text
		if userMsg == fambot.CMD_JOIN {
			EventMessageJoin(webhookObj)
		} else if userMsg == fambot.CMD_SCORE {
			EventMessageScore(webhookObj)
		}
	}

	// Set Bot
	//bot, err := linebot.New(LINE.ChannelSecret, LINE.ChannelToken)
	//if err != nil {
	//log.Println("[LineWebhook] " + err.Error())
	//return
	//}
	//
	//// Reply Message
	//if _, err := bot.ReplyMessage(webhookObj.Events[0].ReplyToken, linebot.NewTextMessage(webhookObj.Events[0].Message.Text)).Do(); err != nil {
	//log.Println("[LineWebhook] " + err.Error())
	//return
	//}
	//
	//// Push Message
	//for i := 0; i < 250; i++ {
	//if _, err := bot.PushMessage(webhookObj.Events[0].Source.UserID, linebot.NewTextMessage("Hello >> "+strconv.FormatInt(int64(i), 10))).Do(); err != nil {
	//	log.Println("[LineWebhook] " + err.Error())
	//	return
	//}
	//}
}

func EventMessageJoin(webhookObj WebhookEvents) {
	var v []byte
	fambot.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			log.Println("[EventMessageJoin]: " + err.Error())
			return err
		}
		err = b.Put([]byte("answer"), []byte("42"))
		if err != nil {
			log.Println("[EventMessageJoin]: " + err.Error())
			return err
		}
		v = b.Get([]byte("answer"))
		return nil
	})
	msg := fmt.Sprintf("The answer is: %s\n", v)
	if _, err := LINE.Bot.ReplyMessage(webhookObj.Events[0].ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
		log.Println("[LineWebhook] " + err.Error())
		return
	}
}

func EventMessageScore(webhookObj WebhookEvents) {

}
