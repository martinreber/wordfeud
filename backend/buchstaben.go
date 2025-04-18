package main

import (
	"fmt"
	"net/http"

	"buchstaben.go/controller"
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

	http.Handle("/letters", controller.EnableCORS(http.HandlerFunc(controller.GetLettersController)))
	http.Handle("/play-move", controller.EnableCORS(http.HandlerFunc(controller.PlayMoveInputController)))
	http.Handle("/reset", controller.EnableCORS(http.HandlerFunc(controller.ResetLettersController)))
	http.Handle("/list", controller.EnableCORS(http.HandlerFunc(controller.ListSessionsController)))
	http.Handle("/create", controller.EnableCORS(http.HandlerFunc(controller.CreateSessionController)))
	http.Handle("/end-session", controller.EnableCORS(http.HandlerFunc(controller.EndSessionController)))
	http.Handle("/delete", controller.EnableCORS(http.HandlerFunc(controller.DeleteSessionController)))
	http.Handle("/played-words", controller.EnableCORS(http.HandlerFunc(controller.PlayedWordsController)))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
