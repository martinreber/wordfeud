package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"buchstaben.go/model"
)

// DataSaver is the interface that wraps the SaveData method.
type DataSaver interface {
	SaveGamesToFile() error
	LoadGamesFromFile() error
}
type FileDataSaver struct {
	FilePath string
}

func (fds *FileDataSaver) SaveGamesToFile() error {
	fmt.Println("Saving games to file...")

	file, err := json.MarshalIndent(model.GlobalPersistence, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling games:", err)
		return err
	}
	return os.WriteFile(fds.FilePath, file, 0644)
}

func (fds *FileDataSaver) LoadGamesFromFile() error {
	fmt.Println("Loading games from file...")

	if _, err := os.Stat(fds.FilePath); os.IsNotExist(err) {
		model.GlobalPersistence = model.GlobalPersistenceStruct{
			Games:       make(map[string]model.UserGame),
			EndedGames:  []model.UserGame{},
			CustomWords: []model.CustomWord{},
		}
		return nil
	}
	file, err := os.ReadFile(fds.FilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &model.GlobalPersistence)
}
