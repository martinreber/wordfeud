import
{
  showMessage,
  handleResponse,
  getElementByIdOrThrow,
  updateTextContent,
  API_BASE_URL
} from '../common/utils.js';
import { LetterPlaySet, UserSession } from '../common/types.js';

document.addEventListener("DOMContentLoaded", () =>
{
  try {
    const username = getUsername();
    fetchLetters();
  } catch (error) {
    if (error instanceof Error) {
      console.error(error.message);
    } else {
      console.error("An unknown error occurred:", error);
    }
    // window.location.href = "../list-sessions/index.html";
  }

  const listSessionButton = document.getElementById("list-sessions-button");
  if (listSessionButton) {
    listSessionButton.addEventListener("click", () => { window.open("../list-sessions/index.html", "_blank"); });
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

    fetch(`${API_BASE_URL}/play-move?username=${username}`, {
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
      .then(response => handleResponse<UserSession>(response))
      .then((data: UserSession) =>
      {
        createUserSessionLayout(username, data);
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
    if (!confirm('Are you sure you want to reset the game?')) {
      return;
    }
    const username = getUsername();
    fetch(`${API_BASE_URL}/reset?username=${username}`, {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data: UserSession) =>
      {
        createUserSessionLayout(username, data);
        showMessage("Game reset successfully.");
      })
      .catch((error) =>
      {
        console.error("Error resetting letters:", error);
      });
  });

  const endSessionButton = document.getElementById("end-session-button");
  if (!endSessionButton) {
    console.error("End Session button not found");
    return;
  }
  endSessionButton.addEventListener("click", () => {
    if (!confirm("Are you sure you want to end the session?")) {
      return;
    }

    const username = getUsername();
    fetch(`${API_BASE_URL}/end-session?username=${username}`, {
      method: "POST",
    })
      .then(response => {
        if (!response.ok) {
          throw new Error("Failed to end session");
        }
        return response.text();
      })
      .then(message => {
        showMessage(message);
        window.location.href = "../list-sessions/index.html"; // Redirect to session list
      })
      .catch(error => {
        console.error("Error ending session:", error);
        showMessage(error.message);
      });
  });
});

function fetchLetters()
{
  const username = getUsername();
  fetch(`${API_BASE_URL}/letters?username=${username}`)
    .then(response => handleResponse<UserSession>(response))
    .then((data: UserSession) =>
    {
      createUserSessionLayout(username, data);
    })
    .catch((error) =>
    {
      console.error("Error fetching letter data:", error);
      showMessage(error.message);
    });
}

function createUserSessionLayout(username: string, data: UserSession)
{

  updateTextContent("username", `Username: ${username}`);
  updateTextContent("session-start-timestamp", `Session Start: ${data.session_start_timestamp}`);
  updateTextContent("last-move-timestamp", `Last Move: ${data.last_move_timestamp}`);
  updateTextContent("overall-value", `Overall Letter Value: ${data.letter_overall_value}`);

  const totalRemaining = data.letters_play_set.reduce(
    (sum, letter) => sum + letter.current_count, 0
  );
  updateTextContent("remaining-letters", `Remaining Letters: ${totalRemaining}`);

  const letterContainer = getElementByIdOrThrow<HTMLElement>("letters-play-set");
  letterContainer.innerHTML = "";
  letterContainer.appendChild(createLettersTable(data));
}

function getUsername()
{
  const urlParams = new URLSearchParams(window.location.search);
  const usernameFromQuery = urlParams.get("username");

  if (usernameFromQuery) {
    const usernameInput = getElementByIdOrThrow<HTMLInputElement>("username");
    usernameInput.value = usernameFromQuery;
    return usernameFromQuery;
  }

  showMessage("No username provided. Please create a session first.");
  throw new Error("Username is required");
}

function createLettersTable(data: UserSession)
{
  const table = document.createElement("table");
  table.classList.add("letter-table");

  let row: HTMLTableRowElement;
  let index = 0;
  data.letters_play_set.forEach((letter) =>
  {
    // if (letter.count === 0) {
    //   // Skip letters with a count of 0
    //   return;
    // }

    if (index % 5 === 0) {
      // Create a new row every 5 letters
      row = document.createElement("tr");
      table.appendChild(row);
    }
    index++;

    const letterCell = document.createElement("td");
    letterCell.textContent = letter.letter;
    letterCell.classList.add("letter-cell");

    const countCell = document.createElement("td");
    countCell.textContent = `${letter.current_count.toString()} / ${letter.original_count.toString()}`;
    countCell.classList.add("count-cell");
    countCell.style.backgroundColor = getBackgroundColor(letter);

    row.appendChild(letterCell);
    row.appendChild(countCell);
  });
  return table;
}

function getBackgroundColor(letter: LetterPlaySet): string
{
  const reminding_percentage = (letter.current_count / letter.original_count) * 100;
  console.log(`Letter: ${letter.letter}, Current Count: ${letter.current_count}, Original Count: ${letter.original_count}, Percentage: ${reminding_percentage}`);
  if (reminding_percentage >= 100) {
    return "lightgreen";
  } else if (reminding_percentage > 75) {
    return "lightblue";
  } else if (reminding_percentage > 50) {
    return "lightpink";
  } else if (reminding_percentage > 25) {
    return "lightyellow";
  } else if (reminding_percentage > 0) {
    return "lightaquamarine";
  } else if (reminding_percentage <= 0) {
    return "lightcoral";
  } else {
    return "black";
  }
}
