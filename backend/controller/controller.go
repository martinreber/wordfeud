package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"buchstaben.go/model"
	"buchstaben.go/service"
)

func ListGamesHandler(c *gin.Context) {
	listGames := service.ListGames()
	c.JSON(http.StatusOK, listGames)
}

func CreateGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := service.CreateGame(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func GetGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	userGame, err := service.GetLetters(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userGame)
}

func PlayMoveHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var playedMoveInput model.PlayedMove
	if err := c.BindJSON(&playedMoveInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedGame, err := service.PlayMoveInput(username, playedMoveInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedGame)
}

func EndGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := service.EndGame(username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Game for user '%s' ended successfully.", username)})
}

func ListEndedGamesHandler(c *gin.Context) {
	listEndedGames := service.ListEndedGames()
	c.JSON(http.StatusOK, listEndedGames)
}

func PlayedWordsHandler(c *gin.Context) {
	wordsCount := service.GetPlayedWords()
	c.JSON(http.StatusOK, wordsCount)
}
