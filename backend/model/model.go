package model

import "sync"

type LetterPlaySet struct {
	Letter        string `json:"letter"`
	OriginalCount uint   `json:"original_count"`
	CurrentCount  uint   `json:"current_count"`
	Value         uint   `json:"value"`
}

type LettersPlaySet []LetterPlaySet

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
	User                  User            `json:"user"`
	LettersPlaySet        []LetterPlaySet `json:"letters_play_set"`
	LastMoveTimestamp     string          `json:"last_move_timestamp"`
	SessionStartTimestamp string          `json:"session_start_timestamp"`
	SessionEndTimestamp   string          `json:"session_end_timestamp"`
	LetterOverAllValue    uint            `json:"letter_overall_value"`
	PlayedMoves           []PlayedMove    `json:"played_moves"`
}

type WordCount struct {
	Word         string `json:"word"`
	CurrentCount int    `json:"current_count"`
}

type GlobalPersistenceStruct struct {
	Sessions      map[User]UserSession `json:"sessions"`
	EndedSessions []UserSession        `json:"ended_sessions"`
}

var (
	GlobalPersistence GlobalPersistenceStruct
	SessionsLock      sync.Mutex
)
