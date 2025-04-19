package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"buchstaben.go/model"
)

const gameFilePath = "../data/games.json"

func SaveGamesToFile() error {
	fmt.Println("Saving games to file...")

	file, err := json.MarshalIndent(model.GlobalPersistence, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling games:", err)
		return err
	}
	return os.WriteFile(gameFilePath, file, 0644)
}

func LoadGamesFromFile() error {
	fmt.Println("Loading games from file...")

	if _, err := os.Stat(gameFilePath); os.IsNotExist(err) {
		model.GlobalPersistence = model.GlobalPersistenceStruct{
			Games:      make(map[model.User]model.UserGame),
			EndedGames: []model.UserGame{},
		}
		return nil
	}
	file, err := os.ReadFile(gameFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &model.GlobalPersistence)
}
