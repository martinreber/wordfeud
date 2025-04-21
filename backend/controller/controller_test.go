package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"buchstaben.go/model"
	"buchstaben.go/persistence"
	"buchstaben.go/service"
)

// setupTestEnvironment initializes a test environment with a temporary data file
func setupTestEnvironment(t *testing.T) (*DataController, *gin.Engine, string) {
	gin.SetMode(gin.TestMode)

	// Create temporary file for test data
	tempFile, err := os.CreateTemp("", "test-games-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()

	// Initialize empty test data
	testData := model.GlobalPersistenceStruct{
		Games:      make(map[string]model.UserGame),
		EndedGames: []model.UserGame{},
	}

	// Write initial test data
	dataBytes, _ := json.Marshal(testData)
	if err := os.WriteFile(tempFilePath, dataBytes, 0644); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Initialize the global persistence
	model.GlobalPersistence = testData

	// Create controller with actual service
	fileSaver := &persistence.FileDataSaver{FilePath: tempFilePath}
	dataService := &service.DataService{Saver: fileSaver}
	controller := &DataController{Service: dataService}

	// Setup router
	router := gin.Default()
	setupRoutes(router, controller)

	return controller, router, tempFilePath
}

// setupRoutes configures all routes for testing
func setupRoutes(router *gin.Engine, controller *DataController) {
	router.GET("/games", controller.ListGamesHandler)
	router.POST("/games/:username", controller.CreateGameHandler)
	router.GET("/games/:username", controller.GetGameHandler)
	router.POST("/games/:username/play-move", controller.PlayMoveHandler)
	router.POST("/games/:username/end-game", controller.EndGameHandler)
	router.GET("/games/end-game", controller.ListEndedGamesHandler)
	router.GET("/played-words", controller.PlayedWordsHandler)
}

// cleanupTestEnvironment removes temporary files
func cleanupTestEnvironment(t *testing.T, filePath string) {
	if err := os.Remove(filePath); err != nil {
		t.Logf("Failed to remove temp file %s: %v", filePath, err)
	}
}

func TestListGamesHandler_Empty(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Test request
	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.ListGame
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response, "Expected empty game list")
}

func TestCreateGameHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Test creating a new game
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions for successful creation
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test duplicate creation
	req = httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions for duplicate creation
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse["error"], "already exists")
}

func TestListGamesHandler_WithGames(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Now test listing games
	req = httptest.NewRequest(http.MethodGet, "/games", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.ListGame
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1, "Expected one game")
	assert.Equal(t, "testuser", response[0].User)
}

func TestGetGameHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test getting the game
	req = httptest.NewRequest(http.MethodGet, "/games/testuser", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UserGame
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.User)
	assert.NotEmpty(t, response.LettersPlaySet, "Expected letter set to be populated")
}

func TestPlayMoveHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get the game to find available letters
	req = httptest.NewRequest(http.MethodGet, "/games/testuser", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var game model.UserGame
	err := json.Unmarshal(w.Body.Bytes(), &game)
	assert.NoError(t, err)

	// Find a letter that exists in the game
	var availableLetter string
	for _, letter := range game.LettersPlaySet {
		if letter.CurrentCount > 0 {
			availableLetter = letter.Letter
			break
		}
	}
	assert.NotEmpty(t, availableLetter, "No available letters found")

	// Play a move with an available letter
	moveData := model.PlayedMove{
		Letters:        availableLetter,
		Words:          []string{"test"},
		PlayedByMyself: true,
		Points:         10,
	}
	moveJSON, _ := json.Marshal(moveData)

	req = httptest.NewRequest(http.MethodPost, "/games/testuser/play-move", bytes.NewBuffer(moveJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedGame model.UserGame
	err = json.Unmarshal(w.Body.Bytes(), &updatedGame)
	assert.NoError(t, err)
	assert.Len(t, updatedGame.PlayedMoves, 1, "Expected one played move")
}

func TestEndGameHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// End the game
	req = httptest.NewRequest(http.MethodPost, "/games/testuser/end-game", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "testuser")
	assert.Contains(t, response["message"], "ended successfully")
}

func TestListEndedGamesHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// End the game
	req = httptest.NewRequest(http.MethodPost, "/games/testuser/end-game", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// List ended games
	req = httptest.NewRequest(http.MethodGet, "/games/end-game", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.ListEndedGame
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1, "Expected one ended game")
	assert.Equal(t, "testuser", response[0].User)
}

func TestPlayedWordsHandler(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get the game to find available letters
	req = httptest.NewRequest(http.MethodGet, "/games/testuser", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var game model.UserGame
	err := json.Unmarshal(w.Body.Bytes(), &game)
	assert.NoError(t, err)

	// Find a letter that exists in the game
	var availableLetter string
	for _, letter := range game.LettersPlaySet {
		if letter.CurrentCount > 0 {
			availableLetter = letter.Letter
			break
		}
	}

	// Play a move with a word
	moveData := model.PlayedMove{
		Letters:        availableLetter,
		Words:          []string{"hello", "world"},
		PlayedByMyself: true,
		Points:         10,
	}
	moveJSON, _ := json.Marshal(moveData)

	req = httptest.NewRequest(http.MethodPost, "/games/testuser/play-move", bytes.NewBuffer(moveJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Check played words
	req = httptest.NewRequest(http.MethodGet, "/played-words", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var words []model.WordCount
	err = json.Unmarshal(w.Body.Bytes(), &words)
	assert.NoError(t, err)
	assert.Len(t, words, 2, "Expected two words")

	// Create a map of words to verify
	wordMap := make(map[string]int)
	for _, word := range words {
		wordMap[word.Word] = word.CurrentCount
	}

	assert.Equal(t, 1, wordMap["hello"], "Expected 'hello' to have count 1")
	assert.Equal(t, 1, wordMap["world"], "Expected 'world' to have count 1")
}

func TestGetGameHandler_NonExistent(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Test getting a non-existent game (should create it)
	req := httptest.NewRequest(http.MethodGet, "/games/newuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UserGame
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", response.User, "Game should be created for new user")
}

func TestPlayMoveHandler_InvalidRequest(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Send invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/games/testuser/play-move", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request body")
}

func TestPlayMoveHandler_InvalidLetter(t *testing.T) {
	_, router, tempFile := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempFile)

	// Create a game first
	req := httptest.NewRequest(http.MethodPost, "/games/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Play a move with an invalid letter
	moveData := model.PlayedMove{
		Letters:        "😊", // Invalid emoji character
		Words:          []string{"test"},
		PlayedByMyself: true,
		Points:         10,
	}
	moveJSON, _ := json.Marshal(moveData)

	req = httptest.NewRequest(http.MethodPost, "/games/testuser/play-move", bytes.NewBuffer(moveJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "not valid")
}
