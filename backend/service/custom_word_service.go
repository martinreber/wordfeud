package service

import (
	"fmt"
	"time"

	"buchstaben.go/model"
)

// Add these functions to your DataService struct

func (ds *DataService) AddCustomWord(newWord model.CustomWord) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	// Check if word already exists
	for _, word := range model.GlobalPersistence.CustomWords {
		if word.Word == newWord.Word {
			return fmt.Errorf("word '%s' already exists", newWord.Word)
		}
	}

	// Add timestamp if not provided
	if newWord.Timestamp == "" {
		newWord.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	model.GlobalPersistence.CustomWords = append(model.GlobalPersistence.CustomWords, newWord)
	return ds.Saver.SaveGamesToFile()
}

func (ds *DataService) GetCustomWords() []model.CustomWord {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	return model.GlobalPersistence.CustomWords
}

func (ds *DataService) DeleteCustomWord(word string) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	found := false
	for i, w := range model.GlobalPersistence.CustomWords {
		if w.Word == word {
			// Remove the word by slicing
			model.GlobalPersistence.CustomWords = append(
				model.GlobalPersistence.CustomWords[:i],
				model.GlobalPersistence.CustomWords[i+1:]...,
			)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("word '%s' not found", word)
	}

	return ds.Saver.SaveGamesToFile()
}
