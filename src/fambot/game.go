package fambot

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func init() {
	initDB()
}

// ==== Game Setting ====
const MINIMUM_PLAYER = 3
const CMD_JOIN = "join"
const CMD_SCORE = "score"

type GameInfo struct {
	RoomID    string
	Players   map[string]PlayerInfo
	Question  QuestionInfo
	UpdatedAt time.Time
}
type PlayerInfo struct {
	PlayerID    string
	ScoreRound  int
	ScoreRoom   int
	IsJoinRound bool
}
type QuestionInfo struct {
	QuestionID   string
	QuestionText string
	Answer       []AnswerInfo
}
type AnswerInfo struct {
	AnswerText string
	Score      int
	Answered   bool
}

func (game *GameInfo) IsStarted() bool {
	return game.NumOfJoinedPlayer() >= MINIMUM_PLAYER
}

// NumOfJoinedPlayer used to count numbers of joined player
func (game *GameInfo) NumOfJoinedPlayer() int {
	var total int
	for _, p := range game.Players {
		if p.IsJoinRound {
			total++
		}
	}
	return total
}

//func (game *GameInfo) IsExpired(ns int64) bool {
//	duration := time.Since(game.UpdatedAt)
//	isExpired := duration.Nanoseconds() > ns
//	if isExpired {
//		game.ResetJoinedPlayer()
//	}
//	return isExpired
//}

func (game *GameInfo) ResetJoinedPlayer() {
	for i, _ := range game.Players {
		v := game.Players[i]
		v.IsJoinRound = false
		game.Players[i] = v
	}
}

func (game *GameInfo) SetNewQuestions() {
	// TODO:
}

func (game *GameInfo) LoadGameInfoByRoomID(rID string) {
	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("GameRoom"))
		if err != nil {
			log.Println("[LoadGameInfoByRoomID]: " + err.Error())
			return err
		}

		v := b.Get([]byte(rID))
		err = json.Unmarshal(v, game)
		if err != nil {
			log.Println("[LoadGameInfoByRoomID]: " + err.Error())
			return err
		}
		return nil
	})
}

func (game *GameInfo) SaveGameInfoByRoomID(rID string) {
	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("GameRoom"))
		if err != nil {
			log.Println("[SaveGameInfoByRoomID]: " + err.Error())
			return err
		}

		v, err := json.Marshal(game)
		if err != nil {
			log.Println("[SaveGameInfoByRoomID]: " + err.Error())
			return err
		}

		err = b.Put([]byte(rID), []byte(v))
		if err != nil {
			log.Println("[SaveGameInfoByRoomID]: " + err.Error())
			return err
		}
		return nil
	})
}

func (game *GameInfo) JoinRoundIfNotJoined(uID string) (isUpdated bool) {
	isUpdated = false
	if !game.Players[uID].IsJoinRound {
		isUpdated = true
		v := game.Players[uID]
		v.IsJoinRound = true
		game.Players[uID] = v
	}
	return isUpdated
}

func (game *GameInfo) CreateUserIfNotListed(uID string) (isUpdated bool) {
	isUpdated = false
	if game.Players == nil {
		game.Players = make(map[string]PlayerInfo)
	}
	if _, exist := game.Players[uID]; !exist {
		isUpdated = true
		game.createUser(uID)
	}
	return isUpdated
}

func (game *GameInfo) createUser(uID string) {
	game.Players[uID] = PlayerInfo{
		PlayerID:    uID,
		ScoreRound:  0,
		ScoreRoom:   0,
		IsJoinRound: false,
	}
}

// ==== Gameplay Here ====
func (game *GameInfo) StartWaitingRoom() {
	// TODO:
}

// ==== Database Setting (BoltDB) ====
var DB *bolt.DB

func initDB() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	DB, _ = bolt.Open("my.db", 0600, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
