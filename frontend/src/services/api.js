const API_URL =
    import.meta.env.VITE_API_URL || "http://localhost:8080";

const WS_URL =
    import.meta.env.VITE_WS_URL || "ws://localhost:8080";

export async function signup(name, email, password) {
    const response = await fetch(`${API_URL}/api/auth/signup`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, email, password }),
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || "Signup failed");
    }

    return data;
}

export async function login(email, password) {
    const response = await fetch(`${API_URL}/api/auth/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || "Login failed");
    }

    return data;
}

export async function createPoll(token, question, options) {
    const response = await fetch(`${API_URL}/api/polls`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ question, options }),
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || "Failed to create poll");
    }

    return data;
}

export async function getPoll(pollId) {
    const response = await fetch(`${API_URL}/api/polls/${pollId}`);
    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || "Failed to load poll");
    }

    return data;
}

export async function vote(token, pollId, optionId) {
    const response = await fetch(
        `${API_URL}/api/polls/${pollId}/vote`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ optionId }),
        }
    );

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || "Vote failed");
    }

    return data;
}

export function connectToPoll(pollId, onMessage, onError) {
    const ws = new WebSocket(
    `${WS_URL}/ws/polls/${pollId}`    );

    ws.onmessage = (event) => {
        try {
            onMessage(JSON.parse(event.data));
        } catch (error) {
            console.error("Invalid WebSocket message:", error);
        }
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);

        if (onError) {
            onError(error);
        }
    };

    return ws;
}
