package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistance"
)

func ListSessionsHandler(w http.ResponseWriter, r *http.Request) {
    model.SessionsLock.Lock()
    defer model.SessionsLock.Unlock()

    listSessions := model.ListSessionsResponse{}
    listSession := model.ListSession{}

    for user := range model.Sessions {
        listSession.User = user
        listSession.LastMoveTimestamp = model.Sessions[user].LastMoveTimestamp
        listSession.SessionStartTimestamp = model.Sessions[user].SessionStartTimestamp
        listSession.RemindingLetters = logic.GetRemindingsLetterCount(model.Sessions[user].LettersPlaySet)
        listSessions.Sessions = append(listSessions.Sessions, listSession)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(listSessions)
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

    if _, exists := model.Sessions[username]; exists {
        http.Error(w, "Session already exists for this username", http.StatusBadRequest)
        return
    }
    model.Sessions[username] = model.UserSession{
        LettersPlaySet:    logic.LoadLettersPlaySet(),
        LastMoveTimestamp: time.Now().Format("2006-01-02 15:04:05"),
        SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
    }
    if err := persistance.SaveSessionsToFile(); err != nil {
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
    if _, exists := model.Sessions[username]; !exists {
        http.Error(w, "Session not found for username", http.StatusNotFound)
        return
    }
    delete(model.Sessions, username)
    if err := persistance.SaveSessionsToFile(); err != nil {
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

    userSession, exists := model.Sessions[username]
    if !exists {
        userSession = model.UserSession{
            LettersPlaySet:    logic.LoadLettersPlaySet(),
            LastMoveTimestamp: time.Now().Format("2006-01-02 15:04:05"),
            SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
        }
        model.Sessions[username] = userSession
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        model.LetterResponse{
            LetterOverAllValue: logic.GetLetterValue(userSession.LettersPlaySet),
            LettersPlaySet:     userSession.LettersPlaySet,
            LastMoveTimestamp:  userSession.LastMoveTimestamp,
        },
    )
}

func PlayMoveInputHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Processing input...")
    username := logic.GetUserNameFromResponse(*r)
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    var input struct {
        String string `json:"string"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    model.SessionsLock.Lock()
    defer model.SessionsLock.Unlock()

    _, exists := model.Sessions[username]
    if !exists {
        http.Error(w, "Session not found for username", http.StatusNotFound)
        return
    }
    lettersPlaySet := model.Sessions[username].LettersPlaySet
    newLettersPlaySet, err := logic.RemoveLetters(lettersPlaySet, input.String)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    model.Sessions[username] = model.UserSession{
        LettersPlaySet:    newLettersPlaySet,
        LastMoveTimestamp: time.Now().Format("2006-01-02 15:04:05"),
    }
    if err := persistance.SaveSessionsToFile(); err != nil {
        http.Error(w, "Failed to save session data", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        model.LetterResponse{
            LetterOverAllValue: logic.GetLetterValue(newLettersPlaySet),
            LettersPlaySet:     newLettersPlaySet,
            LastMoveTimestamp:  model.Sessions[username].LastMoveTimestamp,
            SessionStartTimestamp: model.Sessions[username].SessionStartTimestamp,
        },
    )
}

func ResetLettersHandler(w http.ResponseWriter, r *http.Request) {
    username := logic.GetUserNameFromResponse(*r)
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }
    model.SessionsLock.Lock()
    defer model.SessionsLock.Unlock()
    model.Sessions[username] = model.UserSession{
		LettersPlaySet: logic.LoadLettersPlaySet(),
		LastMoveTimestamp: time.Now().Format("2006-01-02 15:04:05"),
        SessionStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
    if err := persistance.SaveSessionsToFile(); err != nil {
        http.Error(w, "Failed to save session data", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        model.LetterResponse{
            LetterOverAllValue: logic.GetLetterValue(model.Sessions[username].LettersPlaySet),
            LettersPlaySet:       model.Sessions[username].LettersPlaySet,
			LastMoveTimestamp: model.Sessions[username].LastMoveTimestamp,
            SessionStartTimestamp: model.Sessions[username].SessionStartTimestamp,
        },
    )
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
