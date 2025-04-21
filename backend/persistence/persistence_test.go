package persistence

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"buchstaben.go/model"
	"github.com/stretchr/testify/assert"
)

// Helper function to create temporary test file
func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_games.json")

	if content != "" {
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
	}

	return filePath
}

func TestSaveGamesToFile(t *testing.T) {
	// Setup test data
	testFilePath := createTempFile(t, "")

	// Initialize test data
	model.GlobalPersistence = model.GlobalPersistenceStruct{
		Games: map[string]model.UserGame{
			"testuser": {
				User:               "testuser",
				LastMoveTimestamp:  time.Now().Format("2006-01-02 15:04:05"),
				GameStartTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				PlayedMoves:        []model.PlayedMove{},
			},
		},
		EndedGames: []model.UserGame{},
	}

	saver := &FileDataSaver{
		FilePath: testFilePath,
	}

	// Test save
	err := saver.SaveGamesToFile()

	// Verify
	assert.NoError(t, err, "SaveGamesToFile should not return an error")

	// Check if file exists and contains data
	fileInfo, err := os.Stat(testFilePath)
	assert.NoError(t, err, "Should be able to stat the created file")
	assert.Greater(t, fileInfo.Size(), int64(0), "File should not be empty")

	// Test error case with invalid path
	invalidSaver := &FileDataSaver{
		FilePath: "/nonexistent/directory/games.json",
	}

	err = invalidSaver.SaveGamesToFile()
	assert.Error(t, err, "Should return error for invalid file path")
}

func TestLoadGamesFromFile_NonExistent(t *testing.T) {
	// Use non-existent file path
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

	// Reset global persistence to ensure it's initialized by the load function
	model.GlobalPersistence = model.GlobalPersistenceStruct{}

	saver := &FileDataSaver{
		FilePath: nonExistentPath,
	}

	// Test load from non-existent file
	err := saver.LoadGamesFromFile()

	// Verify
	assert.NoError(t, err, "Loading non-existent file should not return error")

	// Check if default structure was initialized
	assert.NotNil(t, model.GlobalPersistence.Games, "Games map should be initialized")
	assert.NotNil(t, model.GlobalPersistence.EndedGames, "EndedGames slice should be initialized")
}

func TestLoadGamesFromFile_ExistingFile(t *testing.T) {
	// Create valid JSON content
	validJSON := `{
        "games": {
            "testuser": {
                "user": "testuser",
                "last_move_timestamp": "2025-04-21 12:00:00",
                "game_start_timestamp": "2025-04-21 11:00:00",
                "letters_play_set": [],
                "played_moves": []
            }
        },
        "ended_games": []
    }`

	// Create temp file with content
	testFilePath := createTempFile(t, validJSON)

	// Reset global persistence
	model.GlobalPersistence = model.GlobalPersistenceStruct{}

	saver := &FileDataSaver{
		FilePath: testFilePath,
	}

	// Test load from existing file
	err := saver.LoadGamesFromFile()

	// Verify
	assert.NoError(t, err, "Loading valid JSON file should not return error")

	// Check if data was loaded correctly
	assert.Len(t, model.GlobalPersistence.Games, 1, "Should load 1 game")

	game, exists := model.GlobalPersistence.Games["testuser"]
	assert.True(t, exists, "Game for 'testuser' should exist")
	assert.Equal(t, "testuser", game.User, "User field should match")
	assert.Equal(t, "2025-04-21 12:00:00", game.LastMoveTimestamp, "Last move timestamp should match")
}

func TestLoadGamesFromFile_InvalidJSON(t *testing.T) {
	// Create invalid JSON content
	invalidJSON := `{
        "games": {
            "testuser": {
                "user": "testuser",
                INVALID_JSON
            }
        }
    }`

	// Create temp file with content
	testFilePath := createTempFile(t, invalidJSON)

	saver := &FileDataSaver{
		FilePath: testFilePath,
	}

	// Test load from file with invalid JSON
	err := saver.LoadGamesFromFile()

	// Verify
	assert.Error(t, err, "Loading invalid JSON should return error")
}

func TestLoadGamesFromFile_ReadError(t *testing.T) {
	// Create a temp directory instead of file
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "directory")
	err := os.Mkdir(dirPath, 0755)
	assert.NoError(t, err, "Should be able to create test directory")

	saver := &FileDataSaver{
		FilePath: dirPath, // Use directory path, which will cause a read error
	}

	// Test load from unreadable path
	err = saver.LoadGamesFromFile()

	// Verify
	assert.Error(t, err, "Reading a directory as file should return error")
}
