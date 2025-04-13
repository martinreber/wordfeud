import { showMessage } from '../common/utils.js';

interface WordEntry {
    word: string;
    count: number;
}
let allWords: WordEntry[] = [];

function containsAllLetters(word: string, letters: string) {
    const wordChars = word.toLowerCase().split('');
    const searchChars = letters.toLowerCase().split('');
    return searchChars.every(char =>
        wordChars.includes(char)
    );
}


function filterAndDisplayWords(filterText = '') {
    const tableBody = document.querySelector("#words-table tbody") as HTMLElement;
    tableBody.innerHTML = "";

    console.log("All words:", allWords); // Debugging line
    console.log("Filter text:", filterText); // Debugging line
    const wordCount = document.getElementById("word-count") as HTMLElement;
    const filteredWords = filterText
        ? allWords.filter(entry => containsAllLetters(entry.word, filterText))
        : allWords;
        wordCount.textContent = filteredWords.length.toString();

    if (filteredWords) {
        console.log("Filtered words:", filteredWords); // Debugging line
        wordCount.textContent = filteredWords.length.toString();
        filteredWords.forEach(entry => {
            const row = document.createElement("tr");

            const wordCell = document.createElement("td");
            wordCell.textContent = entry.word;

            const countCell = document.createElement("td");
            countCell.textContent = entry.count.toString();

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
            filterAndDisplayWords(data.words);
        })
        .catch((error) => {
            console.error("Error fetching played words:", error);
            showMessage(error.message);
        });
}

document.addEventListener("DOMContentLoaded", () => {
    fetchPlayedWords();

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
        filterInput.value = '';
    });
});