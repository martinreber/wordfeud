import { showMessage, handleResponse, getElementByIdOrThrow, API_BASE_URL } from '../common/utils.js';
import { Session } from '../common/types.js';

async function fetchSessions(): Promise<void>
{
    try {
        const response = await fetch(`${API_BASE_URL}/list`);

        const data = await handleResponse<Session[]>(response);

        const tableBody = getElementByIdOrThrow<HTMLTableSectionElement>('sessions-table').querySelector('tbody');
        if (!tableBody) throw new Error('Table body not found');

        tableBody.innerHTML = ""; // Clear existing rows

        // Sort sessions by username (alphabetically)
        console.log(data)
        if (data) {
            data.sort((a, b) => a.user.localeCompare(b.user));
            data.forEach((session) =>
            {
                const row = createSessionRow(session);
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error("Error fetching sessions:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

function createSessionRow(session: Session): HTMLTableRowElement
{
    const row = document.createElement("tr");

    row.appendChild(createUsernameCell(session.user));
    row.appendChild(createCell(session.session_start_timestamp));
    row.appendChild(createCell(session.last_move_timestamp));
    row.appendChild(createCell(session.reminding_letters.toString()));
    row.appendChild(createDeleteCell(session.user));

    return row;
}

function createUsernameCell(username: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    const link = document.createElement("a");
    link.textContent = username;
    link.href = `../session/index.html?username=${encodeURIComponent(username)}`;
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

function createDeleteCell(username: string): HTMLTableCellElement
{
    const cell = document.createElement("td");
    const button = document.createElement("button");
    button.textContent = "Delete";
    button.classList.add("button", "delete-button");
    button.addEventListener("click", () => deleteSession(username));
    cell.appendChild(button);
    return cell;
}

async function deleteSession(username: string): Promise<void>
{
    if (!confirm(`Are you sure you want to delete the session for "${username}"?`)) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/delete?username=${encodeURIComponent(username)}`, {
            method: "DELETE",
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || "Failed to delete session");
        }

        showMessage(`Session for "${username}" deleted successfully.`);
        await fetchSessions();
    } catch (error) {
        console.error("Error deleting session:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

async function createSession(username: string): Promise<void>
{
    try {
        const response = await fetch(`${API_BASE_URL}/create?username=${encodeURIComponent(username)}`, {
            method: "POST",
        });
        await handleResponse(response);
        await fetchSessions();
        window.open(`../session/index.html?username=${encodeURIComponent(username)}`, "_blank");
        showMessage(""); // Clear any previous message
    } catch (error) {
        console.error("Error creating session:", error);
        showMessage(error instanceof Error ? error.message : "An unexpected error occurred");
    }
}

// Event Listeners
document.addEventListener("DOMContentLoaded", () =>
{
    const createSessionButton = getElementByIdOrThrow<HTMLButtonElement>("create-session-button");
    createSessionButton.addEventListener("click", async () =>
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

        await createSession(newUsername);
    });

    const refreshButton = getElementByIdOrThrow<HTMLButtonElement>("refresh-button");
    refreshButton.addEventListener("click", () => fetchSessions());

    const playedWordsButton = getElementByIdOrThrow<HTMLButtonElement>("played-words-button");
    playedWordsButton.addEventListener("click", () =>
    {
        window.open("../played-words/index.html", "_blank");
    });

    // Initial fetch
    fetchSessions();
});