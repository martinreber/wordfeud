import { WordCounts } from './../common/types';
import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';

let allWords: WordCounts = [];

async function fetchPlayedWords(): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/played-words`);
        allWords = await handleResponse<WordCounts>(response);
        filterAndDisplayWords('');
    } catch (error) {
        console.error("Error fetching played words:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
        allWords = [];
    }
}

function filterAndDisplayWords(filterText: string): void {
    const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>("words-table").querySelector('tbody');
    if (!tableBody) throw new Error('Table body not found');

    const filteredWords = filterText.trim() === ''
        ? allWords
        : allWords.filter(entry => containsAllLetters(entry.word, filterText));

    updateWordCount(filteredWords.length);
    displayWords(filteredWords, tableBody);
}

function containsAllLetters(word: string, letters: string): boolean {
    const wordChars = word.toLowerCase().split('');
    const searchChars = letters.toLowerCase().split('');

    return searchChars.every(char => {
        const included = wordChars.includes(char);
        return included;
    });
}


function updateWordCount(count: number): void {
    const wordCount = getElementByIdOrThrow<HTMLElement>("word-count");
    wordCount.textContent = count.toString();
}

function displayWords(words: WordCounts, container: HTMLElement): void {
    container.innerHTML = '';

    words.forEach(entry => {
        const row = document.createElement("tr");
        row.appendChild(createCell(entry.word));
        row.appendChild(createCell(entry.current_count.toString()));
        container.appendChild(row);
    });
}

function createCell(content: string): HTMLTableCellElement {
    const cell = document.createElement("td");
    cell.textContent = content;
    return cell;
}

function setupEventListeners(): void {
    const filterInput = getElementByIdOrThrow<HTMLInputElement>("word-filter");
    filterInput.addEventListener("input", (e) => {
        if (e.target instanceof HTMLInputElement) {
            filterAndDisplayWords(e.target.value);
        }
    });

    const refreshButton = getElementByIdOrThrow<HTMLButtonElement>("refresh-button");
    refreshButton.addEventListener("click", () => {
        fetchPlayedWords();
        filterInput.value = '';
    });
}

document.addEventListener("DOMContentLoaded", () => {
    setupEventListeners();
    fetchPlayedWords();
});