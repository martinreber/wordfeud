import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';
import { WordCount } from '../common/types.js';

let allWords: WordCount[] = [];

function containsAllLetters(word: string, letters: string): boolean {
    const wordChars = word.toLowerCase().split('');
    const searchChars = letters.toLowerCase().split('');
    return searchChars.every(char => wordChars.includes(char));
}

function filterAndDisplayWords(filterText = ''): void {
    const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>("words-table").querySelector('tbody');
    if (!tableBody) throw new Error('Table body not found');

    tableBody.innerHTML = "";

    const filteredWords = filterText
        ? allWords.filter(entry => containsAllLetters(entry.word, filterText))
        : allWords;

    updateWordCount(filteredWords.length);
    displayWords(filteredWords, tableBody);
}

function updateWordCount(count: number): void {
    const wordCount = getElementByIdOrThrow<HTMLElement>("word-count");
    wordCount.textContent = count.toString();
}

function displayWords(words: WordCount[], container: HTMLElement): void {
    words.forEach(entry => {
        const row = document.createElement("tr");
        row.appendChild(createCell(entry.word));
        row.appendChild(createCell(entry.count.toString()));
        container.appendChild(row);
    });
}

function createCell(content: string): HTMLTableCellElement {
    const cell = document.createElement("td");
    cell.textContent = content;
    return cell;
}

async function fetchPlayedWords(): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/played-words`);
        const data = await handleResponse<{ words: WordCount[] }>(response);
        allWords = data.words;
        filterAndDisplayWords();
    } catch (error) {
        console.error("Error fetching played words:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
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