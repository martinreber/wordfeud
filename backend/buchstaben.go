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
		Games:      make(map[model.User]model.UserGame),
		EndedGames: []model.UserGame{},
	}

	if err := persistence.LoadGamesFromFile(); err != nil {
		fmt.Println("Error loading games:", err)
	}

	http.Handle("/letters", controller.EnableCORS(http.HandlerFunc(controller.GetLettersController)))
	http.Handle("/play-move", controller.EnableCORS(http.HandlerFunc(controller.PlayMoveInputController)))
	http.Handle("/reset", controller.EnableCORS(http.HandlerFunc(controller.ResetLettersController)))
	http.Handle("/list", controller.EnableCORS(http.HandlerFunc(controller.ListGamesController)))
	http.Handle("/create", controller.EnableCORS(http.HandlerFunc(controller.CreateGameController)))
	http.Handle("/end-game", controller.EnableCORS(http.HandlerFunc(controller.EndGameController)))
	http.Handle("/delete", controller.EnableCORS(http.HandlerFunc(controller.DeleteGameController)))
	http.Handle("/played-words", controller.EnableCORS(http.HandlerFunc(controller.PlayedWordsController)))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
