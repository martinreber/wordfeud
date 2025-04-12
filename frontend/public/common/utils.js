export function showMessage(message) {
    const messageContainer = document.getElementById("response-message");
    messageContainer.textContent = message;
}

export function getUsername() {
    const urlParams = new URLSearchParams(window.location.search);
    const usernameFromQuery = urlParams.get("username");

    if (usernameFromQuery) {
        return usernameFromQuery;
    }

    showMessage("No username provided. Please create a session first.");
    throw new Error("Username is required");
}