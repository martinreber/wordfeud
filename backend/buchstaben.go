package main

import (
	"fmt"
	"net/http"

	"buchstaben.go/handler"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
)

func main() {
	model.GlobalPersistence = model.GlobalPersistenceStruct{
		Sessions:      make(map[model.User]model.UserSession),
		EndedSessions: []model.UserSession{},
	}

	if err := persistence.LoadSessionsFromFile(); err != nil {
		fmt.Println("Error loading sessions:", err)
	}

	http.Handle("/letters", handler.EnableCORS(http.HandlerFunc(handler.GetLettersHandler)))
	http.Handle("/play-move", handler.EnableCORS(http.HandlerFunc(handler.PlayMoveInputHandler)))
	http.Handle("/reset", handler.EnableCORS(http.HandlerFunc(handler.ResetLettersHandler)))
	http.Handle("/list", handler.EnableCORS(http.HandlerFunc(handler.ListSessionsHandler)))
	http.Handle("/create", handler.EnableCORS(http.HandlerFunc(handler.CreateSessionHandler)))
	http.Handle("/end-session", handler.EnableCORS(http.HandlerFunc(handler.EndSessionHandler)))
	http.Handle("/delete", handler.EnableCORS(http.HandlerFunc(handler.DeleteSessionHandler)))
	http.Handle("/played-words", handler.EnableCORS(http.HandlerFunc(handler.PlayedWordsHandler)))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
