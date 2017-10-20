package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pikomonde/fam100bot/src/fambot"
)

var LINE LineSetting

//var GAMES fambot.GameInfo

func init() {
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

	if webhookObj[0].Type == "join" {
		EventJoin(webhookObj)
	} else if webhookObj[0].Type == "message" {
		if webhookObj[0].Message.Type == "text" {
			userMsg := webhookObj[0].Message.Text
			if userMsg == fambot.CMD_JOIN {
				EventMessageJoin(webhookObj)
			} else if userMsg == fambot.CMD_SCORE {
				EventMessageScore(webhookObj)
			} else {
				EventMessageAny(webhookObj)
			}
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

func SimulateWebhook(userID, msg string) {
	// Sample Webhook
	webhookObj := []WebhookEvent{WebhookEvent{
		ReplyToken: "QWE",
		Type:       "message",
		Timestamp:  1462629479859,
		Source: WebhookSource{
			Type:   "room", // user, group, room
			RoomID: "000000XXXXX",
			UserID: userID,
		},
		Message: WebhookMessage{
			Type: "text",
			Text: msg,
		},
	}}

	EventJoin(webhookObj)
	if webhookObj[0].Message.Type == "text" {
		userMsg := webhookObj[0].Message.Text
		if userMsg == fambot.CMD_JOIN {
			EventMessageJoin(webhookObj)
		} else if userMsg == fambot.CMD_SCORE {
			EventMessageScore(webhookObj)
		} else {
			EventMessageAny(webhookObj)
		}
	}
}

func EventJoin(webhookObj WebhookEvents) {
	// Set Variables Game, UserID, GameRoomID
	var game fambot.GameInfo
	gameRoomID := BOT_TYPE_LINE + webhookObj.GetMetaGameRoomID()

	// Load Game Info Data
	game.LoadGameInfoByRoomID(gameRoomID)

	// Add RoomID
	game.RoomID = gameRoomID

	// Call and Append Players
	// TODO: Hit the LINE API to get list of user in a room
	uID := "QQQ"
	if game.Players == nil {
		game.Players = make(map[string]fambot.PlayerInfo)
	}
	game.CreateUserIfNotListed(uID)

	// Save Game Info Data
	game.SaveGameInfoByRoomID(gameRoomID)
}

func EventMessageJoin(webhookObj WebhookEvents) {
	// Set Variables Game, UserID, GameRoomID
	var game fambot.GameInfo
	userID := BOT_TYPE_LINE + SOURCE_TYPE_USER + webhookObj[0].Source.UserID
	gameRoomID := BOT_TYPE_LINE + webhookObj.GetMetaGameRoomID()

	// Load Game Info Data
	game.LoadGameInfoByRoomID(gameRoomID)

	// Set Game Info value
	if ok, _ := game.IsStarted(); !ok {
		// Hosting the game
		if game.NumOfJoinedPlayer() == 0 {
			game.Println("User " + userID + " is the host")
			game.ResetJoinedPlayer()
			game.SetNewQuestions()
			game.UpdatedAt = time.Now()
			defer game.HostGame(gameRoomID)
		}

		// Set Join Round Info To determined number of player join
		game.CreateUserIfNotListed(userID)
		if isUpdated := game.JoinRoundIfNotJoined(userID); isUpdated {
			game.Println("User " + userID + " join the game")
		}

		// Save Game Info Data
		game.SaveGameInfoByRoomID(gameRoomID)
	}

	//game.Players = append(game.Players, fambot.PlayerInfo{
	//	PlayerID:    "qweqwe",
	//	RoomScore:   100,
	//	RoundScore:  30,
	//	IsJoinRound: true,
	//})
	//game.Players = append(game.Players, fambot.PlayerInfo{
	//	PlayerID:    "qweqweWWW",
	//	RoomScore:   1020,
	//	RoundScore:  320,
	//	IsJoinRound: false,
	//})
	//fambot.DB.Update(func(tx *bolt.Tx) error {
	//	b, err := tx.CreateBucketIfNotExists([]byte("GameRoom"))
	//	if err != nil {
	//		log.Println("[EventMessageJoin]: " + err.Error())
	//		return err
	//	}
	//
	//	v, _ := json.Marshal(game)
	//	err = b.Put([]byte("LNRM000000XXXXX"), []byte(v))
	//	if err != nil {
	//		log.Println("[EventMessageJoin]: " + err.Error())
	//		return err
	//	}
	//	return nil
	//})
	//var game2 fambot.GameInfo
	//msg := fmt.Sprintf("%v\n", game)
	//fmt.Println(msg)
	//fmt.Println(game.Players[0].PlayerID)
	//fmt.Println(game.NumOfJoinedPlayer())
	//if _, err := LINE.Bot.ReplyMessage(webhookObj.Events[0].ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
	//	log.Println("[LineWebhook] " + err.Error())
	//	return
	//}
}

func EventMessageScore(webhookObj WebhookEvents) {
	// TODO:
}

func EventMessageAny(webhookObj WebhookEvents) {
	// Set Variables Game, UserID, GameRoomID
	var game fambot.GameInfo
	userID := BOT_TYPE_LINE + SOURCE_TYPE_USER + webhookObj[0].Source.UserID
	gameRoomID := BOT_TYPE_LINE + webhookObj.GetMetaGameRoomID()

	// Load Game Info Data
	game.LoadGameInfoByRoomID(gameRoomID)

	// Set Game Info value
	if ok, _ := game.IsStarted(); ok {
		// Set Join Round Info To determined number of player join
		game.CreateUserIfNotListed(userID)
		v := game.Players[userID]
		v.IsJoinGame = true
		game.Players[userID] = v

		// TODO: Add game.Answer(userID, msg)
		game.Println("User " + userID + " answer " + webhookObj[0].Message.Text)

		// Save Game Info Data
		game.SaveGameInfoByRoomID(gameRoomID)
	}
}
