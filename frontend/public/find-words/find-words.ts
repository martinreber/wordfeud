import { WordCounts } from '../common/types';
import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';

let allWords: WordCounts = [];

async function findWords(filterText: string): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/find-words?letters=${encodeURIComponent(filterText)}`);
        const filteredWords = await handleResponse<WordCounts>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>("words-table").querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        displayWords(filteredWords, tableBody);
    } catch (error) {
        console.error("Error fetching find words:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
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
            if (e.target.value.length < 2) {

                return;
            }
            findWords(e.target.value); // Pass the filter text to the backend
        }
    });

    const clearButton = getElementByIdOrThrow<HTMLButtonElement>("clear-button");
    clearButton.addEventListener("click", () => {
        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>("words-table").querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        displayWords([], tableBody);
        filterInput.value = '';
        showMessage('');
    });
}

document.addEventListener("DOMContentLoaded", () => {
    setupEventListeners();
});