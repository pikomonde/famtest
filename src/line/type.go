package line

import "github.com/line/line-bot-sdk-go/linebot"

type LineSetting struct {
	ChannelSecret string
	ChannelToken  string
	Bot           *linebot.Client
}

// ==== Webhook ====
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
