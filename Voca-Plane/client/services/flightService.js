import { fetchApi } from '../api/api.js';

export const flightService = {
    async getAllFlights(page = 1, limit = 10) {
        return fetchApi(`/flights?page=${page}&limit=${limit}`);
    },

    async getAllFlightsFull() {
        return fetchApi('/flights/all');
    },

    async searchFlights(params, page = 1, limit = 10) {
        const query = new URLSearchParams({ ...params, page, limit }).toString();
        return fetchApi(`/flights/search?${query}`);
    },

    async getFlightById(id) {
        return fetchApi(`/flights/${id}`);
    },

    async getFlightSeats(id) {
        return fetchApi(`/flights/${id}/seats`);
    }
};
