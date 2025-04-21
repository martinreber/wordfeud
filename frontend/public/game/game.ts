import
{
  showMessage,
  handleResponse,
  getElementByIdOrThrow,
  updateTextContent,
  API_BASE_URL
} from '../common/utils.js';
import { LetterPlaySet, UserGame } from '../common/types.js';

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
    // window.location.href = "../list-games/index.html";
  }

  const listGameButton = document.getElementById("list-games-button");
  if (listGameButton) {
    listGameButton.addEventListener("click", () => { window.open("../list-games/index.html", "_blank"); });
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
    const inputPointsElement = document.getElementById("input-points") as HTMLInputElement;
    const inputPoints = inputPointsElement ? inputPointsElement.value : "";
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

    if (!inputPoints) {
      showMessage("Please enter the points you scored.");
      return;
    }

    // check if inputWord have only letters and , characters
    if (!/^[a-z\,]+$/.test(inputWord)) {
      showMessage("Word can only contain letters and , to split multiple words.");
      return;
    }

    fetch(`${API_BASE_URL}/games/${username}/play-move`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        letters: inputString,
        word: "", // deprecated
        words: inputWord.split(","),
        playedByMyself: isPlayedByMyself,
        points: parseInt(inputPoints),
      }),
    })
      .then(async response =>
      {
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Failed to play move');
        }
        return handleResponse<UserGame>(response);
      })
      .then((data: UserGame) =>
      {
        createUserGameLayout(username, data);
        const player = isPlayedByMyself ? "myself" : "opponent";
        showMessage(`Word "${inputWord}" successfully played by ${player} with ${inputPoints} points.`);
        inputStringElement.value = "";
        inputWordElement.value = "";
        inputPointsElement.value = "";
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
    fetch(`${API_BASE_URL}/games/${username}/reset`, {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data: UserGame) =>
      {
        createUserGameLayout(username, data);
        showMessage("Game reset successfully.");
      })
      .catch((error) =>
      {
        console.error("Error resetting letters:", error);
      });
  });

  const endGameButton = document.getElementById("end-game-button");
  if (!endGameButton) {
    console.error("End Game button not found");
    return;
  }
  endGameButton.addEventListener("click", () =>
  {
    if (!confirm("Are you sure you want to end the game?")) {
      return;
    }

    const username = getUsername();
    fetch(`${API_BASE_URL}/games/${username}/end`, {
      method: "POST",
    })
      .then(response =>
      {
        if (!response.ok) {
          throw new Error("Failed to end game");
        }
        return response.text();
      })
      .then(message =>
      {
        showMessage(message);
        window.location.href = "../list-games/index.html"; // Redirect to game list
      })
      .catch(error =>
      {
        console.error("Error ending game:", error);
        showMessage(error.message);
      });
  });
});

function fetchLetters()
{
  const username = getUsername();
  fetch(`${API_BASE_URL}/games/${username}`)
    .then(response => handleResponse<UserGame>(response))
    .then((data: UserGame) =>
    {
      createUserGameLayout(username, data);
    })
    .catch((error) =>
    {
      console.error("Error fetching letter data:", error);
      showMessage(error.message);
    });
}

function createUserGameLayout(username: string, data: UserGame)
{

  updateTextContent("username", `Username: ${username}`);
  updateTextContent("game-start-timestamp", `Game Start: ${data.game_start_timestamp}`);
  updateTextContent("last-move-timestamp", `Last Move: ${data.last_move_timestamp}`);
  updateTextContent("overall-value", `Overall Letter Value: ${data.letter_overall_value}`);

  const playerToggleElement = document.getElementById("player-toggle") as HTMLInputElement;
  playerToggleElement.checked = !playerToggleElement.checked;
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

  showMessage("No username provided. Please create a game first.");
  throw new Error("Username is required");
}

function createLettersTable(data: UserGame)
{
  const table = document.createElement("table");
  table.classList.add("letter-table");

  let row: HTMLTableRowElement;
  let index = 0;
  data.letters_play_set.forEach((letter) =>
  {
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
    countCell.style.color = "#333";

    row.appendChild(letterCell);
    row.appendChild(countCell);
  });
  return table;
}

function getBackgroundColor(letter: LetterPlaySet): string
{
  if (letter.current_count >= 5) {
    return "#00ff00";
  } else if (letter.current_count == 4) {
    return "#bfff00";
  } else if (letter.current_count == 3) {
    return "#ffff00";
  } else if (letter.current_count == 2) {
    return "#ff8000";
  } else if (letter.current_count == 1) {
    return "#ff4000";
  } else if (letter.current_count <= 0) {
    return "#ff0000";
  }
  else {
    return "black";
  }
}
