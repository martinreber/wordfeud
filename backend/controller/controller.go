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

func ListGamesController(w http.ResponseWriter, r *http.Request) {
	listGames := service.ListGames()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listGames); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateGameController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating new game...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.CreateGame(username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteGameController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete game")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.DeleteGame(username); err != nil {
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
	userGame, err := service.GetLetters(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userGame); err != nil {
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

	updatedGame, err := service.PlayMoveInput(username, playedMoveInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedGame); err != nil {
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
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()
	newGame := model.UserGame{
		User:               username,
		LettersPlaySet:     logic.LoadLettersPlaySet(),
		LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
		GameStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		LetterOverAllValue: logic.GetLetterValue(logic.LoadLettersPlaySet()),
		PlayedMoves:        []model.PlayedMove{},
	}
	model.GlobalPersistence.Games[username] = newGame

	if err := persistence.SaveGamesToFile(); err != nil {
		http.Error(w, "Failed to save game data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newGame); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func EndGameController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ending game...")
	username := logic.GetUserNameFromResponse(*r)
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if err := service.EndGame(username); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Game for user '%s' ended successfully.", username)
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
