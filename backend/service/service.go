package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
)

func ListSessions() []model.ListSession {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	listSessions := []model.ListSession{}
	for user, session := range model.GlobalPersistence.Sessions {
		listSessions = append(listSessions, model.ListSession{
			User:                  user,
			LastMoveTimestamp:     session.LastMoveTimestamp,
			SessionStartTimestamp: session.SessionStartTimestamp,
			RemindingLetters:      logic.GetRemindingsLetterCount(session.LettersPlaySet),
		})
	}
	return listSessions
}

func CreateSession(username model.User) error {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	if _, exists := model.GlobalPersistence.Sessions[username]; exists {
		return fmt.Errorf("session already exists for this username")
	}

	model.GlobalPersistence.Sessions[username] = model.UserSession{
		User:                  username,
		LettersPlaySet:        logic.LoadLettersPlaySet(),
		LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		LetterOverAllValue:    0, // Replace with logic.GetLetterValue() if needed
		PlayedMoves:           []model.PlayedMove{},
	}

	return persistence.SaveSessionsToFile()
}

func DeleteSession(username model.User) error {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	if _, exists := model.GlobalPersistence.Sessions[username]; !exists {
		return fmt.Errorf("session not found for username")
	}

	delete(model.GlobalPersistence.Sessions, username)
	return persistence.SaveSessionsToFile()
}

func EndSession(username model.User) error {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	session, exists := model.GlobalPersistence.Sessions[username]
	if !exists {
		return fmt.Errorf("session not found for username")
	}

	// Set the SessionEndTimestamp
	session.SessionEndTimestamp = time.Now().Format("2006-01-02 15:04:05")

	// Move the session to EndedSessions and remove it from active sessions
	model.GlobalPersistence.EndedSessions = append(model.GlobalPersistence.EndedSessions, session)
	delete(model.GlobalPersistence.Sessions, username)

	return persistence.SaveSessionsToFile()
}
func GetLetters(username model.User) (model.UserSession, error) {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	userSession, exists := model.GlobalPersistence.Sessions[username]
	if !exists {
		// Create a new session if it doesn't exist
		userSession = model.UserSession{
			User:                  username,
			LettersPlaySet:        logic.LoadLettersPlaySet(),
			LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
			SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			LetterOverAllValue:    logic.GetLetterValue(logic.LoadLettersPlaySet()),
			PlayedMoves:           []model.PlayedMove{},
		}
		model.GlobalPersistence.Sessions[username] = userSession

		// Save the new session to file
		if err := persistence.SaveSessionsToFile(); err != nil {
			return model.UserSession{}, fmt.Errorf("failed to save session data: %w", err)
		}
	}

	return userSession, nil
}

func PlayMoveInput(username model.User, playedMoveInput model.PlayedMove) (model.UserSession, error) {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	session, exists := model.GlobalPersistence.Sessions[username]
	if !exists {
		return model.UserSession{}, fmt.Errorf("session not found for username")
	}

	playedMoveInput.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	newLettersPlaySet, err := logic.RemoveLetters(session.LettersPlaySet, playedMoveInput.Letters)
	if err != nil {
		return model.UserSession{}, err
	}

	updatedSession := model.UserSession{
		User:                  username,
		LettersPlaySet:        newLettersPlaySet,
		LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SessionStartTimestamp: session.SessionStartTimestamp,
		LetterOverAllValue:    logic.GetLetterValue(newLettersPlaySet),
		PlayedMoves:           append(session.PlayedMoves, playedMoveInput),
	}
	model.GlobalPersistence.Sessions[username] = updatedSession

	if err := persistence.SaveSessionsToFile(); err != nil {
		return model.UserSession{}, fmt.Errorf("failed to save session data: %w", err)
	}
	return updatedSession, nil
}

func GetPlayedWords() []model.WordCount {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	wordCounts := make(map[string]int)

	// Count words in active sessions
	for _, session := range model.GlobalPersistence.Sessions {
		for _, move := range session.PlayedMoves {
			word := strings.ToLower(move.Word)
			wordCounts[word]++
		}
	}

	// Count words in ended sessions
	for _, session := range model.GlobalPersistence.EndedSessions {
		for _, move := range session.PlayedMoves {
			word := strings.ToLower(move.Word)
			wordCounts[word]++
		}
	}

	wordsCount := make([]model.WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		wordsCount = append(wordsCount, model.WordCount{Word: word, CurrentCount: count})
	}
	sort.Slice(wordsCount, func(i, j int) bool {
		return wordsCount[i].Word < wordsCount[j].Word
	})

	return wordsCount
}
