import { showMessage } from '../common/utils.js';

// let allWords = []; // Store all words for filtering

function containsAllLetters(word, letters) {
    const wordChars = word.toLowerCase().split('');
    const searchChars = letters.toLowerCase().split('');
    return searchChars.every(char =>
        wordChars.includes(char)
    );
}

function filterAndDisplayWords(allWords, filterText = '') {
    const tableBody = document.querySelector("#words-table tbody") as HTMLElement;
    tableBody.innerHTML = "";

    console.log("All words:", allWords); // Debugging line
    console.log("Filter text:", filterText); // Debugging line
    const filteredWords = filterText
        ? allWords.filter(entry => containsAllLetters(entry.word, filterText))
        : allWords;

    const wordCount = document.getElementById("word-count") as HTMLElement;

    if (filteredWords) {
        console.log("Filtered words:", filteredWords); // Debugging line
        wordCount.textContent = filteredWords.length;
        filteredWords.forEach(entry => {
            const row = document.createElement("tr");

            const wordCell = document.createElement("td");
            wordCell.textContent = entry.word;

            const countCell = document.createElement("td");
            countCell.textContent = entry.count;

            row.appendChild(wordCell);
            row.appendChild(countCell);
            tableBody.appendChild(row);
        });
    }
}

function fetchPlayedWords() {
    fetch("http://localhost:8080/played-words")
        .then((response) => {
            if (!response.ok) {
                return response.json().then((errorData) => {
                    throw new Error(errorData.message || "Failed to fetch played words.");
                });
            }
            return response.json();
        })
        .then((data) => {
            console.log("Fetched played words:", data); // Debugging line
            filterAndDisplayWords(data.words); // Initial display
        })
        .catch((error) => {
            console.error("Error fetching played words:", error);
            showMessage(error.message);
        });
}

document.addEventListener("DOMContentLoaded", () => {
    fetchPlayedWords();

    // Add filter input handler
    const filterInput = document.getElementById("word-filter") as HTMLInputElement;
    filterInput.addEventListener("input", (e) => {
        if (e.target instanceof HTMLInputElement) {
            filterAndDisplayWords(e.target.value);
        }
    });

    const refreshButton = document.getElementById("refresh-button");
    if (!refreshButton) {
        console.error("Refresh button not found");
        return;
    }
    refreshButton.addEventListener("click", () => {
        fetchPlayedWords();
        filterInput.value = ''; // Clear filter on refresh
    });
});