package fambot

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func init() {
	initDB()
}

var ConsoleVersion = false

// ==== Game Setting ====
const (
	MinimumPlayer     = 3
	RoundPerGame      = 3
	QuorumDuration    = 120 * time.Second
	RoundDuration     = 90 * time.Second
	DelayBetweenRound = 5 * time.Second
	TickDuration      = 10 * time.Second

	CMD_JOIN  = "join"
	CMD_SCORE = "score"
)

type GameInfo struct {
	RoomID    string
	Players   map[string]PlayerInfo
	Round     RoundInfo
	UpdatedAt time.Time
}
type PlayerInfo struct {
	PlayerID   string
	ScoreGame  int
	ScoreRoom  int
	IsJoinGame bool
}
type RoundInfo struct {
	QuestionID   string
	QuestionText string
	Answer       []AnswerInfo
	IsStarted    bool
}
type AnswerInfo struct {
	AnswerText string
	Score      int
	Answered   bool
}

func (game *GameInfo) IsStarted() (bool, int) {
	// TODO: change condition to IsStarter, also consider for isRoundStarted
	return game.NumOfJoinedPlayer() >= MinimumPlayer, 0
}

func (game *GameInfo) IsLastRound() bool {
	// TODO: implement isLastRound()
	return true
}

// NumOfJoinedPlayer used to count numbers of joined player
func (game *GameInfo) NumOfJoinedPlayer() int {
	var total int
	for _, p := range game.Players {
		if p.IsJoinGame {
			total++
		}
	}
	return total
}

func (game *GameInfo) ResetJoinedPlayer() {
	for i, _ := range game.Players {
		v := game.Players[i]
		v.IsJoinGame = false
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
			log.Println("[LoadGameInfoByRoomID][CreateBucketIfNotExists]: " + err.Error())
			return err
		}

		v := b.Get([]byte(rID))
		err = json.Unmarshal(v, game)
		if err != nil {
			log.Println("[LoadGameInfoByRoomID][Unmarshal game]: " + err.Error())
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
	if !game.Players[uID].IsJoinGame {
		isUpdated = true
		v := game.Players[uID]
		v.IsJoinGame = true
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
		PlayerID:   uID,
		ScoreGame:  0,
		ScoreRoom:  0,
		IsJoinGame: false,
	}
}

func (game *GameInfo) Println(msg string) {
	if ConsoleVersion {
		fmt.Println("==> " + msg)
	}
}

// ==== Gameplay Here ====
func (game *GameInfo) HostGame(rID string) {
	// TODO: Learn how to make threadng (event based)
	tick := time.Tick(TickDuration)
	quorumEnd := time.After(QuorumDuration)
	var roundEnd <-chan time.Time
	for {
		select {
		case <-quorumEnd:
			game.LoadGameInfoByRoomID(rID)
			if ok, _ := game.IsStarted(); !ok {
				game.Println("WAKTU HABIS, PERMAINAN DIBATALKAN")
				//game.ResetJoinedPlayer()
			}
		case <-roundEnd:
			game.LoadGameInfoByRoomID(rID)
			if ok, _ := game.IsStarted(); !ok {
				if !game.IsLastRound() {
					game.Println("WAKTU HABIS, RONDE SELANJUTNYA")
				} else {
					game.Println("WAKTU HABIS, PERMAIANAN USAI")
					//game.PrintRoundScore()
					//game.ResetJoinedPlayer()
				}
			}

			roundEnd = time.After(RoundDuration)
		case <-tick:
			game.LoadGameInfoByRoomID(rID)
			if ok, _ := game.IsStarted(); !ok {
				fmt.Println("BELOM MULAI")
			} else {
				if roundEnd == nil {
					roundEnd = time.After(RoundDuration)
				}
			}
		}
	}
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
