package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
)

func ListSessionsHandler(w http.ResponseWriter, r *http.Request) {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	listSessions := model.ListSessions{}
	listSession := model.ListSession{}

	for user := range model.GlobalPersistence.Sessions {
		listSession.User = user
		listSession.LastMoveTimestamp = model.GlobalPersistence.Sessions[user].LastMoveTimestamp
		listSession.SessionStartTimestamp = model.GlobalPersistence.Sessions[user].SessionStartTimestamp
		listSession.RemindingLetters = logic.GetRemindingsLetterCount(model.GlobalPersistence.Sessions[user].LettersPlaySet)
		listSessions.Sessions = append(listSessions.Sessions, listSession)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listSessions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating new session...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	if _, exists := model.GlobalPersistence.Sessions[username]; exists {
		http.Error(w, "Session already exists for this username", http.StatusBadRequest)
		return
	}
	model.GlobalPersistence.Sessions[username] = model.UserSession{
		User:                  username,
		LettersPlaySet:        logic.LoadLettersPlaySet(),
		LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		LetterOverAllValue:    logic.GetLetterValue(logic.LoadLettersPlaySet()),
		PlayedMoves:           []model.PlayedMove{},
	}
	if err := persistence.SaveSessionsToFile(); err != nil {
		http.Error(w, "Failed to save session data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete session")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if _, exists := model.GlobalPersistence.Sessions[username]; !exists {
		http.Error(w, "Session not found for username", http.StatusNotFound)
		return
	}
	delete(model.GlobalPersistence.Sessions, username)
	if err := persistence.SaveSessionsToFile(); err != nil {
		http.Error(w, "Failed to save session data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetLettersHandler(w http.ResponseWriter, r *http.Request) {
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	userSession, exists := model.GlobalPersistence.Sessions[username]
	if !exists {
		userSession = model.UserSession{
			User:                  username,
			LettersPlaySet:        logic.LoadLettersPlaySet(),
			LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
			SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			LetterOverAllValue:    logic.GetLetterValue(logic.LoadLettersPlaySet()),
			PlayedMoves:           []model.PlayedMove{},
		}
		model.GlobalPersistence.Sessions[username] = userSession
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userSession)
}

func PlayMoveInputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing input...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var playedMoveInput = model.PlayedMove{}
	if err := json.NewDecoder(r.Body).Decode(&playedMoveInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	playedMoveInput.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Received input: %+v\n", playedMoveInput)

	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	if _, exists := model.GlobalPersistence.Sessions[username]; !exists {
		http.Error(w, "Session not found for username", http.StatusNotFound)
		return
	}
	lettersPlaySet := model.GlobalPersistence.Sessions[username].LettersPlaySet
	newLettersPlaySet, err := logic.RemoveLetters(lettersPlaySet, playedMoveInput.Letters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedSession := model.UserSession{
		User:                  username,
		LettersPlaySet:        newLettersPlaySet,
		LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SessionStartTimestamp: model.GlobalPersistence.Sessions[username].SessionStartTimestamp,
		LetterOverAllValue:    logic.GetLetterValue(newLettersPlaySet),
		PlayedMoves:           append(model.GlobalPersistence.Sessions[username].PlayedMoves, playedMoveInput),
	}
	model.GlobalPersistence.Sessions[username] = updatedSession

	if err := persistence.SaveSessionsToFile(); err != nil {
		http.Error(w, "Failed to save session data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSession); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ResetLettersHandler(w http.ResponseWriter, r *http.Request) {
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()
	newSession := model.UserSession{
		User:                  username,
		LettersPlaySet:        logic.LoadLettersPlaySet(),
		LastMoveTimestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		LetterOverAllValue:    logic.GetLetterValue(logic.LoadLettersPlaySet()),
		PlayedMoves:           []model.PlayedMove{},
	}
	model.GlobalPersistence.Sessions[username] = newSession

	if err := persistence.SaveSessionsToFile(); err != nil {
		http.Error(w, "Failed to save session data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newSession); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func EndSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ending session...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	session, exists := model.GlobalPersistence.Sessions[username]
	if !exists {
		http.Error(w, "Session not found for username", http.StatusNotFound)
		return
	}

	// Set the SessionEndTimestamp
	session.SessionEndTimestamp = time.Now().Format("2006-01-02 15:04:05")

	// Move the session to EndedSessions and remove it from active sessions
	model.GlobalPersistence.EndedSessions = append(model.GlobalPersistence.EndedSessions, session)
	delete(model.GlobalPersistence.Sessions, username)

	// Save the updated sessions to file
	if err := persistence.SaveSessionsToFile(); err != nil {
		http.Error(w, "Failed to save session data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Session for user '%s' ended successfully.", username)
}

func PlayedWordsHandler(w http.ResponseWriter, r *http.Request) {
	model.SessionsLock.Lock()
	defer model.SessionsLock.Unlock()

	wordCounts := make(map[string]int)
	for _, session := range model.GlobalPersistence.Sessions {
		for _, move := range session.PlayedMoves {
			word := strings.ToLower(move.Word)
			wordCounts[word]++
		}
	}
	fmt.Printf("Word counts: %v\n", wordCounts)
	wordsCount := make([]model.WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		wordsCount = append(wordsCount, model.WordCount{Word: word, CurrentCount: count})
	}
	sort.Slice(wordsCount, func(i, j int) bool {
		return wordsCount[i].Word < wordsCount[j].Word
	})
	fmt.Printf("wordCounts: %v\n", wordCounts)
	fmt.Printf("Sorted word counts: %+v\n", wordsCount)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wordsCount); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
