package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type LetterCount struct {
	Letter string `json:"letter"`
	Count uint `json:"count"`
	count uint `json:"value"`
}

type LettersCount []LetterCount

type LetterResponse struct {
	LetterOverAllValue uint `json:"letter_overall_value"`
	LettersCount LettersCount `json:"letters_count"`
}

var (
    sessions     map[string]LettersCount
    sessionsLock sync.Mutex
)

var (
    letters     LettersCount
    lettersLock sync.Mutex
)

func loadLetters() (LettersCount) {
	lettersCount := LettersCount{
		{Letter: "a", Count: 5, count: 1},
		{Letter: "b", Count: 2, count: 2},
		{Letter: "c", Count: 2, count: 4},
		{Letter: "d", Count: 5, count: 1},
		{Letter: "e", Count: 14, count: 1},
		{Letter: "f", Count: 2, count: 4},
		{Letter: "g", Count: 3, count: 2},
		{Letter: "h", Count: 4, count: 2},
		{Letter: "i", Count: 6, count: 1},
		{Letter: "j", Count: 1, count: 6},
		{Letter: "k", Count: 2, count: 4},
		{Letter: "l", Count: 3, count: 2},
		{Letter: "m", Count: 4, count: 3},
		{Letter: "n", Count: 9, count: 1},
		{Letter: "o", Count: 3, count: 2},
		{Letter: "p", Count: 1, count: 5},
		{Letter: "q", Count: 1, count: 10},
		{Letter: "r", Count: 6, count: 1},
		{Letter: "s", Count: 7, count: 1},
		{Letter: "t", Count: 6, count: 1},
		{Letter: "u", Count: 6, count: 1},
		{Letter: "v", Count: 1, count: 6},
		{Letter: "w", Count: 1, count: 3},
		{Letter: "x", Count: 1, count: 8},
		{Letter: "y", Count: 1, count: 10},
		{Letter: "z", Count: 1, count: 3},
		{Letter: "ä", Count: 1, count: 6},
		{Letter: "ö", Count: 1, count: 8},
		{Letter: "ü", Count: 1, count: 6},
	}
	return lettersCount
}

func removeLetters(lettersCount LettersCount, inputString string) LettersCount {
	for _, letter := range inputString {
		for i, l := range lettersCount {
			if l.Letter == string(letter) {
				if lettersCount[i].Count == 0 {
					fmt.Println("Letter ", l.Letter, " is not available anymore.")
					break
				}
				lettersCount[i].Count--
				break
			}
		}
	}
	return lettersCount
}

func getValue(LettersCount LettersCount) uint {
	value := uint(0)
	for _, l := range LettersCount {
		value += l.Count * l.count
	}
	return value
}

func getRemindingLetterCount(lettersCount LettersCount) uint {
	value := uint(0)
	for _, l := range lettersCount {
		value += l.Count
	}
	return value
}

func getLettersHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    sessionsLock.Lock()
    defer sessionsLock.Unlock()

    letters, exists := sessions[username]
    if !exists {
        letters = loadLetters() // Initialize letters for the new session
        sessions[username] = letters
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        LetterResponse{
            LetterOverAllValue: getValue(letters),
            LettersCount:       letters,
        },
    )
}

func processInputHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    var input struct {
        String string `json:"string"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    sessionsLock.Lock()
    defer sessionsLock.Unlock()

    letters, exists := sessions[username]
    if !exists {
        http.Error(w, "Session not found for username", http.StatusNotFound)
        return
    }

    letters = removeLetters(letters, input.String)
    sessions[username] = letters

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        LetterResponse{
            LetterOverAllValue: getValue(letters),
            LettersCount:       letters,
        },
    )
}

func resetLettersHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    sessionsLock.Lock()
    defer sessionsLock.Unlock()

    sessions[username] = loadLetters()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(
        LetterResponse{
            LetterOverAllValue: getValue(sessions[username]),
            LettersCount:       sessions[username],
        },
    )
}

func deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete session")
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	sessionsLock.Lock()
	defer sessionsLock.Unlock()

    w.Header().Set("Content-Type", "application/json")

	if _, exists := sessions[username]; !exists {
		http.Error(w, "Session not found for username", http.StatusNotFound)
		return
	}

	delete(sessions, username)
	w.WriteHeader(http.StatusOK)
}

type ListSessionsResponse struct {
	Sessions []ListSession `json:"sessions"`
}

type ListSession struct {
	User string `json:"user"`
	RemindingLetters uint `json:"reminding_letters"`
}

func listSessionsHandler(w http.ResponseWriter, r *http.Request) {

    sessionsLock.Lock()
    defer sessionsLock.Unlock()

	listSessions := ListSessionsResponse{}
	listSession := ListSession{}

	for user := range sessions {
		listSession.User = user
		listSession.RemindingLetters = getRemindingLetterCount(sessions[user])
		listSessions.Sessions = append(listSessions.Sessions, listSession)
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(listSessions)
}

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    sessionsLock.Lock()
    defer sessionsLock.Unlock()

    if _, exists := sessions[username]; exists {
        http.Error(w, "Session already exists for this username" , http.StatusBadRequest)
        return
    }

    sessions[username] = loadLetters() // Initialize a new session
    w.WriteHeader(http.StatusCreated)
}
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    letters = loadLetters()
    sessions = make(map[string]LettersCount)

    http.Handle("/letters", enableCORS(http.HandlerFunc(getLettersHandler)))
    http.Handle("/process", enableCORS(http.HandlerFunc(processInputHandler)))
    http.Handle("/reset", enableCORS(http.HandlerFunc(resetLettersHandler)))
	http.Handle("/list", enableCORS(http.HandlerFunc(listSessionsHandler)))
	http.Handle("/create", enableCORS(http.HandlerFunc(createSessionHandler)))
	http.Handle("/delete", enableCORS(http.HandlerFunc(deleteSessionHandler)))

    fmt.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}