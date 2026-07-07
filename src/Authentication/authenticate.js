const authenticate = async () => {
    const apiBaseUrl = (
        process.env.REACT_APP_API_BASE_URL
        || process.env.REACT_APP_BACKEND_API
        || 'https://filiya-backend.onrender.com/api'
    ).replace(/\/+$/, '');

    const response = await fetch(`${apiBaseUrl}/auth/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({}),
    })
    if (!response.ok) {
        return response.json()
    }
    return response.json()
}

export default authenticate;
