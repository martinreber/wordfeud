function fetchLetters() {
  const username = getUsername();
  fetch(`http://localhost:8080/letters?username=${username}`)
    .then((response) => response.json())
    .then((data) => {
      const sessionTableBody = document.querySelector("#session-table tbody");
      const letterContainer = document.getElementById("letters-play-set");
      const overallValue = document.getElementById("overall-value");
      const remainingLetters = document.getElementById("remaining-letters");

      sessionTableBody.innerHTML = "";
      const row = document.createElement("tr");
      const usernameCell = document.createElement("td");
      usernameCell.textContent = username;

      const sessionStartCell = document.createElement("td");
      sessionStartCell.textContent = data.session_start_timestamp;

      const timestampCell = document.createElement("td");
      timestampCell.textContent = data.last_move_timestamp;

      row.appendChild(usernameCell);
      row.appendChild(sessionStartCell);
      row.appendChild(timestampCell);
      sessionTableBody.appendChild(row);

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
    window.location.href = "index.html"; // Redirect to the sessions list page if no username is provided
  }

  // Add event listeners for buttons
  document
    .getElementById("list-sessions-button")
    .addEventListener("click", () => {
      window.open("index.html", "_blank");
    });

  document.getElementById("play-move-button").addEventListener("click", () => {
    const username = getUsername();
    const inputString = document.getElementById("input-string").value;

    if (!inputString) {
      showMessage("Please enter a string to process.");
      return;
    }

    fetch(`http://localhost:8080/play-move?username=${username}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ string: inputString }),
    })
      .then((response) => {
        if (!response.ok) {
          // If the response is not OK, throw an error with the response message
          return response.json().then((data) => {
            throw new Error(data.message || "Failed to play the move.");
          });
        }
        return response.json(); // Parse the JSON if the response is OK
      })
      .then((data) => {
        // Update the UI with the new data
        const letterContainer = document.getElementById("letters-play-set");
        const overallValue = document.getElementById("overall-value");
        const remainingLetters = document.getElementById("remaining-letters");
        const sessionTableBody = document.querySelector("#session-table tbody");

        // Update the session table with the new timestamp
        sessionTableBody.innerHTML = "";
        const row = document.createElement("tr");
        const usernameCell = document.createElement("td");
        usernameCell.textContent = username;
        const timestampCell = document.createElement("td");
        timestampCell.textContent = data.last_move_timestamp;

        row.appendChild(usernameCell);
        row.appendChild(timestampCell);
        sessionTableBody.appendChild(row);

        // Update the overall letter value
        overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

        // Calculate and display the total remaining letters
        const totalRemaining = data.letters_play_set.reduce(
          (sum, letter) => sum + letter.count,
          0
        );
        remainingLetters.textContent = `Remaining Letters: ${totalRemaining}`;

        // Update the letter table
        letterContainer.innerHTML = "";
        letterContainer.appendChild(createLettersTable(data));
      })
      .catch((error) => {
        // Show the error message in the UI
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

function showMessage(message) {
  const messageContainer = document.getElementById("response-message");
  messageContainer.textContent = message; // Set the message text
}
