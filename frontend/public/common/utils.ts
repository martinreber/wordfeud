export function showMessage(message: string): void {
    const messageContainer = document.getElementById("response-message");
    if (messageContainer) {
        messageContainer.textContent = message;
    }
}

export function getUsername(): string {
    const urlParams = new URLSearchParams(window.location.search);
    const usernameFromQuery = urlParams.get("username");

    if (usernameFromQuery) {
        return usernameFromQuery;
    }

    showMessage("No username provided. Please create a session first.");
    throw new Error("Username is required");
}