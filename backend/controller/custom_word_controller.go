package controller

import (
	"fmt"
	"net/http"

	"buchstaben.go/model"
	"github.com/gin-gonic/gin"
)

func (dc *DataController) GetCustomWordsHandler(c *gin.Context) {
	customWords := dc.Service.GetCustomWords()
	c.JSON(http.StatusOK, customWords)
}

func (dc *DataController) AddCustomWordHandler(c *gin.Context) {
	var newWords model.CustomWords
	if err := c.BindJSON(&newWords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := dc.Service.AddCustomWords(newWords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (dc *DataController) DeleteCustomWordHandler(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Word is required"})
		return
	}

	if err := dc.Service.DeleteCustomWord(word); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Word '%s' deleted successfully.", word)})
}
