package model

import "sync"

type LetterPlaySet struct {
	Letter        string `json:"letter"`
	OriginalCount uint   `json:"original_count"`
	CurrentCount  uint   `json:"current_count"`
	Value         uint   `json:"value"`
}

type LettersPlaySet []LetterPlaySet

type ListGame struct {
	User               string `json:"user"`
	LastMoveTimestamp  string `json:"last_move_timestamp"`
	GameStartTimestamp string `json:"game_start_timestamp"`
	RemindingLetters   uint   `json:"reminding_letters"`
}

type ListEndedGame struct {
	User               string `json:"user"`
	LastMoveTimestamp  string `json:"last_move_timestamp"`
	GameStartTimestamp string `json:"game_start_timestamp"`
}

type PlayedMove struct {
	Letters        string   `json:"letters"`
	Words          []string `json:"words"`
	PlayedByMyself bool     `json:"played_by_myself"`
	Timestamp      string   `json:"timestamp"`
	Points         uint     `json:"points"`
}

type UserGame struct {
	User               string          `json:"user"`
	LettersPlaySet     []LetterPlaySet `json:"letters_play_set"`
	LastMoveTimestamp  string          `json:"last_move_timestamp"`
	GameStartTimestamp string          `json:"game_start_timestamp"`
	GameEndTimestamp   string          `json:"game_end_timestamp"`
	LetterOverAllValue uint            `json:"letter_overall_value"`
	PlayedMoves        []PlayedMove    `json:"played_moves"`
}

type WordCount struct {
	Word         string `json:"word"`
	CurrentCount int    `json:"current_count"`
}

type GlobalPersistenceStruct struct {
	Games      map[string]UserGame `json:"games"`
	EndedGames []UserGame          `json:"ended_games"`
}

var (
	GlobalPersistence GlobalPersistenceStruct
	GamesLock         sync.Mutex
)
