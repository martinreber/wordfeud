package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
)

func ListGames() []model.ListGame {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	listGames := []model.ListGame{}
	for user, game := range model.GlobalPersistence.Games {
		listGames = append(listGames, model.ListGame{
			User:               user,
			LastMoveTimestamp:  game.LastMoveTimestamp,
			GameStartTimestamp: game.GameStartTimestamp,
			RemindingLetters:   logic.GetRemindingsLetterCount(game.LettersPlaySet),
		})
	}
	return listGames
}

func CreateGame(username string) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	if _, exists := model.GlobalPersistence.Games[username]; exists {
		return fmt.Errorf("game already exists for this username")
	}

	model.GlobalPersistence.Games[username] = model.UserGame{
		User:               username,
		LettersPlaySet:     logic.LoadLettersPlaySet(),
		LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
		GameStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		LetterOverAllValue: 0, // Replace with logic.GetLetterValue() if needed
		PlayedMoves:        []model.PlayedMove{},
	}

	return persistence.SaveGamesToFile()
}

func DeleteGame(username string) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	if _, exists := model.GlobalPersistence.Games[username]; !exists {
		return fmt.Errorf("game not found for username")
	}

	delete(model.GlobalPersistence.Games, username)
	return persistence.SaveGamesToFile()
}

func EndGame(username string) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	game, exists := model.GlobalPersistence.Games[username]
	if !exists {
		return fmt.Errorf("game not found for username")
	}

	// Set the GameEndTimestamp
	game.GameEndTimestamp = time.Now().Format("2006-01-02 15:04:05")

	// Move the game to EndedGames and remove it from active games
	model.GlobalPersistence.EndedGames = append(model.GlobalPersistence.EndedGames, game)
	delete(model.GlobalPersistence.Games, username)

	return persistence.SaveGamesToFile()
}
func GetLetters(username string) (model.UserGame, error) {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	userGame, exists := model.GlobalPersistence.Games[username]
	if !exists {
		// Create a new game if it doesn't exist
		userGame = model.UserGame{
			User:               username,
			LettersPlaySet:     logic.LoadLettersPlaySet(),
			LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
			GameStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			LetterOverAllValue: logic.GetLetterValue(logic.LoadLettersPlaySet()),
			PlayedMoves:        []model.PlayedMove{},
		}
		model.GlobalPersistence.Games[username] = userGame

		// Save the new game to file
		if err := persistence.SaveGamesToFile(); err != nil {
			return model.UserGame{}, fmt.Errorf("failed to save game data: %w", err)
		}
	}

	return userGame, nil
}

func PlayMoveInput(username string, playedMoveInput model.PlayedMove) (model.UserGame, error) {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	game, exists := model.GlobalPersistence.Games[username]
	if !exists {
		return model.UserGame{}, fmt.Errorf("game not found for username")
	}

	playedMoveInput.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	newLettersPlaySet, err := logic.RemoveLetters(game.LettersPlaySet, playedMoveInput.Letters)
	if err != nil {
		return model.UserGame{}, err
	}

	updatedGame := model.UserGame{
		User:               username,
		LettersPlaySet:     newLettersPlaySet,
		LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
		GameStartTimestamp: game.GameStartTimestamp,
		LetterOverAllValue: logic.GetLetterValue(newLettersPlaySet),
		PlayedMoves:        append(game.PlayedMoves, playedMoveInput),
	}
	model.GlobalPersistence.Games[username] = updatedGame

	if err := persistence.SaveGamesToFile(); err != nil {
		return model.UserGame{}, fmt.Errorf("failed to save game data: %w", err)
	}
	return updatedGame, nil
}

func GetPlayedWords() []model.WordCount {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	wordCounts := make(map[string]int)

	// Count words in active games
	for _, game := range model.GlobalPersistence.Games {
		for _, move := range game.PlayedMoves {
			word := strings.ToLower(move.Word)
			wordCounts[word]++
		}
	}

	// Count words in ended games
	for _, game := range model.GlobalPersistence.EndedGames {
		for _, move := range game.PlayedMoves {
			word := strings.ToLower(move.Word)
			wordCounts[word]++
		}
	}

	wordsCount := make([]model.WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		wordsCount = append(wordsCount, model.WordCount{Word: word, CurrentCount: count})
	}
	sort.Slice(wordsCount, func(i, j int) bool {
		return wordsCount[i].Word < wordsCount[j].Word
	})

	return wordsCount
}
