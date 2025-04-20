package main

import (
	"fmt"

	"buchstaben.go/controller"
	"buchstaben.go/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := persistence.LoadGamesFromFile(); err != nil {
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
	r.GET("/games", controller.ListGamesHandler)
	r.GET("/games/:username", controller.GetGameHandler)
	r.POST("/games/:username", controller.CreateGameHandler)
	r.POST("/games/:username/play-move", controller.PlayMoveHandler)
	r.GET("/games/end-game", controller.ListEndedGamesHandler)
	r.POST("/games/:username/end-game", controller.EndGameHandler)
	r.GET("/played-words", controller.PlayedWordsHandler)

	fmt.Println("Starting server on :8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
