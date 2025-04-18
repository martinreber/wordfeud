package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
	"buchstaben.go/service"
)

func ListSessionsController(w http.ResponseWriter, r *http.Request) {
	listSessions := service.ListSessions()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listSessions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateSessionController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating new session...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.CreateSession(username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteSessionController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete session")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.DeleteSession(username); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetLettersController(w http.ResponseWriter, r *http.Request) {
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	userSession, err := service.GetLetters(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userSession); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func PlayMoveInputController(w http.ResponseWriter, r *http.Request) {
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var playedMoveInput model.PlayedMove
	if err := json.NewDecoder(r.Body).Decode(&playedMoveInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedSession, err := service.PlayMoveInput(username, playedMoveInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSession); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ResetLettersController(w http.ResponseWriter, r *http.Request) {
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

func EndSessionController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ending session...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.EndSession(username); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Session for user '%s' ended successfully.", username)
}

func PlayedWordsController(w http.ResponseWriter, r *http.Request) {
	wordsCount := service.GetPlayedWords()

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
