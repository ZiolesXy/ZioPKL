import { fetchApi } from '../api/api.js';

export const authService = {
    async register(data) {
        return fetchApi('/auth/register', {
            method: 'POST',
            body: data
        });
    },

    async login(data) {
        return fetchApi('/auth/login', {
            method: 'POST',
            body: data
        });
    },

    async refreshToken(refreshTokenStr) {
        return fetchApi('/auth/refresh', {
            method: 'POST',
            body: { refresh_token: refreshTokenStr }
        });
    }
};
