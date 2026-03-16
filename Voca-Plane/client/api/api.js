import { BASE_URL } from '../config/config.js';
import { getCookie } from '../utils/cookie.js';

export async function fetchApi(endpoint, options = {}) {
    // Determine the full URL
    const url = endpoint.startsWith('http') ? endpoint : `${BASE_URL}${endpoint}`;
    
    // Default headers
    const headers = new Headers(options.headers || {});
    
    // Auto-attach JSON Content-Type if there's a body and no Content-Type was explicitly set
    if (options.body && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json');
    }
    
    // Auto-attach Authorization token
    const token = getCookie('access_token');
    if (token) {
        headers.set('Authorization', `Bearer ${token}`);
    }

    const config = {
        ...options,
        headers,
    };

    // If body is object and not FormData, stringify it
    if (config.body && typeof config.body === 'object' && !(config.body instanceof FormData)) {
        config.body = JSON.stringify(config.body);
    }

    try {
        const response = await fetch(url, config);
        
        let data;
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.indexOf("application/json") !== -1) {
            data = await response.json();
        } else {
            data = await response.text();
        }

        if (!response.ok) {
            // Throw the generic response error wrapper (usually has data.message or similar)
            throw new Error((data && data.message) ? data.message : `HTTP error! status: ${response.status}`);
        }

        return data; // Usually data is { status, message, data: {...} } based on the gin backend
    } catch (error) {
        console.error('Fetch API Error:', error);
        throw error;
    }
}
