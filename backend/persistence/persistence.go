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
	LoadWordListFromFile() error
}
type FileDataSaver struct {
	GameFilePath     string
	WordListFilePath string
}

func (fds *FileDataSaver) SaveGamesToFile() error {
	fmt.Println("Saving games to file...")

	file, err := json.MarshalIndent(model.GlobalPersistence, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling games:", err)
		return err
	}
	return os.WriteFile(fds.GameFilePath, file, 0644)
}

func (fds *FileDataSaver) LoadGamesFromFile() error {
	fmt.Println("Loading game from file...")

	if _, err := os.Stat(fds.GameFilePath); os.IsNotExist(err) {
		model.GlobalPersistence = model.GlobalPersistenceStruct{
			Games:       make(map[string]model.UserGame),
			EndedGames:  []model.UserGame{},
			CustomWords: []model.CustomWord{},
		}
		return nil
	}
	file, err := os.ReadFile(fds.GameFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &model.GlobalPersistence)
}

func (fds *FileDataSaver) LoadWordListFromFile() error {
	fmt.Println("Loading word list from file...")

	if _, err := os.Stat(fds.WordListFilePath); os.IsNotExist(err) {
		fmt.Println("Word list file does not exist")
		model.GlobalWordMap = model.WordMap{}
		return nil
	}
	file, err := os.ReadFile(fds.WordListFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &model.GlobalWordMap)
}
