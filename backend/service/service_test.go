package service

import (
	"fmt"
	"reflect"
	"testing"

	"buchstaben.go/logic"
	"buchstaben.go/model"
)

// MockDataSaver implements the DataSaver interface for testing
type MockDataSaver struct {
	SaveCalled bool
	SaveError  error
	LoadCalled bool
	LoadError  error
}

func (m *MockDataSaver) SaveGamesToFile() error {
	m.SaveCalled = true
	return m.SaveError
}

func (m *MockDataSaver) LoadGamesFromFile() error {
	m.LoadCalled = true
	return m.LoadError
}

func setupTestEnvironment() (*DataService, *MockDataSaver) {
	// Reset global state
	model.GlobalPersistence = model.GlobalPersistenceStruct{
		Games:      make(map[string]model.UserGame),
		EndedGames: []model.UserGame{},
	}

	mock := &MockDataSaver{}
	service := &DataService{
		Saver: mock,
	}

	return service, mock
}

func TestListGames(t *testing.T) {
	service, _ := setupTestEnvironment()

	// Test with empty games
	games := service.ListGames()
	if len(games) != 0 {
		t.Errorf("Expected 0 games, got %d", len(games))
	}

	// Add a test game
	model.GlobalPersistence.Games["testuser"] = model.UserGame{
		User:               "testuser",
		LastMoveTimestamp:  "2025-04-21 12:00:00",
		GameStartTimestamp: "2025-04-21 11:00:00",
		LettersPlaySet:     []model.LetterPlaySet{{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1}},
	}

	// Test with one game
	games = service.ListGames()
	if len(games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(games))
	}
	if games[0].User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", games[0].User)
	}
	if games[0].RemindingLetters != 3 {
		t.Errorf("Expected 3 remaining letters, got %d", games[0].RemindingLetters)
	}
}

func TestCreateGame(t *testing.T) {
	service, mock := setupTestEnvironment()

	// Test successful creation
	err := service.CreateGame("testuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.SaveCalled {
		t.Error("Expected SaveGamesToFile to be called")
	}

	game, exists := model.GlobalPersistence.Games["testuser"]
	if !exists {
		t.Error("Game was not created")
		return
	}

	if game.User != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", game.User)
	}

	// Test duplicate user creation
	mock.SaveCalled = false // Reset flag
	err = service.CreateGame("testuser")

	if err == nil {
		t.Error("Expected error for duplicate user, got nil")
	}
	if mock.SaveCalled {
		t.Error("SaveGamesToFile should not be called for duplicate user")
	}

	// Test error during save
	mock.SaveError = fmt.Errorf("save error")
	err = service.CreateGame("newuser")

	if err == nil || err.Error() != "save error" {
		t.Errorf("Expected 'save error', got %v", err)
	}
}

func TestDeleteGame(t *testing.T) {
	service, mock := setupTestEnvironment()

	// Add a test game
	model.GlobalPersistence.Games["testuser"] = model.UserGame{
		User: "testuser",
	}

	// Test successful deletion
	err := service.DeleteGame("testuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.SaveCalled {
		t.Error("Expected SaveGamesToFile to be called")
	}

	if _, exists := model.GlobalPersistence.Games["testuser"]; exists {
		t.Error("Game was not deleted")
	}

	// Test non-existent user
	mock.SaveCalled = false
	err = service.DeleteGame("nonexistent")

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
	if mock.SaveCalled {
		t.Error("SaveGamesToFile should not be called for non-existent user")
	}
}

func TestEndGame(t *testing.T) {
	service, mock := setupTestEnvironment()

	// Test ending non-existent game
	err := service.EndGame("nonexistent")

	if err == nil {
		t.Error("Expected error for non-existent game, got nil")
	}

	// Add a test game
	model.GlobalPersistence.Games["testuser"] = model.UserGame{
		User: "testuser",
		PlayedMoves: []model.PlayedMove{
			{Letters: "abc", Words: []string{"test"}},
		},
	}

	// Test successful end
	err = service.EndGame("testuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.SaveCalled {
		t.Error("Expected SaveGamesToFile to be called")
	}

	if _, exists := model.GlobalPersistence.Games["testuser"]; exists {
		t.Error("Game was not removed from active games")
	}

	if len(model.GlobalPersistence.EndedGames) != 1 {
		t.Error("Game was not added to ended games")
	}

	if model.GlobalPersistence.EndedGames[0].User != "testuser" {
		t.Errorf("Expected ended game for user 'testuser', got '%s'",
			model.GlobalPersistence.EndedGames[0].User)
	}

	if model.GlobalPersistence.EndedGames[0].GameEndTimestamp == "" {
		t.Error("Game end timestamp was not set")
	}
}

func TestGetLetters(t *testing.T) {
	service, mock := setupTestEnvironment()

	// Test with non-existent user (should create new game)
	game, err := service.GetLetters("newuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.SaveCalled {
		t.Error("Expected SaveGamesToFile to be called")
	}

	if game.User != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", game.User)
	}

	// Verify game was stored in global state
	if _, exists := model.GlobalPersistence.Games["newuser"]; !exists {
		t.Error("Game was not created in global state")
	}

	// Test with existing user
	mock.SaveCalled = false
	_, err = service.GetLetters("newuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if mock.SaveCalled {
		t.Error("SaveGamesToFile should not be called for existing user")
	}

	// Test with save error
	mock.SaveError = fmt.Errorf("save error")
	_, err = service.GetLetters("anotheruser")

	if err == nil || err.Error() != "failed to save game data: save error" {
		t.Errorf("Expected 'failed to save game data: save error', got %v", err)
	}
}

