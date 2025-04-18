package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"buchstaben.go/model"
)

const sessionFilePath = "../data/sessions.json"

func SaveSessionsToFile() error {
	fmt.Println("Saving sessions to file...")

	file, err := json.MarshalIndent(model.GlobalPersistence, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling sessions:", err)
		return err
	}
	return os.WriteFile(sessionFilePath, file, 0644)
}

func LoadSessionsFromFile() error {
	fmt.Println("Loading sessions from file...")

	if _, err := os.Stat(sessionFilePath); os.IsNotExist(err) {
		model.GlobalPersistence = model.GlobalPersistenceStruct{
			Sessions:      make(map[model.User]model.UserSession),
			EndedSessions: []model.UserSession{},
		}
		return nil
	}
	file, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &model.GlobalPersistence)
}
