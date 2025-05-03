package service

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"buchstaben.go/logic"
	"buchstaben.go/model"
	"buchstaben.go/persistence"
)

type DataService struct {
	Saver persistence.DataSaver
}

func (ds *DataService) ListGames() []model.ListGame {
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

func (ds *DataService) CreateGame(username string) error {
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
		LetterOverAllValue: 0,
		PlayedMoves:        []model.PlayedMove{},
	}
	return ds.Saver.SaveGamesToFile()
}

func (ds *DataService) DeleteGame(username string) error {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	if _, exists := model.GlobalPersistence.Games[username]; !exists {
		return fmt.Errorf("game not found for username")
	}

	delete(model.GlobalPersistence.Games, username)
	return ds.Saver.SaveGamesToFile()
}

func (ds *DataService) EndGame(username string) error {
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

	return ds.Saver.SaveGamesToFile()
}
func (ds *DataService) GetLetters(username string) (model.UserGame, error) {
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
		if err := ds.Saver.SaveGamesToFile(); err != nil {
			return model.UserGame{}, fmt.Errorf("failed to save game data: %w", err)
		}
	}

	return userGame, nil
}

func (ds *DataService) PlayMove(username string, playedMove model.PlayedMove) (model.UserGame, error) {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	game, exists := model.GlobalPersistence.Games[username]
	if !exists {
		return model.UserGame{}, fmt.Errorf("game not found for username")
	}

	playedMove.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	newLettersPlaySet, err := logic.RemoveLetters(game.LettersPlaySet, playedMove.Letters)
	if err != nil {
		return model.UserGame{}, err
	}

	updatedGame := model.UserGame{
		User:               username,
		LettersPlaySet:     newLettersPlaySet,
		LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
		GameStartTimestamp: game.GameStartTimestamp,
		LetterOverAllValue: logic.GetLetterValue(newLettersPlaySet),
		PlayedMoves:        append(game.PlayedMoves, playedMove),
	}
	model.GlobalPersistence.Games[username] = updatedGame

	if err := ds.Saver.SaveGamesToFile(); err != nil {
		return model.UserGame{}, fmt.Errorf("failed to save game data: %w", err)
	}
	return updatedGame, nil
}

func (ds *DataService) ListEndedGames() []model.ListEndedGame {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	listEndedGames := []model.ListEndedGame{}
	for _, endedGame := range model.GlobalPersistence.EndedGames {
		listEndedGames = append(listEndedGames, model.ListEndedGame{
			User:               endedGame.User,
			LastMoveTimestamp:  endedGame.LastMoveTimestamp,
			GameStartTimestamp: endedGame.GameStartTimestamp,
		})
	}
	return listEndedGames
}

func (ds *DataService) GetPlayedWords(letters string) []model.WordCount {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	unfilteredWordCounts := make(map[string]int)

	// Count words in active games
	for _, game := range model.GlobalPersistence.Games {
		countWords(game.PlayedMoves, unfilteredWordCounts)
	}

	// Count words in ended games
	for _, endedGame := range model.GlobalPersistence.EndedGames {
		countWords(endedGame.PlayedMoves, unfilteredWordCounts)
	}

	// add words from DWDS Word list
	for word := range model.GlobalWordMap {
		if len(word) > 15 || len(word) < 2 {
			continue
		}
		// if word contains alpha characters and hyphen, keep it, otherwise skip
		if !regexp.MustCompile(`^[a-zA-Z-]+$`).MatchString(word) {
			continue
		}
		if _, exists := unfilteredWordCounts[word]; !exists {
			unfilteredWordCounts[word] = 0
		}
	}

	// convert map to slice of WordCount
	wordsCount := make([]model.WordCount, 0, len(unfilteredWordCounts))
	for word := range unfilteredWordCounts {
		if !buildWordOutOfLetters(word, letters) {
			continue
		}
		wordsCount = append(wordsCount, model.WordCount{
			Word:         word,
			CurrentCount: 0,
		})
	}
	// Sort the wordsCount slice by the length of the words in descending order,
	// and for words with the same length, sort them alphabetically in ascending order.
	sort.Slice(wordsCount, func(i, j int) bool {
		if len(wordsCount[i].Word) == len(wordsCount[j].Word) {
			return wordsCount[i].Word < wordsCount[j].Word // Sort alphabetically if lengths are equal
		}
		return len(wordsCount[i].Word) > len(wordsCount[j].Word) // Sort by length in descending order
	})
	if len(wordsCount) > 100 {
		wordsCount = wordsCount[:100]
	}

	return wordsCount
}

func (ds *DataService) FindWords(letters string) []model.WordCount {
	model.GamesLock.Lock()
	defer model.GamesLock.Unlock()

	unfilteredWordCounts := make(map[string]int)

	// add words from DWDS Word list
	for word := range model.GlobalWordMap {
		if len(word) > 15 || len(word) < 2 {
			continue
		}
		// if word contains alpha characters and hyphen, keep it, otherwise skip
		if !regexp.MustCompile(`^[a-zA-Z-]+$`).MatchString(word) {
			continue
		}
		if _, exists := unfilteredWordCounts[word]; !exists {
			unfilteredWordCounts[word] = 0
		}
	}

	// convert map to slice of WordCount
	wordsCount := make([]model.WordCount, 0, len(unfilteredWordCounts))
	for word := range unfilteredWordCounts {
		if !buildWordOutOfLetters(word, letters) {
			continue
		}
		wordsCount = append(wordsCount, model.WordCount{
			Word:         word,
			CurrentCount: 0,
		})
	}
	// Sort the wordsCount slice by the length of the words in descending order,
	// and for words with the same length, sort them alphabetically in ascending order.
	sort.Slice(wordsCount, func(i, j int) bool {
		if len(wordsCount[i].Word) == len(wordsCount[j].Word) {
			return wordsCount[i].Word < wordsCount[j].Word // Sort alphabetically if lengths are equal
		}
		return len(wordsCount[i].Word) > len(wordsCount[j].Word) // Sort by length in descending order
	})
	if len(wordsCount) > 100 {
		wordsCount = wordsCount[:100]
	}

	return wordsCount
}

// check if the word can be built out of the letters
func buildWordOutOfLetters(word, letters string) bool {
	// Create a map to count occurrences of each character in the letters string
	letterCounts := make(map[rune]int)
	for _, char := range letters {
		letterCounts[char]++
	}

	// Check if each character in the word can be formed using the letters
	for _, char := range word {
		if letterCounts[char] > 0 {
			letterCounts[char]-- // Decrement the count for the character
		} else {
			return false // Character not found or insufficient occurrences
		}
	}

	return true
}

func countWords(playedMoves []model.PlayedMove, wordCounts map[string]int) {
	for _, move := range playedMoves {
		if move.Words != nil {
			for _, word := range move.Words {
				word = strings.ToLower(word)
				wordCounts[word]++
			}
		}
	}
}