func TestPlayMove(t *testing.T) {
	service, mock := setupTestEnvironment()

	// Test with non-existent user
	move := model.PlayedMove{
		Letters:        "abc",
		Words:          []string{"test"},
		PlayedByMyself: true,
	}

	_, err := service.PlayMove("nonexistent", move)

	if err == nil {
		t.Error("Expected error for non-existent game, got nil")
	}

	// Add a test game with letters
	model.GlobalPersistence.Games["testuser"] = model.UserGame{
		User:           "testuser",
		LettersPlaySet: logic.LoadLettersPlaySet(),
		PlayedMoves:    []model.PlayedMove{},
	}

	// Test successful move
	updatedGame, err := service.PlayMove("testuser", move)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.SaveCalled {
		t.Error("Expected SaveGamesToFile to be called")
	}

	if len(updatedGame.PlayedMoves) != 1 {
		t.Error("Move was not added to played moves")
	}

	// Test with invalid letters (not available)
	move.Letters = "zzzzzz" // More z's than available
	_, err = service.PlayMove("testuser", move)

	if err == nil {
		t.Error("Expected error for unavailable letters, got nil")
	}

	// Test with save error
	mock.SaveError = fmt.Errorf("save error")
	move.Letters = "a"
	_, err = service.PlayMove("testuser", move)

	if err == nil || err.Error() != "failed to save game data: save error" {
		t.Errorf("Expected 'failed to save game data: save error', got %v", err)
	}
}

func TestListEndedGames(t *testing.T) {
	service, _ := setupTestEnvironment()

	// Test with no ended games
	endedGames := service.ListEndedGames()

	if len(endedGames) != 0 {
		t.Errorf("Expected 0 ended games, got %d", len(endedGames))
	}

	// Add ended games
	model.GlobalPersistence.EndedGames = []model.UserGame{
		{
			User:               "user1",
			LastMoveTimestamp:  "2025-04-21 12:00:00",
			GameStartTimestamp: "2025-04-21 11:00:00",
		},
		{
			User:               "user2",
			LastMoveTimestamp:  "2025-04-21 13:00:00",
			GameStartTimestamp: "2025-04-21 11:30:00",
		},
	}

	// Test with ended games
	endedGames = service.ListEndedGames()

	if len(endedGames) != 2 {
		t.Errorf("Expected 2 ended games, got %d", len(endedGames))
	}

	if endedGames[0].User != "user1" || endedGames[1].User != "user2" {
		t.Error("Incorrect ended games returned")
	}
}

func TestGetPlayedWords(t *testing.T) {
	service, _ := setupTestEnvironment()

	// Test with no words
	words := service.GetPlayedWords()

	if len(words) != 0 {
		t.Errorf("Expected 0 words, got %d", len(words))
	}

	// Add active game with words
	model.GlobalPersistence.Games["user1"] = model.UserGame{
		PlayedMoves: []model.PlayedMove{
			{Words: []string{"hello", "WORLD"}},
			{Words: []string{"hello"}},
		},
	}

	// Add ended game with words
	model.GlobalPersistence.EndedGames = []model.UserGame{
		{
			PlayedMoves: []model.PlayedMove{
				{Words: []string{"test", "Case"}},
			},
		},
	}

	// Test with words
	words = service.GetPlayedWords()

	if len(words) != 4 {
		t.Errorf("Expected 4 unique words, got %d", len(words))
	}

	// Create expected map to verify counts
	expectedCounts := map[string]int{
		"hello": 2,
		"world": 1,
		"test":  1,
		"case":  1,
	}

	// Verify counts
	for _, word := range words {
		expected, exists := expectedCounts[word.Word]
		if !exists {
			t.Errorf("Unexpected word: %s", word.Word)
			continue
		}

		if word.CurrentCount != expected {
			t.Errorf("Expected count %d for word '%s', got %d",
				expected, word.Word, word.CurrentCount)
		}
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name       string
		moves      []model.PlayedMove
		wantCounts map[string]int
	}{
		{
			name:       "empty moves list",
			moves:      []model.PlayedMove{},
			wantCounts: map[string]int{},
		},
		{
			name: "single word",
			moves: []model.PlayedMove{
				{Words: []string{"Hello"}},
			},
			wantCounts: map[string]int{
				"hello": 1,
			},
		},
		{
			name: "multiple words in one move",
			moves: []model.PlayedMove{
				{Words: []string{"Hello", "World"}},
			},
			wantCounts: map[string]int{
				"hello": 1,
				"world": 1,
			},
		},
		{
			name: "duplicate words across moves",
			moves: []model.PlayedMove{
				{Words: []string{"Hello"}},
				{Words: []string{"Hello"}},
			},
			wantCounts: map[string]int{
				"hello": 2,
			},
		},
		{
			name: "nil words field",
			moves: []model.PlayedMove{
				{Words: nil},
			},
			wantCounts: map[string]int{},
		},
		{
			name: "mixed case words",
			moves: []model.PlayedMove{
				{Words: []string{"Hello", "WORLD", "HeLLo"}},
			},
			wantCounts: map[string]int{
				"hello": 2,
				"world": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounts := make(map[string]int)
			countWords(tt.moves, gotCounts)

			if !reflect.DeepEqual(gotCounts, tt.wantCounts) {
				t.Errorf("countWords() = %v, want %v", gotCounts, tt.wantCounts)
			}
		})
	}
}
