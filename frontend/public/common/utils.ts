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

    showMessage("No username provided. Please create a game first.");
    throw new Error("Username is required");
}

export async function handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
            const errorData = await response.json();
            throw new Error(errorData.message || "Request failed");
        } else {
            const errorText = await response.text();
            throw new Error(errorText || "Request failed");
        }
    }

    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        return response.json();
    }

    // Return empty object for non-JSON responses
    return {} as T;
}

export function getElementByIdOrThrow<T extends HTMLElement>(id: string): T {
    const element = document.getElementById(id) as T;
    if (!element) {
        throw new Error(`Element with id '${id}' not found`);
    }
    return element;
}

export function updateTextContent(id: string, text: string): void {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = text;
    }
}

export const API_BASE_URL = "http://localhost:8080";