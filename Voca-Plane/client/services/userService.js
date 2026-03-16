import { fetchApi } from '../api/api.js';

export const userService = {
    async getProfile() {
        return fetchApi('/user/profile');
    },

    async updateProfile(data) {
        return fetchApi('/user/profile', {
            method: 'PATCH',
            body: data
        });
    }
};
