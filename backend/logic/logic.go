package logic

import (
	"fmt"

	"buchstaben.go/model"
)

func LoadLettersPlaySet() model.LettersPlaySet {
	lettersPlaySet := model.LettersPlaySet{
		{Letter: "a", OriginalCount: 5, CurrentCount: 5, Value: 1},
		{Letter: "b", OriginalCount: 2, CurrentCount: 2, Value: 2},
		{Letter: "c", OriginalCount: 2, CurrentCount: 2, Value: 4},
		{Letter: "d", OriginalCount: 5, CurrentCount: 5, Value: 1},
		{Letter: "e", OriginalCount: 14, CurrentCount: 14, Value: 1},
		{Letter: "f", OriginalCount: 2, CurrentCount: 2, Value: 4},
		{Letter: "g", OriginalCount: 3, CurrentCount: 3, Value: 2},
		{Letter: "h", OriginalCount: 4, CurrentCount: 4, Value: 2},
		{Letter: "i", OriginalCount: 6, CurrentCount: 6, Value: 1},
		{Letter: "j", OriginalCount: 1, CurrentCount: 1, Value: 6},
		{Letter: "k", OriginalCount: 2, CurrentCount: 2, Value: 4},
		{Letter: "l", OriginalCount: 3, CurrentCount: 3, Value: 2},
		{Letter: "m", OriginalCount: 4, CurrentCount: 4, Value: 3},
		{Letter: "n", OriginalCount: 9, CurrentCount: 9, Value: 1},
		{Letter: "o", OriginalCount: 3, CurrentCount: 3, Value: 2},
		{Letter: "p", OriginalCount: 1, CurrentCount: 1, Value: 5},
		{Letter: "q", OriginalCount: 1, CurrentCount: 1, Value: 10},
		{Letter: "r", OriginalCount: 6, CurrentCount: 6, Value: 1},
		{Letter: "s", OriginalCount: 7, CurrentCount: 7, Value: 1},
		{Letter: "t", OriginalCount: 6, CurrentCount: 6, Value: 1},
		{Letter: "u", OriginalCount: 6, CurrentCount: 6, Value: 1},
		{Letter: "v", OriginalCount: 1, CurrentCount: 1, Value: 6},
		{Letter: "w", OriginalCount: 1, CurrentCount: 1, Value: 3},
		{Letter: "x", OriginalCount: 1, CurrentCount: 1, Value: 8},
		{Letter: "y", OriginalCount: 1, CurrentCount: 1, Value: 10},
		{Letter: "z", OriginalCount: 1, CurrentCount: 1, Value: 3},
		{Letter: "ä", OriginalCount: 1, CurrentCount: 1, Value: 6},
		{Letter: "ö", OriginalCount: 1, CurrentCount: 1, Value: 8},
		{Letter: "ü", OriginalCount: 1, CurrentCount: 1, Value: 6},
		{Letter: "*", OriginalCount: 2, CurrentCount: 2, Value: 0},
	}
	return lettersPlaySet
}

func RemoveLetters(lettersPlaySet model.LettersPlaySet, inputString string) (model.LettersPlaySet, error) {
	previousLettersPlaySet := lettersPlaySet
	for _, letter := range inputString {
		isValidLetter := false
		for i, l := range lettersPlaySet {
			if l.Letter == string(letter) {
				if lettersPlaySet[i].CurrentCount == 0 {
					fmt.Println("Letter ", l.Letter, " is not available anymore.")
					return previousLettersPlaySet, fmt.Errorf("letter %q is not available anymore, \"Play Move\" ignored", l.Letter)
				}
				lettersPlaySet[i].CurrentCount--
				isValidLetter = true
				break
			}
		}
		if !isValidLetter {
			fmt.Println("Letter ", string(letter), " is not valid.")
			return previousLettersPlaySet, fmt.Errorf("letter %q is not valid, \"Play Move\" ignored", string(letter))
		}
	}
	return lettersPlaySet, nil
}

func GetLetterValue(lettersPlaySet model.LettersPlaySet) uint {
	value := uint(0)
	for _, l := range lettersPlaySet {
		value += l.CurrentCount * l.Value
	}
	return value
}

func GetRemindingsLetterCount(lettersPlaySet model.LettersPlaySet) uint {
	remindingLetterCount := uint(0)
	for _, l := range lettersPlaySet {
		remindingLetterCount += l.CurrentCount
	}
	return remindingLetterCount
}
