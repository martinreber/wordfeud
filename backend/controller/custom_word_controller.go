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
	var newWord model.CustomWord
	if err := c.BindJSON(&newWord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if newWord.Word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Word cannot be empty"})
		return
	}

	if err := dc.Service.AddCustomWord(newWord); err != nil {
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
