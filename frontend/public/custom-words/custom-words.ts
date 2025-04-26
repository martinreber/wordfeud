import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';
import { CustomWord } from '../common/types.js';

// Event Listeners
document.addEventListener("DOMContentLoaded", () => {
    const addWordButton = getElementByIdOrThrow<HTMLButtonElement>("add-word-button");
    addWordButton.addEventListener("click", async () => {
        const newWordInput = getElementByIdOrThrow<HTMLInputElement>("new-words");
        const newWord = newWordInput.value.trim();

        if (!newWord) {
            showMessage("Please enter a word.");
            return;
        }

        await addCustomWord(newWord);
        newWordInput.value = "";
    });

    // Initial fetch
    fetchCustomWords();
});

async function fetchCustomWords(): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/custom-words`);
        const data = await handleResponse<CustomWord[]>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>('custom-words-table').querySelector('tbody');
        if (!tableBody) throw new Error('Custom words table body not found');

        tableBody.innerHTML = ""; // Clear existing rows

        // Sort custom words alphabetically
        if (data) {
            data.sort((a, b) => a.word.localeCompare(b.word));
            data.forEach((customWord) => {
                const row = createCustomWordRow(customWord);
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error("Error fetching custom words:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

function createCustomWordRow(customWord: CustomWord): HTMLTableRowElement {
    const row = document.createElement("tr");

    // Word cell
    const wordCell = document.createElement("td");
    wordCell.textContent = customWord.word;
    row.appendChild(wordCell);

    // Timestamp cell
    const timestampCell = document.createElement("td");
    timestampCell.textContent = customWord.timestamp;
    row.appendChild(timestampCell);

    // Delete button cell
    const actionCell = document.createElement("td");
    const deleteButton = document.createElement("button");
    deleteButton.textContent = "Delete";
    deleteButton.classList.add("delete-word-button");
    deleteButton.addEventListener("click", () => deleteCustomWord(customWord.word));
    actionCell.appendChild(deleteButton);
    row.appendChild(actionCell);

    return row;
}

async function addCustomWord(word: string): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/custom-words`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                word: word,
            }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to add word');
        }

        showMessage(`Word "${word}" added successfully.`);
        await fetchCustomWords();
    } catch (error) {
        console.error("Error adding custom word:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

async function deleteCustomWord(word: string): Promise<void> {
    if (!confirm(`Are you sure you want to delete the word "${word}"?`)) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/custom-words/${encodeURIComponent(word)}`, {
            method: "DELETE",
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || "Failed to delete word");
        }

        showMessage(`Word "${word}" deleted successfully.`);
        await fetchCustomWords();
    } catch (error) {
        console.error("Error deleting word:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}