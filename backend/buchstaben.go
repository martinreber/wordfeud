package main

import (
	"fmt"

	"buchstaben.go/controller"
	"buchstaben.go/persistence"
	"buchstaben.go/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const gameFilePath = "../data/games.json"

func main() {
	fileSaver := &persistence.FileDataSaver{FilePath: gameFilePath}
	dataService := service.DataService{Saver: fileSaver}
	dataController := controller.DataController{Service: &dataService}

	if err := fileSaver.LoadGamesFromFile(); err != nil {
		fmt.Println("Error loading games from file:", err)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type"}
	r.Use(cors.New(config))

	// API routes
	r.GET("/games", dataController.ListGamesHandler)
	r.GET("/games/:username", dataController.GetGameHandler)
	r.POST("/games/:username", dataController.CreateGameHandler)
	r.POST("/games/:username/play-move", dataController.PlayMoveHandler)
	r.GET("/games/end-game", dataController.ListEndedGamesHandler)
	r.POST("/games/:username/end-game", dataController.EndGameHandler)
	r.GET("/played-words", dataController.PlayedWordsHandler)

	r.GET("/custom-words", dataController.GetCustomWordsHandler)
	r.POST("/custom-words", dataController.AddCustomWordHandler)
	r.DELETE("/custom-words/:word", dataController.DeleteCustomWordHandler)

	fmt.Println("Starting server on :8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
