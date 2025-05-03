import { WordCounts } from './../common/types';
import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';

let allWords: WordCounts = [];

async function fetchPlayedWords(filterText: string): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/played-words?filter=${encodeURIComponent(filterText)}`);
        const filteredWords = await handleResponse<WordCounts>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>("words-table").querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        updateWordCount(filteredWords.length);
        displayWords(filteredWords, tableBody);
    } catch (error) {
        console.error("Error fetching played words:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
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
            fetchPlayedWords(e.target.value); // Pass the filter text to the backend
        }
    });

    const refreshButton = getElementByIdOrThrow<HTMLButtonElement>("refresh-button");
    refreshButton.addEventListener("click", () => {
        fetchPlayedWords(''); // Clear the filter and fetch all words
        filterInput.value = '';
    });
}

document.addEventListener("DOMContentLoaded", () => {
    setupEventListeners();
    fetchPlayedWords(''); // Initial fetch with no filter
});