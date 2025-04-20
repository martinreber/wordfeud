import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';
import { EndedGame, Game } from '../common/types.js';

async function fetchGames(): Promise<void>
{
    try {
        const response = await fetch(`${API_BASE_URL}/games/end-game`);

        const data = await handleResponse<EndedGame[]>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>('games-table').querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        tableBody.innerHTML = "";

        if (data) {
            data.sort((a, b) => a.user.localeCompare(b.user));
            data.forEach((endedGame) =>
            {
                const row = createGameRow(endedGame);
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error("Error fetching games:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

function createGameRow(endedGame: EndedGame): HTMLTableRowElement
{
    const row = document.createElement("tr");
    row.appendChild(createCell(endedGame.user));
    row.appendChild(createCell(endedGame.game_start_timestamp));
    row.appendChild(createCell(endedGame.last_move_timestamp));
    return row;
}

function createCell(content: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    cell.textContent = content;
    return cell;
}

// Event Listeners
document.addEventListener("DOMContentLoaded", () =>
{
    fetchGames();
});