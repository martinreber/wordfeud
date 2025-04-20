import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';
import { Game } from '../common/types.js';

async function fetchGames(): Promise<void>
{
    try {
        const response = await fetch(`${API_BASE_URL}/games`);

        const data = await handleResponse<Game[]>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>('games-table').querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        tableBody.innerHTML = ""; // Clear existing rows

        // Sort games by username (alphabetically)
        console.log(data)
        if (data) {
            data.sort((a, b) => a.user.localeCompare(b.user));
            data.forEach((game) =>
            {
                const row = createGameRow(game);
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error("Error fetching games:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

function createGameRow(game: Game): HTMLTableRowElement
{
    const row = document.createElement("tr");

    row.appendChild(createUsernameCell(game.user));
    row.appendChild(createCell(game.game_start_timestamp));
    row.appendChild(createCell(game.last_move_timestamp));
    row.appendChild(createCell(game.reminding_letters.toString()));
    row.appendChild(EndGameCell(game.user));

    return row;
}

function createUsernameCell(username: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    const link = document.createElement("a");
    link.textContent = username;
    link.href = `../game/index.html?username=${encodeURIComponent(username)}`;
    link.target = "_blank";
    cell.appendChild(link);
    return cell;
}

function createCell(content: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    cell.textContent = content;
    return cell;
}

function EndGameCell(username: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    const button = document.createElement("button");
    button.textContent = "End Game";
    button.classList.add("button", "end-game-button");
    button.addEventListener("click", () => endGame(username));
    cell.appendChild(button);
    return cell;
}

async function endGame(username: string): Promise<void>
{
    if (!confirm(`Are you sure you want to terminate the game with "${username}"?`)) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/games/${encodeURIComponent(username)}/end-game`, {
            method: "POST",
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || "Failed to end game");
        }

        showMessage(`Game for "${username}" ended successfully.`);
        await fetchGames();
    } catch (error) {
        console.error("Error ending game:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

async function createGame(username: string): Promise<void>
{
    try {
        const response = await fetch(`${API_BASE_URL}/games/${encodeURIComponent(username)}`, {
            method: "POST",
        });
        await handleResponse(response);
        await fetchGames();
        window.open(`../game/index.html?username=${encodeURIComponent(username)}`, "_blank");
        showMessage(""); // Clear any previous message
    } catch (error) {
        console.error("Error creating game:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

// Event Listeners
document.addEventListener("DOMContentLoaded", () =>
{
    const createGameButton = getElementByIdOrThrow<HTMLButtonElement>("create-game-button");
    createGameButton.addEventListener("click", async () =>
    {
        const newUsernameInput = getElementByIdOrThrow<HTMLInputElement>("new-username");
        const newUsername = newUsernameInput.value.trim();

        if (!newUsername) {
            showMessage("Please enter a username.");
            return;
        }

        if (newUsername.length > 20) {
            showMessage("Username cannot exceed 20 characters.");
            return;
        }

        await createGame(newUsername);
    });

    const refreshButton = getElementByIdOrThrow<HTMLButtonElement>("refresh-button");
    refreshButton.addEventListener("click", () => fetchGames());

    const playedWordsButton = getElementByIdOrThrow<HTMLButtonElement>("played-words-button");
    playedWordsButton.addEventListener("click", () =>
    {
        window.open("../played-words/index.html", "_blank");
    });

    // Initial fetch
    fetchGames();
});