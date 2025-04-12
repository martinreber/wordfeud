package model

import "sync"

type LetterPlaySet struct {
	Letter string `json:"letter"`
	Count  uint   `json:"count"`
	Value  uint   `json:"value"`
}

type LettersPlaySet []LetterPlaySet

type ListSessions struct {
	Sessions []ListSession `json:"sessions"`
}

type ListSession struct {
	User                  User   `json:"user"`
	LastMoveTimestamp     string `json:"last_move_timestamp"`
	SessionStartTimestamp string `json:"session_start_timestamp"`
	RemindingLetters      uint   `json:"reminding_letters"`
}

type User string

type PlayedMove struct {
	Letters        string `json:"letters"`
	Word           string `json:"word"`
	PlayedByMyself bool   `json:"played_by_myself"`
	Timestamp      string `json:"timestamp"`
}

type UserSession struct {
	LettersPlaySet        []LetterPlaySet `json:"letters_play_set"`
	LastMoveTimestamp     string          `json:"last_move_timestamp"`
	SessionStartTimestamp string          `json:"session_start_timestamp"`
	LetterOverAllValue    uint            `json:"letter_overall_value"`
	PlayedMoves           []PlayedMove    `json:"played_moves"`
}

var (
	Sessions     map[User]UserSession
	SessionsLock sync.Mutex
)
