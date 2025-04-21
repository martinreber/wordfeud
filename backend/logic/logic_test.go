package logic

import (
	"testing"

	"buchstaben.go/model"
	"github.com/stretchr/testify/assert"
)

func TestLoadLettersPlaySet(t *testing.T) {
	result := LoadLettersPlaySet()

	// Check that we get all expected letters
	assert.Len(t, result, 30, "Should have 30 different letters in the play set")

	// Verify a few specific letters
	letterMap := make(map[string]model.LetterPlaySet)
	for _, l := range result {
		letterMap[l.Letter] = l
	}

	// Check 'a'
	a, exists := letterMap["a"]
	assert.True(t, exists, "Letter 'a' should exist")
	assert.Equal(t, uint(5), a.OriginalCount, "Letter 'a' should have original count of 5")
	assert.Equal(t, uint(5), a.CurrentCount, "Letter 'a' should have current count of 5")
	assert.Equal(t, uint(1), a.Value, "Letter 'a' should have value of 1")

	// Check 'e' - highest frequency
	e, exists := letterMap["e"]
	assert.True(t, exists, "Letter 'e' should exist")
	assert.Equal(t, uint(14), e.OriginalCount, "Letter 'e' should have original count of 14")
	assert.Equal(t, uint(1), e.Value, "Letter 'e' should have value of 1")

	// Check 'q' - high value
	q, exists := letterMap["q"]
	assert.True(t, exists, "Letter 'q' should exist")
	assert.Equal(t, uint(10), q.Value, "Letter 'q' should have value of 10")

	// Check special German letters
	_, exists = letterMap["ä"]
	assert.True(t, exists, "Letter 'ä' should exist")
	_, exists = letterMap["ö"]
	assert.True(t, exists, "Letter 'ö' should exist")
	_, exists = letterMap["ü"]
	assert.True(t, exists, "Letter 'ü' should exist")

	// Check wildcard
	wildcard, exists := letterMap["*"]
	assert.True(t, exists, "Wildcard '*' should exist")
	assert.Equal(t, uint(0), wildcard.Value, "Wildcard '*' should have value of 0")
	assert.Equal(t, uint(2), wildcard.OriginalCount, "Wildcard '*' should have original count of 2")
}

func TestRemoveLetters(t *testing.T) {
	testCases := []struct {
		name          string
		inputSet      model.LettersPlaySet
		inputString   string
		expectedError bool
		errorContains string
	}{
		{
			name: "Remove valid single letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 5, Value: 1},
				{Letter: "b", OriginalCount: 2, CurrentCount: 2, Value: 2},
			},
			inputString:   "a",
			expectedError: false,
		},
		{
			name: "Remove valid multiple letters",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 5, Value: 1},
				{Letter: "b", OriginalCount: 2, CurrentCount: 2, Value: 2},
			},
			inputString:   "ab",
			expectedError: false,
		},
		{
			name: "Remove unavailable letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 0, Value: 1},
				{Letter: "b", OriginalCount: 2, CurrentCount: 2, Value: 2},
			},
			inputString:   "a",
			expectedError: true,
			errorContains: "not available",
		},
		{
			name: "Remove invalid letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 5, Value: 1},
				{Letter: "b", OriginalCount: 2, CurrentCount: 2, Value: 2},
			},
			inputString:   "c",
			expectedError: true,
			errorContains: "not valid",
		},
		{
			name: "Remove German letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 5, Value: 1},
				{Letter: "ö", OriginalCount: 1, CurrentCount: 1, Value: 8},
			},
			inputString:   "ö",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalSetCopy := make(model.LettersPlaySet, len(tc.inputSet))
			copy(originalSetCopy, tc.inputSet)

			result, err := RemoveLetters(tc.inputSet, tc.inputString)

			if tc.expectedError {
				assert.Error(t, err, "Should return an error")
				assert.Contains(t, err.Error(), tc.errorContains, "Error should contain expected text")

				// When error occurs, the original set should be returned unchanged
				assert.Equal(t, originalSetCopy, result, "When error occurs, original set should be returned")
			} else {
				assert.NoError(t, err, "Should not return an error")

				// For each letter in the input string, its count should decrease by 1
				letterCounts := make(map[string]int)
				for _, letter := range tc.inputString {
					letterCounts[string(letter)]++
				}

				for _, letterSet := range result {
					originalLetter := findLetterInSet(originalSetCopy, letterSet.Letter)
					assert.NotNil(t, originalLetter, "Letter should exist in original set")

					expectedCount := originalLetter.CurrentCount
					if count, exists := letterCounts[letterSet.Letter]; exists {
						expectedCount -= uint(count)
					}

					assert.Equal(t, expectedCount, letterSet.CurrentCount,
						"Current count for letter %s should decrease by %d",
						letterSet.Letter, letterCounts[letterSet.Letter])
				}
			}
		})
	}
}

func TestGetLetterValue(t *testing.T) {
	testCases := []struct {
		name          string
		inputSet      model.LettersPlaySet
		expectedValue uint
	}{
		{
			name:          "Empty set",
			inputSet:      model.LettersPlaySet{},
			expectedValue: 0,
		},
		{
			name: "Single letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1},
			},
			expectedValue: 3, // 3 * 1 = 3
		},
		{
			name: "Multiple letters",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1},
				{Letter: "q", OriginalCount: 1, CurrentCount: 1, Value: 10},
				{Letter: "z", OriginalCount: 1, CurrentCount: 1, Value: 3},
			},
			expectedValue: 16, // (3 * 1) + (1 * 10) + (1 * 3) = 16
		},
		{
			name: "With wildcard",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1},
				{Letter: "*", OriginalCount: 2, CurrentCount: 2, Value: 0},
			},
			expectedValue: 3, // (3 * 1) + (2 * 0) = 3
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetLetterValue(tc.inputSet)
			assert.Equal(t, tc.expectedValue, result,
				"Letter value calculation should match expected value")
		})
	}
}

func TestGetRemindingsLetterCount(t *testing.T) {
	testCases := []struct {
		name          string
		inputSet      model.LettersPlaySet
		expectedCount uint
	}{
		{
			name:          "Empty set",
			inputSet:      model.LettersPlaySet{},
			expectedCount: 0,
		},
		{
			name: "Single letter",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1},
			},
			expectedCount: 3,
		},
		{
			name: "Multiple letters",
			inputSet: model.LettersPlaySet{
				{Letter: "a", OriginalCount: 5, CurrentCount: 3, Value: 1},
				{Letter: "b", OriginalCount: 2, CurrentCount: 1, Value: 2},
				{Letter: "c", OriginalCount: 2, CurrentCount: 0, Value: 4},
			},
			expectedCount: 4, // 3 + 1 + 0 = 4
		},
		{
			name:          "Full set",
			inputSet:      LoadLettersPlaySet(),
			expectedCount: 102, // Sum of all current counts in default set
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetRemindingsLetterCount(tc.inputSet)
			assert.Equal(t, tc.expectedCount, result,
				"Remaining letter count should match expected count")
		})
	}
}

// Helper function to find a letter in a set
func findLetterInSet(set model.LettersPlaySet, letter string) *model.LetterPlaySet {
	for i, l := range set {
		if l.Letter == letter {
			return &set[i]
		}
	}
	return nil
}
