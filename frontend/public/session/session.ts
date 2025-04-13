import { showMessage } from '../common/utils.js';

function fetchLetters()
{
  const username = getUsername();
  fetch(`http://localhost:8080/letters?username=${username}`)
    .then((response) =>
    {
      if (!response.ok) {
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
          return response.json().then((errorData) =>
          {
            throw new Error(errorData.message || "Failed to play the move.");
          });
        } else {
          return response.text().then((errorText) =>
          {
            throw new Error(errorText || "Failed to play the move.");
          });
        }
      }
      return response.json();
    })
    .then((data) =>
    {
      const sessionTableBody = document.querySelector("#session-table tbody");
      const remainingLetters = document.getElementById("remaining-letters");

      const usernameUi = document.getElementById("user-name");
      if (usernameUi) {
        usernameUi.textContent = `Username: ${username}`;
      }

      const sessionStart = document.getElementById("session-start-timestamp");
      if (sessionStart) {
        sessionStart.textContent = `Session Start: ${data.session_start_timestamp}`;
      }


      const lastMoveTimestamp = document.getElementById("last-move-timestamp");
      if (lastMoveTimestamp) {
        lastMoveTimestamp.textContent = `Last Move: ${data.last_move_timestamp}`;
      }

      const overallValue = document.getElementById("overall-value");
      if (overallValue) {
        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;
      }

      const totalRemaining = data.letters_play_set.reduce(
        (sum, letter) => sum + letter.count,
        0
      );

      const letterContainer = document.getElementById("letters-play-set");
      if (remainingLetters) {
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;
      }
      if (letterContainer) {
        letterContainer.innerHTML = "";
        letterContainer.appendChild(createLettersTable(data));
      }
    })
    .catch((error) =>
    {
      console.error("Error fetching letter data:", error);
    });
}

document.addEventListener("DOMContentLoaded", () =>
{
  // Initialize the page
  try {
    const username = getUsername(); // Get the username from the URL
    fetchLetters(); // Fetch and display the data for the username
  } catch (error) {
    console.error(error.message);
    window.location.href = "../index.html"; // Redirect to the sessions list page if no username is provided
  }

  const listSessionButton = document.getElementById("list-sessions-button");
  if (listSessionButton) {
    listSessionButton.addEventListener("click", () => { window.open("../index.html", "_blank"); });
  }

  const playMoveButton = document.getElementById("play-move-button");
  if (!playMoveButton) {
    console.error("Play Move button not found");
    return;
  }
  playMoveButton.addEventListener("click", () =>
  {
    const username = getUsername();
    const inputStringElement = document.getElementById("input-string") as HTMLInputElement;
    const inputString = inputStringElement ? inputStringElement.value : "";
    const inputWordElement = document.getElementById("input-word") as HTMLInputElement;
    const inputWord = inputWordElement ? inputWordElement.value : "";
    const playerToggleElement = document.getElementById("player-toggle") as HTMLInputElement;
    const isPlayedByMyself = playerToggleElement ? playerToggleElement.checked : false;

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
      .then((response) =>
      {
        if (!response.ok) {
          // Check the content type of the response
          const contentType = response.headers.get("content-type");
          if (contentType && contentType.includes("application/json")) {
            // Handle JSON error response
            return response.json().then((errorData) =>
            {
              throw new Error(errorData.message || "Failed to play the move.");
            });
          } else {
            // Handle plain text error response
            return response.text().then((errorText) =>
            {
              throw new Error(errorText || "Failed to play the move.");
            });
          }
        }
        return response.json();
      })
      .then((data) =>
      {
        // Update the UI with the new data
        const letterContainer = document.getElementById("letters-play-set");

        const usernameUi = document.getElementById("user-name") as HTMLInputElement;
        usernameUi.textContent = `Username: ${username}`;

        const sessionStart = document.getElementById("session-start-timestamp") as HTMLInputElement;
        sessionStart.textContent = `Session Start: ${data.session_start_timestamp}`;

        const lastMoveTimestamp = document.getElementById("last-move-timestamp") as HTMLInputElement;
        lastMoveTimestamp.textContent = `Last Move: ${data.last_move_timestamp}`;


        const overallValue = document.getElementById("overall-value") as HTMLInputElement;
        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

        const totalRemaining = data.letters_play_set.reduce(
          (sum, letter) => sum + letter.count,
          0
        );
        const remainingLetters = document.getElementById("remaining-letters") as HTMLInputElement;
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

        if (letterContainer) {
          letterContainer.innerHTML = "";
          letterContainer.appendChild(createLettersTable(data));
        }

        const playerToggle = document.getElementById("player-toggle") as HTMLInputElement;
        if (playerToggle) {
          playerToggle.checked = !playerToggle.checked;
        }
      })
      .catch((error) =>
      {
        showMessage(error.message);
      });
  });

  const resetButton = document.getElementById("reset-button");
  if (!resetButton) {
    console.error("Reset button not found");
    return;
  }
  resetButton.addEventListener("click", () =>
  {
    const username = getUsername();
    fetch(`http://localhost:8080/reset?username=${username}`, {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data) =>
      {
        const letterContainer = document.getElementById("letters-play-set");

        const overallValue = document.getElementById("overall-value") as HTMLInputElement;
        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

        const totalRemaining = data.letters_play_set.reduce(
          (sum, letter) => sum + letter.count,
          0
        );
        const remainingLetters = document.getElementById("remaining-letters") as HTMLInputElement;
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

        if (letterContainer) {
          letterContainer.innerHTML = "";
          letterContainer.appendChild(createLettersTable(data));
        }
        const usernameUi = document.getElementById("user-name") as HTMLInputElement;
        usernameUi.textContent = `Username: ${username}`;
        if (letterContainer) {
          letterContainer.innerHTML = "";
          letterContainer.appendChild(createLettersTable(data));
        }
      })
      .catch((error) =>
      {
        console.error("Error resetting letters:", error);
      });
  });
});

function getUsername()
{
  const urlParams = new URLSearchParams(window.location.search);
  const usernameFromQuery = urlParams.get("username");

  if (usernameFromQuery) {
    const usernameInput = document.getElementById("username") as HTMLInputElement;
    if (usernameInput) {
      usernameInput.value = usernameFromQuery; // Set the username
    }
    return usernameFromQuery;
  }

  showMessage("No username provided. Please create a session first.");
  throw new Error("Username is required");
}

function createLettersTable(data)
{
  const table = document.createElement("table");
  table.classList.add("letter-table"); // Add a class for styling the table

  let row;
  let index = 0;
  data.letters_play_set.forEach((letter) =>
  {
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
