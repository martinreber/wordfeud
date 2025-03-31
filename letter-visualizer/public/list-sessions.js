function fetchSessions() {
    fetch("http://localhost:8080/list")
        .then((response) => response.json())
        .then((data) => {
            const tableBody = document.querySelector("#sessions-table tbody");
            tableBody.innerHTML = ""; // Clear existing rows

            // Sort sessions by username (alphabetically)
            data.sessions.sort((a, b) => a.user.localeCompare(b.user));

            data.sessions.forEach((session) => {
                const row = document.createElement("tr");

                const usernameCell = document.createElement("td");
                const usernameLink = document.createElement("a");
                usernameLink.textContent = session.user;
                usernameLink.href = `session.html?username=${encodeURIComponent(session.user)}`;
                usernameLink.target = "_blank"; // Open in a new tab
                usernameCell.appendChild(usernameLink);

                const sessionStartCell = document.createElement("td");
                sessionStartCell.textContent = session.session_start_timestamp;

                const timestampCell = document.createElement("td");
                timestampCell.textContent = session.last_move_timestamp

                // Add the remaining letters count
                const remainingLettersCell = document.createElement("td");
                remainingLettersCell.textContent = session.reminding_letters;

                // Add a delete button
                const deleteCell = document.createElement("td");
                const deleteButton = document.createElement("button");
                deleteButton.textContent = "Delete";
                deleteButton.classList.add("button", "delete-button");
                deleteButton.addEventListener("click", () => deleteSession(session.user));
                deleteCell.appendChild(deleteButton);

                // Append cells to the row
                row.appendChild(usernameCell);
                row.appendChild(sessionStartCell);
                row.appendChild(timestampCell);
                row.appendChild(remainingLettersCell);
                row.appendChild(deleteCell);
                tableBody.appendChild(row);
            });
        })
        .catch((error) => {
            console.error("Error fetching sessions:", error);
        });
}

function deleteSession(username) {
    if (!confirm(`Are you sure you want to delete the session for "${username}"?`)) {
        return; // Exit if the user cancels the confirmation
    }

    fetch(`http://localhost:8080/delete?username=${encodeURIComponent(username)}`, {
        method: "DELETE",
    })
        .then((response) => {
            if (response.ok) {
                showMessage(`Session for "${username}" deleted successfully.`);
                fetchSessions(); // Refresh the session list
            } else {
                return response.json().then((data) => {
                    const errorMessage = data.message || "Failed to delete the session.";
                    showMessage(errorMessage);
                });
            }
        })
        .catch((error) => {
            console.error("Error deleting session:", error);
            showMessage("An unexpected error occurred. Please try again.");
        });
}
// Create a new session
document.getElementById("create-session-button").addEventListener("click", () => {
    const newUsernameInput = document.getElementById("new-username");
    const newUsername = newUsernameInput.value.trim();

    if (!newUsername) {
        showMessage("Please enter a username.");
        return;
    }

    if (newUsername.length > 20) {
        showMessage("Username cannot exceed 20 characters.");
        return;
    }

    // Create a new session by sending a request to the server
    fetch(`http://localhost:8080/create?username=${encodeURIComponent(newUsername)}`, {
        method: "POST",
    })
    .then((response) => {
        if (response.ok) {
            // Redirect to the session page with the username as a query parameter
            fetchSessions();
            window.open(`session.html?username=${encodeURIComponent(newUsername)}`, "_blank");
            showMessage(""); // Clear any previous message
        } else {
            // Extract the error message from the response body
            const errorMessage = "username already exists";
            showMessage(errorMessage); // Show the error message in the UI
        }
    })
    .catch((error) => {
        console.error("Error creating session:", error);
        showMessage("An unexpected error occurred. Please try again.");
    });
});

// Add event listener for the "Refresh" button
document.getElementById("refresh-button").addEventListener("click", () => {
    fetchSessions();
});

// Fetch sessions on page load
fetchSessions();


function showMessage(message) {
    const messageContainer = document.getElementById("response-message");
    messageContainer.textContent = message; // Set the message text
}