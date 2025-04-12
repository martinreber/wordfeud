import { showMessage } from '../common/utils.js';

function fetchLetters() {
  const username = getUsername();
  fetch(`http://localhost:8080/letters?username=${username}`)
  .then((response) => {
    if (!response.ok) {
      // Check the content type of the response
      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        // Handle JSON error response
        return response.json().then((errorData) => {
          throw new Error(errorData.message || "Failed to play the move.");
        });
      } else {
        // Handle plain text error response
        return response.text().then((errorText) => {
          throw new Error(errorText || "Failed to play the move.");
        });
      }
    }
    return response.json();
  })
  .then((data) => {
      const sessionTableBody = document.querySelector("#session-table tbody");
      const letterContainer = document.getElementById("letters-play-set");
      const overallValue = document.getElementById("overall-value");
      const remainingLetters = document.getElementById("remaining-letters");

      const usernameUi = document.getElementById("user-name");
      usernameUi.textContent = `Username: ${username}`;

      const sessionStart = document.getElementById("session-start-timestamp");
      sessionStart.textContent = `Session Start: ${data.session_start_timestamp}`;

      const lastMoveTimestamp = document.getElementById("last-move-timestamp");
      lastMoveTimestamp.textContent = `Last Move: ${data.last_move_timestamp}`;

      overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

      const totalRemaining = data.letters_play_set.reduce(
        (sum, letter) => sum + letter.count,
        0
      );
      remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

      letterContainer.innerHTML = "";
      letterContainer.appendChild(createLettersTable(data));
    })
    .catch((error) => {
      console.error("Error fetching letter data:", error);
    });
}

document.addEventListener("DOMContentLoaded", () => {
  // Initialize the page
  try {
    const username = getUsername(); // Get the username from the URL
    fetchLetters(); // Fetch and display the data for the username
  } catch (error) {
    console.error(error.message);
    window.location.href = "../index.html"; // Redirect to the sessions list page if no username is provided
  }

  // Add event listeners for buttons
  document
    .getElementById("list-sessions-button")
    .addEventListener("click", () => {
      window.open("../index.html", "_blank");
    });

  document.getElementById("play-move-button").addEventListener("click", () => {
    const username = getUsername();
    const inputString = document.getElementById("input-string").value;
    const inputWord = document.getElementById("input-word").value;
    const isPlayedByMyself = document.getElementById("player-toggle").checked;

    if (!inputString) {
      showMessage("Please enter the letters you played.");
      return;
  }

  if (!inputWord) {
      showMessage("Please enter the word you formed.");
      return;
  }

    fetch(`http://localhost:8080/play-move?username=${username}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        letters: inputString,
        word: inputWord,
        playedByMyself: isPlayedByMyself,
    }),
    })
    .then((response) => {
      if (!response.ok) {
        // Check the content type of the response
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
          // Handle JSON error response
          return response.json().then((errorData) => {
            throw new Error(errorData.message || "Failed to play the move.");
          });
        } else {
          // Handle plain text error response
          return response.text().then((errorText) => {
            throw new Error(errorText || "Failed to play the move.");
          });
        }
      }
      return response.json();
    })
    .then((data) => {
        // Update the UI with the new data
        const letterContainer = document.getElementById("letters-play-set");
        const overallValue = document.getElementById("overall-value");
        const remainingLetters = document.getElementById("remaining-letters");

        const usernameUi = document.getElementById("user-name");
        usernameUi.textContent = `Username: ${username}`;

        const sessionStart = document.getElementById("session-start-timestamp");
        sessionStart.textContent = `Session Start: ${data.session_start_timestamp}`;

        const lastMoveTimestamp = document.getElementById("last-move-timestamp");
        lastMoveTimestamp.textContent = `Last Move: ${data.last_move_timestamp}`;


        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

        const totalRemaining = data.letters_play_set.reduce(
          (sum, letter) => sum + letter.count,
          0
        );
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

        letterContainer.innerHTML = "";
        letterContainer.appendChild(createLettersTable(data));

        const playerToggle = document.getElementById("player-toggle");
        playerToggle.checked = !playerToggle.checked;
      })
    .catch((error) => {
      showMessage(error.message);
    });
  });

  document.getElementById("reset-button").addEventListener("click", () => {
    const username = getUsername();
    fetch(`http://localhost:8080/reset?username=${username}`, {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data) => {
        const letterContainer = document.getElementById("letters-play-set");
        const overallValue = document.getElementById("overall-value");
        const remainingLetters = document.getElementById("remaining-letters");

        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

        const totalRemaining = data.letters_play_set.reduce(
          (sum, letter) => sum + letter.count,
          0
        );
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

        letterContainer.innerHTML = "";
        letterContainer.appendChild(createLettersTable(data));
      })
      .catch((error) => {
        console.error("Error resetting letters:", error);
      });
  });
});

function getUsername() {
  const urlParams = new URLSearchParams(window.location.search);
  const usernameFromQuery = urlParams.get("username");

  if (usernameFromQuery) {
    const usernameInput = document.getElementById("username");
    if (usernameInput) {
      usernameInput.value = usernameFromQuery; // Set the username
    }
    return usernameFromQuery;
  }

  showMessage("No username provided. Please create a session first.");
  throw new Error("Username is required");
}

function createLettersTable(data) {
  const table = document.createElement("table");
  table.classList.add("letter-table"); // Add a class for styling the table

  let row;
  let index = 0;
  data.letters_play_set.forEach((letter) => {
    if (letter.count === 0) {
      // Skip letters with a count of 0
      return;
    }

    if (index % 5 === 0) {
      // Create a new row every 5 letters
      row = document.createElement("tr");
      table.appendChild(row);
    }
    index++;

    // Create a cell for the letter
    const letterCell = document.createElement("td");
    letterCell.textContent = letter.letter;
    letterCell.classList.add("letter-cell"); // Add a class for styling

    // Create a cell for the count
    const countCell = document.createElement("td");
    countCell.textContent = letter.count;
    countCell.classList.add("count-cell"); // Add a class for styling

    // Append both cells to the row
    row.appendChild(letterCell);
    row.appendChild(countCell);
  });
  return table;
}
