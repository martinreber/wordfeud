import { showMessage } from '../common/utils.js';

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
            const tableBody = document.querySelector("#words-table tbody");
            tableBody.innerHTML = "";

            const wordCount = document.getElementById("word-count");
            wordCount.textContent = data.length;

            data.forEach(entry => {

                const row = document.createElement("tr");

                const wordCell = document.createElement("td");
                wordCell.textContent = entry.word;

                const countCell = document.createElement("td");
                countCell.textContent = entry.count;

                row.appendChild(wordCell);
                row.appendChild(countCell);
                tableBody.appendChild(row);
            });
        })
        .catch((error) => {
            console.error("Error fetching played words:", error);
            showMessage(error.message);
        });
}

document.addEventListener("DOMContentLoaded", () => {
    fetchPlayedWords();

    document.getElementById("refresh-button").addEventListener("click", () => {
        fetchPlayedWords();
    });
});


document.addEventListener("DOMContentLoaded", fetchPlayedWords);