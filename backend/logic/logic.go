package logic

import (
	"fmt"
	"net/http"

	"buchstaben.go/model"
)

func LoadLettersPlaySet() (model.LettersPlaySet) {
	lettersPlaySet := model.LettersPlaySet{
		{Letter: "a", Count: 5, Value: 1},
		{Letter: "b", Count: 2, Value: 2},
		{Letter: "c", Count: 2, Value: 4},
		{Letter: "d", Count: 5, Value: 1},
		{Letter: "e", Count: 14, Value: 1},
		{Letter: "f", Count: 2, Value: 4},
		{Letter: "g", Count: 3, Value: 2},
		{Letter: "h", Count: 4, Value: 2},
		{Letter: "i", Count: 6, Value: 1},
		{Letter: "j", Count: 1, Value: 6},
		{Letter: "k", Count: 2, Value: 4},
		{Letter: "l", Count: 3, Value: 2},
		{Letter: "m", Count: 4, Value: 3},
		{Letter: "n", Count: 9, Value: 1},
		{Letter: "o", Count: 3, Value: 2},
		{Letter: "p", Count: 1, Value: 5},
		{Letter: "q", Count: 1, Value: 10},
		{Letter: "r", Count: 6, Value: 1},
		{Letter: "s", Count: 7, Value: 1},
		{Letter: "t", Count: 6, Value: 1},
		{Letter: "u", Count: 6, Value: 1},
		{Letter: "v", Count: 1, Value: 6},
		{Letter: "w", Count: 1, Value: 3},
		{Letter: "x", Count: 1, Value: 8},
		{Letter: "y", Count: 1, Value: 10},
		{Letter: "z", Count: 1, Value: 3},
		{Letter: "ä", Count: 1, Value: 6},
		{Letter: "ö", Count: 1, Value: 8},
		{Letter: "ü", Count: 1, Value: 6},
		{Letter: "*", Count: 2, Value: 0},
	}
	return lettersPlaySet
}

func RemoveLetters(lettersPlaySet model.LettersPlaySet, inputString string) (model.LettersPlaySet, error) {
    previousLettersPlaySet := lettersPlaySet
	for _, letter := range inputString {
		for i, l := range lettersPlaySet {
			if l.Letter == string(letter) {
				if lettersPlaySet[i].Count == 0 {
					fmt.Println("Letter ", l.Letter, " is not available anymore.")
					return previousLettersPlaySet, fmt.Errorf("letter %s is not available anymore, \"Play Move\" ignored", l.Letter)
				}
				lettersPlaySet[i].Count--
				break
			}
		}
	}
	return lettersPlaySet, nil
}

func GetLetterValue(lettersPlaySet model.LettersPlaySet) uint {
	value := uint(0)
	for _, l := range lettersPlaySet {
		value += l.Count * l.Value
	}
	return value
}

func GetRemindingsLetterCount(lettersPlaySet model.LettersPlaySet) uint {
	remindingLetterCount := uint(0)
	for _, l := range lettersPlaySet {
		remindingLetterCount += l.Count
	}
	return remindingLetterCount
}

func GetUserNameFromResponse(r http.Request) model.User {
    return model.User(r.URL.Query().Get("username"))
}