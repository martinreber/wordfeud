package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"buchstaben.go/model"
	"buchstaben.go/service"
)

type DataController struct {
	Service *service.DataService
}

func (dc *DataController) ListGamesHandler(c *gin.Context) {
	listGames := dc.Service.ListGames()
	c.JSON(http.StatusOK, listGames)
}

func (dc *DataController) CreateGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := dc.Service.CreateGame(username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (dc *DataController) GetGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	userGame, err := dc.Service.GetLetters(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userGame)
}

func (dc *DataController) PlayMoveHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var playedMove model.PlayedMove
	if err := c.BindJSON(&playedMove); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedGame, err := dc.Service.PlayMove(username, playedMove)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedGame)
}

func (dc *DataController) EndGameHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := dc.Service.EndGame(username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Game for user '%s' ended successfully.", username)})
}

func (dc *DataController) ListEndedGamesHandler(c *gin.Context) {
	listEndedGames := dc.Service.ListEndedGames()
	c.JSON(http.StatusOK, listEndedGames)
}

func (dc *DataController) PlayedWordsHandler(c *gin.Context) {
	wordsCount := dc.Service.GetPlayedWords()
	c.JSON(http.StatusOK, wordsCount)
}
