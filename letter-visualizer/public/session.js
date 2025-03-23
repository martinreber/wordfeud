function fetchLetters() {
  const username = getUsername();
  fetch(`http://localhost:8080/letters?username=${username}`)
      .then((response) => response.json())
      .then((data) => {
          const letterContainer = document.getElementById("letters-count");
          const overallValue = document.getElementById("overall-value");
          const remainingLetters = document.getElementById("remaining-letters");

          overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

          const totalRemaining = data.letters_count.reduce((sum, letter) => sum + letter.count, 0);
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
  document.getElementById("list-sessions-button").addEventListener("click", () => {
      window.open("index.html", "_blank");
  });

  document.getElementById("process-button").addEventListener("click", () => {
      const username = getUsername();
      const inputString = document.getElementById("input-string").value;

      if (!inputString) {
          alert("Please enter a string to process.");
          return;
      }

      fetch(`http://localhost:8080/process?username=${username}`, {
          method: "POST",
          headers: {
              "Content-Type": "application/json",
          },
          body: JSON.stringify({ string: inputString }),
      })
          .then((response) => response.json())
          .then((data) => {
              const letterContainer = document.getElementById("letters-count");
              const overallValue = document.getElementById("overall-value");

              overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

              letterContainer.innerHTML = "";
              letterContainer.appendChild(createLettersTable(data));
          })
          .catch((error) => {
              console.error("Error processing input string:", error);
          });
  });

  document.getElementById("reset-button").addEventListener("click", () => {
      const username = getUsername();
      fetch(`http://localhost:8080/reset?username=${username}`, {
          method: "POST",
      })
          .then((response) => response.json())
          .then((data) => {
              const letterContainer = document.getElementById("letters-count");
              const overallValue = document.getElementById("overall-value");

              overallValue.textContent = `Overall Letter Value: ${data.letter_overall_value}`;

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

  alert("No username provided. Please create a session first.");
  throw new Error("Username is required");
}

function createLettersTable(data) {
  const table = document.createElement("table");
  table.classList.add("letter-table"); // Add a class for styling the table

  let row;
  let index = 0;
  data.letters_count.forEach((letter) => {
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