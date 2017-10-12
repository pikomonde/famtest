package fambot

import "time"

// ==== Game Setting ====
const MINIMUM_PLAYER = 3
const CMD_JOIN = "join"
const CMD_SCORE = "score"

type GameInfo struct {
	RoomID    string
	Players   []PlayerInfo
	Question  QuestionInfo
	UpdatedAt time.Time
}
type PlayerInfo struct {
	PlayerID string
	Score    int
	IsJoin   bool
}
type QuestionInfo struct {
	QuestionID string
	Answered   int
}

func (game GameInfo) IsStarted() bool {
	return game.NumOfJoinedPlayer() >= MINIMUM_PLAYER
}
func (game GameInfo) NumOfJoinedPlayer() int {
	var total int
	for _, p := range game.Players {
		if p.IsJoin {
			total++
		}
	}
	return total
}
func (game GameInfo) IsExpired(ns int64) bool {
	duration := time.Since(game.UpdatedAt)
	isExpired := duration.Nanoseconds() > ns
	if isExpired {
		game.ResetJoinedPlayer()
	}
	return isExpired
}
func (game GameInfo) ResetJoinedPlayer() {
	for i, _ := range game.Players {
		game.Players[i].IsJoin = false
	}
}
