package model

import "sync"

type LetterPlaySet struct {
    Letter string `json:"letter"`
    Count  uint   `json:"count"`
    Value  uint   `json:"value"`
}

type LettersPlaySet []LetterPlaySet

type LetterResponse struct {
    LetterOverAllValue uint           `json:"letter_overall_value"`
    LettersPlaySet     LettersPlaySet `json:"letters_play_set"`
    LastMoveTimestamp  string         `json:"last_move_timestamp"`
    SessionStartTimestamp string         `json:"session_start_timestamp"`
}

type ListSessionsResponse struct {
    Sessions []ListSession `json:"sessions"`
}

type ListSession struct {
    User               User   `json:"user"`
    LastMoveTimestamp  string `json:"last_move_timestamp"`
    SessionStartTimestamp string `json:"session_start_timestamp"`
    RemindingLetters   uint   `json:"reminding_letters"`
}

type User string

type UserSession struct {
    LettersPlaySet    []LetterPlaySet `json:"letters_play_set"`
    LastMoveTimestamp string          `json:"last_move_timestamp"`
    SessionStartTimestamp string          `json:"session_start_timestamp"`
}

var (
    Sessions     map[User]UserSession
    SessionsLock sync.Mutex
)